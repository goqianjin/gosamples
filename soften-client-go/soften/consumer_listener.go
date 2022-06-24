package soften

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
	"github.com/panjf2000/ants/v2"
	"github.com/shenqianjin/soften-client-go/soften/checker"
	"github.com/shenqianjin/soften-client-go/soften/config"
	"github.com/shenqianjin/soften-client-go/soften/internal"
	"github.com/shenqianjin/soften-client-go/soften/message"
	"github.com/shenqianjin/soften-client-go/soften/topic"
)

type consumeListener struct {
	// pulsar.Consumer
	client            *client
	logger            log.Logger
	messageCh         chan ConsumerMessage // channel used to deliver message to application
	enables           *internal.StatusEnables
	concurrency       *config.ConcurrencyPolicy
	generalHandlers   *generalConsumeHandlers
	levelHandlers     map[internal.TopicLevel]*leveledConsumeHandlers
	checkers          map[internal.CheckType]*checker.Checkpoint
	startListenerOnce sync.Once
}

func newConsumeListener(cli *client, conf config.ConsumerConfig) (*consumeListener, error) {
	logTopic := conf.Topics[0]
	if len(conf.Topics) > 1 {
		logTopic = logTopic + "+" + strconv.Itoa(len(conf.Topics)-1)
	}
	facade := &consumeListener{
		client:      cli,
		messageCh:   make(chan ConsumerMessage, 10),
		logger:      cli.logger.SubLogger(log.Fields{"Topic": logTopic}),
		concurrency: conf.Concurrency,
	}
	// collect enables
	facade.enables = facade.collectEnables(&conf)
	// initialize checkers
	/*if checkers, err := newConsumeCheckers(facade.logger, facade.enables, checkpointMap); err != nil {
		return nil, err
	} else {
		facade.checkers = checkers
	}*/
	// initialize general handlers
	generalHdOptions := facade.formatGeneralHandlersOptions(conf.Topics[0], &conf)
	if handlers, err := newGeneralConsumeHandlers(cli, generalHdOptions); err != nil {
		return nil, err
	} else {
		facade.generalHandlers = handlers
	}
	// initialize level related handlers
	facade.levelHandlers = make(map[internal.TopicLevel]*leveledConsumeHandlers, len(conf.Levels))
	for _, level := range conf.Levels {
		suffix, err := topic.NameSuffixOf(level)
		if err != nil {
			return nil, err
		}
		options := facade.formatLeveledHandlersOptions(conf.Topics[0]+suffix, &conf)
		if handlers, err := newLeveledConsumeHandlers(cli, options, facade.generalHandlers.deadHandler); err != nil {
			return nil, err
		} else {
			facade.levelHandlers[level] = handlers
		}
	}
	// initialize status multiStatusConsumer
	if len(conf.Levels) == 1 {
		level := conf.Levels[0]
		if _, err := newMultiStatusConsumer(facade.logger, cli, level, &conf, facade.messageCh, facade.levelHandlers[level]); err != nil {
			return nil, err
		}
	} else {
		if _, err := newMultiLevelConsumer(facade.logger, cli, &conf, facade.messageCh, facade.levelHandlers); err != nil {
			return nil, err
		}
	}

	facade.logger.Infof("complete initializing consume listener: topic: %v, levels: %v", conf.Topics, conf.Levels)
	return facade, nil
}

func (c *consumeListener) Start(ctx context.Context, handler Handler, checkpoints ...checker.Checkpoint) error {
	// convert handler
	premiumHandler := func(message pulsar.Message) HandleStatus {
		success, err := handler(message)
		if success {
			return HandleStatusOk
		} else {
			return HandleStatusFail.Err(err)
		}
	}
	// forward the call to c.SubscribePremium
	return c.StartPremium(ctx, premiumHandler, checkpoints...)
}

// StartPremium blocking to consume message one by one. it returns error if any parameters is invalid
func (c *consumeListener) StartPremium(ctx context.Context, handler PremiumHandler, checkpoints ...checker.Checkpoint) error {
	// validate handler
	if handler == nil {
		return errors.New("handler parameter is nil")
	}
	// validate checkpoints
	checkpointMap, err := checker.Validator.ValidateConsumeCheckpoint(checkpoints)
	if err != nil {
		return err
	}
	c.startListenerOnce.Do(func() {
		// initialize checkers
		c.checkers = c.collectCheckers(c.enables, checkpointMap)
		// initialize task pool
		pool, onceErr := ants.NewPool(int(c.concurrency.CorePoolSize), ants.WithExpiryDuration(60*time.Second))
		if onceErr != nil {
			err = onceErr
			return
		}
		// receive msg and then consume one by one
		c.internalStartInPool(ctx, handler, pool)
	})
	return nil
}

