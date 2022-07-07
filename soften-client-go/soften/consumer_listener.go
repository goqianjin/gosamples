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
	client               *client
	logger               log.Logger
	messageCh            chan ConsumerMessage // channel used to deliver message to application
	enables              *internal.StatusEnables
	concurrency          *config.ConcurrencyPolicy
	generalHandlers      *generalConsumeHandlers
	levelHandlers        map[internal.TopicLevel]*leveledConsumeHandlers
	checkers             map[internal.CheckType]*wrappedCheckpoint
	startListenerOnce    sync.Once
	closeListenerOnce    sync.Once
	multiLeveledConsumer *multiLeveledConsumer
	leveledConsumer      *leveledConsumer
	logTopics            string
	logLevels            string
	metrics              *internal.ListenMetrics
}

func newConsumeListener(cli *client, conf config.ConsumerConfig) (*consumeListener, error) {
	logTopic := conf.Topics[0]
	if len(conf.Topics) > 1 {
		logTopic = logTopic + "+" + strconv.Itoa(len(conf.Topics)-1)
	}
	logLevels := internal.TopicLevelParser.FormatList(conf.Levels)
	topicLogger := cli.logger.SubLogger(log.Fields{"Topic": logTopic})
	listener := &consumeListener{
		client:      cli,
		messageCh:   make(chan ConsumerMessage, 10),
		logger:      topicLogger.SubLogger(log.Fields{"level": logLevels}),
		concurrency: conf.Concurrency,
		metrics:     cli.metricsProvider.GetListenMetrics(logTopic, logLevels),
		logTopics:   logTopic,
		logLevels:   logLevels,
	}
	// collect enables
	listener.enables = listener.collectEnables(&conf)
	// initialize checkers
	/*if checkers, err := newConsumeCheckers(listener.logger, listener.enables, checkpointMap); err != nil {
		return nil, err
	} else {
		listener.checkers = checkers
	}*/
	// initialize general handlers
	generalHdOptions := listener.formatGeneralHandlersOptions(conf.Topics[0], &conf)
	if handlers, err := newGeneralConsumeHandlers(cli, generalHdOptions); err != nil {
		return nil, err
	} else {
		listener.generalHandlers = handlers
	}
	// initialize level related handlers
	listener.levelHandlers = make(map[internal.TopicLevel]*leveledConsumeHandlers, len(conf.Levels))
	for _, level := range conf.Levels {
		suffix, err := topic.NameSuffixOf(level)
		if err != nil {
			return nil, err
		}
		options := listener.formatLeveledHandlersOptions(conf.Topics[0]+suffix, &conf)
		if handlers, err := newLeveledConsumeHandlers(cli, options, listener.generalHandlers.deadHandler); err != nil {
			return nil, err
		} else {
			listener.levelHandlers[level] = handlers
		}
	}
	// initialize status leveledConsumer
	if len(conf.Levels) == 1 {
		level := conf.Levels[0]
		if _, err := newSingleLeveledConsumer(topicLogger, cli, level, &conf, listener.messageCh, listener.levelHandlers[level]); err != nil {
			return nil, err
		}
	} else {
		if _, err := newMultiLeveledConsumer(topicLogger, cli, &conf, listener.messageCh, listener.levelHandlers); err != nil {
			return nil, err
		}
	}

	listener.logger.Infof("created consume listener. topics: %v", conf.Topics)
	listener.metrics.ListenersOpened.Inc()
	return listener, nil
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
		c.logger.Info("started to listening...")
		c.metrics.ListenersRunning.Inc()
		// receive msg and then consume one by one
		c.internalStartInPool(ctx, handler, pool)
		c.metrics.ListenersRunning.Dec()
		c.logger.Info("ended to listening")
	})
	return nil
}

func (c *consumeListener) internalStartInPool(ctx context.Context, handler PremiumHandler, pool *ants.Pool) {
	// receive msg and submit task
	count := 0
	for {
		select {
		case msg, ok := <-c.messageCh:
			if !ok {
				return
			}
			count++
			// As pool.Submit is blocking, err happens only if pool is closed.
			// Namely, the 'err != nil' condition is never meet.
			if err := pool.Submit(func() { c.consume(handler, msg) }); err != nil {
				c.logger.Errorf("submit msg failed. err: %v", err)
				//msg.Consumer.Nack(msg.Message)
				//return
			}
			//c.logger.Infof("consume end -------+++++++++++++++++++++++++++- %d", count)
		case <-ctx.Done():
			c.logger.Warnf("closed soften listener")
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
	/*if checkHandler := c.parseHandler(checker.CheckTypePrePending, message); checkHandler != nil {
		if checkHandler.Decide(message.ConsumerMessage, checker.CheckStatusPassed) {
			// return to skip biz handler if check handle succeeded
			return
		}
		return
	}*/
	// pre-check to handle in turn
	for _, checkType := range checker.PreCheckTypes() {
		if checkpoint, ok := c.checkers[checkType]; ok && checkpoint.Before != nil {
			start := time.Now()
			checkStatus := checkpoint.Before(message)
			latency := time.Now().Sub(start).Seconds()
			if handledDeferFunc := checkStatus.GetHandledDefer(); handledDeferFunc != nil {
				defer handledDeferFunc()
			}
			checkpoint.metrics.CheckLatency.Observe(latency)
			if !checkStatus.IsPassed() {
				checkpoint.metrics.CheckRejected.Inc()
				continue
			} else {
				checkpoint.metrics.CheckPassed.Inc()
			}
			if checkHandler := c.parseHandler(checkType, message); checkHandler != nil {
				if checkHandler.Decide(message.ConsumerMessage, checkStatus) {
					// return to skip biz handler if check handle succeeded
					return
				}
			}

		}
	}

	message.Consumer.Name()
	// do biz handle
	start := time.Now()
	bizHandleStatus := handler(message)
	time.Now().Sub(start).Seconds()

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
			start := time.Now()
			checkStatus := checkpoint.After(message, bizHandleStatus.getErr())
			latency := time.Now().Sub(start).Seconds()
			if handledDeferFunc := checkStatus.GetHandledDefer(); handledDeferFunc != nil {
				defer handledDeferFunc()
			}
			checkpoint.metrics.CheckLatency.Observe(latency)
			if !checkStatus.IsPassed() {
				checkpoint.metrics.CheckRejected.Inc()
				continue
			}
			checkpoint.metrics.CheckPassed.Inc()
			if checkHandler := c.parseHandler(checkType, message); checkHandler != nil {
				if checkHandler.Decide(message.ConsumerMessage, checkStatus) {
					// return if check handle succeeded
					return
				}
			}

		}
	}

	// here means to let application client Ack/Nack message
	return
}

func (c *consumeListener) collectCheckers(enables *internal.StatusEnables, checkpointMap map[internal.CheckType]*checker.Checkpoint) map[internal.CheckType]*wrappedCheckpoint {
	checkers := make(map[internal.CheckType]*wrappedCheckpoint)
	if enables.RerouteEnable {
		c.tryLoadConfiguredChecker(&checkers, checker.CheckTypePreReroute, checkpointMap)
		c.tryLoadConfiguredChecker(&checkers, checker.CheckTypePostReroute, checkpointMap)
	}
	if enables.PendingEnable {
		c.tryLoadConfiguredChecker(&checkers, checker.CheckTypePrePending, checkpointMap)
		c.tryLoadConfiguredChecker(&checkers, checker.CheckTypePostPending, checkpointMap)
	}
	if enables.BlockingEnable {
		c.tryLoadConfiguredChecker(&checkers, checker.CheckTypePreBlocking, checkpointMap)
		c.tryLoadConfiguredChecker(&checkers, checker.CheckTypePostBlocking, checkpointMap)
	}
	if enables.RetryingEnable {
		c.tryLoadConfiguredChecker(&checkers, checker.CheckTypePreRetrying, checkpointMap)
		c.tryLoadConfiguredChecker(&checkers, checker.CheckTypePostRetrying, checkpointMap)
	}
	if enables.DeadEnable {
		c.tryLoadConfiguredChecker(&checkers, checker.CheckTypePreDead, checkpointMap)
		c.tryLoadConfiguredChecker(&checkers, checker.CheckTypePostDead, checkpointMap)
	}
	if enables.DiscardEnable {
		c.tryLoadConfiguredChecker(&checkers, checker.CheckTypePreDiscard, checkpointMap)
		c.tryLoadConfiguredChecker(&checkers, checker.CheckTypePostDiscard, checkpointMap)
	}
	if enables.UpgradeEnable {
		c.tryLoadConfiguredChecker(&checkers, checker.CheckTypePreUpgrade, checkpointMap)
		c.tryLoadConfiguredChecker(&checkers, checker.CheckTypePostUpgrade, checkpointMap)
	}
	if enables.DegradeEnable {
		c.tryLoadConfiguredChecker(&checkers, checker.CheckTypePreDegrade, checkpointMap)
		c.tryLoadConfiguredChecker(&checkers, checker.CheckTypePostDegrade, checkpointMap)
	}
	return checkers
}