func (c *consumeListener) internalStartInPool(ctx context.Context, handler PremiumHandler, pool *ants.Pool) {
	// receive msg and submit task
	for {
		select {
		case msg, ok := <-c.messageCh:
			if !ok {
				return
			}

			if err := pool.Submit(func() { c.consume(handler, msg) }); err != nil {
				c.logger.Errorf("submit msg failed. err: %v", err)
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (c *consumeListener) internalStartInParallel(ctx context.Context, handler PremiumHandler) {
	concurrencyChan := make(chan bool, c.concurrency.CorePoolSize)
	for {
		select {
		case msg, ok := <-c.messageCh:
			if !ok {
				return
			}
			concurrencyChan <- true
			go func(msg ConsumerMessage) {
				c.consume(handler, msg)
				<-concurrencyChan
			}(msg)
		case <-ctx.Done():
			return
		}

	}
}

func (c *consumeListener) consume(handler PremiumHandler, message ConsumerMessage) {
	// pre-check to handle in turn
	for _, checkType := range checker.PreCheckTypes() {
		if checkpoint, ok := c.checkers[checkType]; ok && checkpoint.Before != nil {
			checkStatus := checkpoint.Before(message)
			if handledDeferFunc := checkStatus.GetHandledDefer(); handledDeferFunc != nil {
				defer handledDeferFunc()
			}
			if !checkStatus.IsPassed() {
				continue
			}
			if checkHandler := c.parseHandler(checkType, message); checkHandler != nil {
				if checkHandler.Handle(message.ConsumerMessage, checkStatus) {
					// return to skip biz handler if check handle succeeded
					return
				}
			}

		}
	}

	// do biz handle
	bizHandleStatus := handler(message)

	// post-check to route - for obvious goto action
	if bizHandleStatus.getGotoAction() != "" {
		if ok := c.handleMessageGotoAction(message, bizHandleStatus.getGotoAction(), checker.CheckStatusPassed); ok {
			return
		}
	}

	// post-check to route - for obvious checkers or configured checkers
	postCheckTypesInTurn := checker.DefaultPostCheckTypes()
	if len(bizHandleStatus.getCheckTypes()) > 0 {
		postCheckTypesInTurn = bizHandleStatus.getCheckTypes()
	}
	for _, checkType := range postCheckTypesInTurn {
		if checkpoint, ok := c.checkers[checkType]; ok && checkpoint.After != nil {
			checkStatus := checkpoint.After(message, bizHandleStatus.getErr())
			if handledDeferFunc := checkStatus.GetHandledDefer(); handledDeferFunc != nil {
				defer handledDeferFunc()
			}
			if !checkStatus.IsPassed() {
				continue
			}
			if checkHandler := c.parseHandler(checkType, message); checkHandler != nil {
				if checkHandler.Handle(message.ConsumerMessage, checkStatus) {
					// return if check handle succeeded
					return
				}
			}

		}
	}

	// here means to let application client Ack/Nack message
	return
}

func (c *consumeListener) collectCheckers(enables *internal.StatusEnables, checkpointMap map[internal.CheckType]*checker.Checkpoint) map[internal.CheckType]*checker.Checkpoint {
	checkers := make(map[internal.CheckType]*checker.Checkpoint)
	if enables.RerouteEnable {
		checkers[checker.CheckTypeReroute] = c.parseConfiguredChecker(checker.CheckTypeReroute, checkpointMap)
	}
	if enables.PendingEnable {
		checkers[checker.CheckTypePending] = c.parseConfiguredChecker(checker.CheckTypePending, checkpointMap)
	}
	if enables.BlockingEnable {
		checkers[checker.CheckTypeBlocking] = c.parseConfiguredChecker(checker.CheckTypeBlocking, checkpointMap)
	}
	if enables.RetryingEnable {
		checkers[checker.CheckTypeRetrying] = c.parseConfiguredChecker(checker.CheckTypeRetrying, checkpointMap)
	}
	if enables.DeadEnable {
		checkers[checker.CheckTypeDead] = c.parseConfiguredChecker(checker.CheckTypeDead, checkpointMap)
	}
	if enables.DiscardEnable {
		checkers[checker.CheckTypeDiscard] = c.parseConfiguredChecker(checker.CheckTypeDiscard, checkpointMap)
	}
	if enables.UpgradeEnable {
		checkers[checker.CheckTypeUpgrade] = c.parseConfiguredChecker(checker.CheckTypeUpgrade, checkpointMap)
	}
	if enables.DegradeEnable {
		checkers[checker.CheckTypeDegrade] = c.parseConfiguredChecker(checker.CheckTypeDegrade, checkpointMap)
	}
	return checkers
}

func (c *consumeListener) parseConfiguredChecker(checkType internal.CheckType, checkpointMap map[internal.CheckType]*checker.Checkpoint) *checker.Checkpoint {
	if ckp, ok := checkpointMap[checkType]; ok {
		return ckp
	}
	return checker.NilCheckpoint
}

func (c *consumeListener) parseBeforeChecker(checkType internal.CheckType, checkpointMap map[internal.CheckType]*checker.Checkpoint) *checker.Checkpoint {
	if ckp, ok := checkpointMap[checkType]; ok && ckp.Before != nil {
		return ckp
	}
	return nil
}

func (c *consumeListener) parseAfterChecker(checkType internal.CheckType, checkpointMap map[internal.CheckType]*checker.Checkpoint) *checker.Checkpoint {
	if ckp, ok := checkpointMap[checkType]; ok && ckp.After != nil {
		return ckp
	}
	return nil
}

func (c *consumeListener) parseHandler(checkType internal.CheckType, message ConsumerMessage) internalHandler {
	l := message.Level()
	switch checkType {
	case checker.CheckTypePreDiscard:
		return c.generalHandlers.discardHandler
	case checker.CheckTypePreDead:
		return c.generalHandlers.deadHandler
	case checker.CheckTypePreUpgrade:
		return c.levelHandlers[l].upgradeHandler
	case checker.CheckTypePreDegrade:
		return c.levelHandlers[l].degradeHandler
	case checker.CheckTypePreBlocking:
		return c.levelHandlers[l].blockingHandler
	case checker.CheckTypePrePending:
		return c.levelHandlers[l].pendingHandler
	case checker.CheckTypePreRetrying:
		return c.levelHandlers[l].retryingHandler
	case checker.CheckTypePreReroute:
		return c.generalHandlers.rerouteHandler
	}
	return nil
}

func (c *consumeListener) handleMessageGotoAction(consumerMessage ConsumerMessage, messageGoto internal.MessageGoto, cheStatus checker.CheckStatus) (routed bool) {
	l := consumerMessage.Level()
	msg := consumerMessage.ConsumerMessage
	switch messageGoto {
	case message.GotoDone:
		return c.generalHandlers.doneHandler.Handle(msg, cheStatus)
	case message.GotoPending:
		return c.enables.PendingEnable && c.levelHandlers[l].pendingHandler.Handle(msg, cheStatus)
	case message.GotoBlocking:
		return c.enables.BlockingEnable && c.levelHandlers[l].blockingHandler.Handle(msg, cheStatus)
	case message.GotoRetrying:
		return c.enables.RetryingEnable && c.levelHandlers[l].retryingHandler.Handle(msg, cheStatus)
	case message.GotoDead:
		return c.enables.DeadEnable && c.generalHandlers.deadHandler.Handle(msg, cheStatus)
	case message.GotoDiscard:
		return c.enables.DiscardEnable && c.generalHandlers.discardHandler.Handle(msg, cheStatus)
	case message.GotoUpgrade:
		return c.enables.UpgradeEnable && c.levelHandlers[l].upgradeHandler.Handle(msg, cheStatus)
	case message.GotoDegrade:
		return c.enables.DegradeEnable && c.levelHandlers[l].degradeHandler.Handle(msg, cheStatus)
	default:
		c.logger.Warnf("invalid msg goto action: %v", messageGoto)
	}
	return false
}

func (c *consumeListener) Close() {
	/*c.closeOnce.Do(func() {
		var wg sync.WaitGroup
		wg.Add(len(c.multiStatusConsumer))
		for _, con := range c.multiStatusConsumer {
			go func(multiStatusConsumeFacade Consumer) {
				defer wg.Done()
				multiStatusConsumeFacade.Close()
			}(con)
		}
		wg.Wait()
		close(c.closeCh)
		c.client.handlers.Del(c)
		c.dlq.close()
		c.rlq.close()
	})*/
}

func (c *consumeListener) collectEnables(conf *config.ConsumerConfig) *internal.StatusEnables {
	enables := internal.StatusEnables{
		ReadyEnable:    true,
		DeadEnable:     conf.DeadEnable,
		DiscardEnable:  conf.DiscardEnable,
		BlockingEnable: conf.BlockingEnable,
		PendingEnable:  conf.PendingEnable,
		RetryingEnable: conf.RetryingEnable,
		RerouteEnable:  conf.RerouteEnable,
		UpgradeEnable:  conf.UpgradeEnable,
		DegradeEnable:  conf.DegradeEnable,
	}
	return &enables
}

func (c *consumeListener) formatGeneralHandlersOptions(topic string, config *config.ConsumerConfig) generalConsumeHandlerOptions {
	options := generalConsumeHandlerOptions{
		Topic:         topic,
		DiscardEnable: config.BlockingEnable,
		DeadEnable:    config.RetryEnable,
		RerouteEnable: config.RerouteEnable,
	}
	return options
}

func (c *consumeListener) formatLeveledHandlersOptions(leveledTopic string, config *config.ConsumerConfig) leveledConsumeHandlerOptions {
	options := leveledConsumeHandlerOptions{
		Topic:             leveledTopic,
		BlockingEnable:    config.BlockingEnable,
		Blocking:          config.Blocking,
		PendingEnable:     config.PendingEnable,
		Pending:           config.Pending,
		RetryingEnable:    config.RetryEnable,
		Retrying:          config.Retrying,
		UpgradeEnable:     config.UpgradeEnable,
		UpgradeTopicLevel: config.UpgradeTopicLevel,
		DegradeEnable:     config.DegradeEnable,
		DegradeTopicLevel: config.DegradeTopicLevel,
	}
	return options
}