func (c *consumeListener) tryLoadConfiguredChecker(checkers *map[internal.CheckType]*wrappedCheckpoint, checkType internal.CheckType, checkpointMap map[internal.CheckType]*checker.Checkpoint) {
	if ckp, ok := checkpointMap[checkType]; ok {
		metrics := c.client.metricsProvider.GetListenerTypedCheckMetrics(c.logTopics, c.logLevels, checkType)
		(*checkers)[checkType] = newWrappedCheckpoint(ckp, metrics)
	}
}

func (c *consumeListener) parseHandler(checkType internal.CheckType, message ConsumerMessage) internalDecider {
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
		return c.generalHandlers.doneHandler.Decide(msg, cheStatus)
	case message.GotoPending:
		return c.enables.PendingEnable && c.levelHandlers[l].pendingHandler.Decide(msg, cheStatus)
	case message.GotoBlocking:
		return c.enables.BlockingEnable && c.levelHandlers[l].blockingHandler.Decide(msg, cheStatus)
	case message.GotoRetrying:
		return c.enables.RetryingEnable && c.levelHandlers[l].retryingHandler.Decide(msg, cheStatus)
	case message.GotoDead:
		return c.enables.DeadEnable && c.generalHandlers.deadHandler.Decide(msg, cheStatus)
	case message.GotoDiscard:
		return c.enables.DiscardEnable && c.generalHandlers.discardHandler.Decide(msg, cheStatus)
	case message.GotoUpgrade:
		return c.enables.UpgradeEnable && c.levelHandlers[l].upgradeHandler.Decide(msg, cheStatus)
	case message.GotoDegrade:
		return c.enables.DegradeEnable && c.levelHandlers[l].degradeHandler.Decide(msg, cheStatus)
	default:
		c.logger.Warnf("invalid msg goto action: %v", messageGoto)
	}
	return false
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

func (c *consumeListener) Close() {
	c.closeListenerOnce.Do(func() {
		if c.leveledConsumer != nil {
			c.leveledConsumer.Close()
		}
		if c.multiLeveledConsumer != nil {
			for _, con := range c.multiLeveledConsumer.levelConsumers {
				con.Close()
			}
		}
		for _, chk := range c.checkers {
			chk.Close()
		}
		if c.generalHandlers != nil {
			c.generalHandlers.Close()
		}
		for _, hds := range c.levelHandlers {
			hds.Close()
		}
		c.logger.Info("closed consumer listener")
		c.metrics.ListenersOpened.Dec()
	})

}

// ------ helper ------

type wrappedCheckpoint struct {
	*checker.Checkpoint
	metrics *internal.TypedCheckMetrics
}

func newWrappedCheckpoint(ckp *checker.Checkpoint, metrics *internal.TypedCheckMetrics) *wrappedCheckpoint {
	wrappedCkp := &wrappedCheckpoint{Checkpoint: ckp, metrics: metrics}
	wrappedCkp.metrics.CheckersOpened.Inc()
	return wrappedCkp
}

func (c *wrappedCheckpoint) Close() {
	c.metrics.CheckersOpened.Dec()
}
