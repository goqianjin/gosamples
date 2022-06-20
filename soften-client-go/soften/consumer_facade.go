package soften

import (
	"fmt"
	"strconv"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
	"github.com/shenqianjin/soften-client-go/soften/checker"
	"github.com/shenqianjin/soften-client-go/soften/config"
	"github.com/shenqianjin/soften-client-go/soften/internal"
	"github.com/shenqianjin/soften-client-go/soften/message"
	"github.com/shenqianjin/soften-client-go/soften/topic"
)

type consumeFacade struct {
	// pulsar.Consumer
	client          *client
	logger          log.Logger
	messageCh       chan ConsumerMessage // channel used to deliver message to application
	enables         *internal.StatusEnables
	checkers        *consumeCheckers
	generalHandlers *generalConsumeHandlers
	levelHandlers   map[internal.TopicLevel]*leveledConsumeHandlers
}

func (c *consumeFacade) collectEnables(conf *config.ConsumerConfig) *internal.StatusEnables {
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
func (c *consumeFacade) formatGeneralHandlersOptions(topic string, config *config.ConsumerConfig) generalConsumeHandlerOptions {
	options := generalConsumeHandlerOptions{
		Topic:         topic,
		DiscardEnable: config.BlockingEnable,
		DeadEnable:    config.RetryEnable,
		RerouteEnable: config.RerouteEnable,
	}
	return options
}

func (c *consumeFacade) formatLeveledHandlersOptions(levelTopic string, config *config.ConsumerConfig) leveledConsumeHandlerOptions {
	options := leveledConsumeHandlerOptions{
		Topic:             levelTopic,
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

func newMultiStatusConsumeFacade(c *client, conf config.ConsumerConfig, handler PremiumHandler, checkpointMap map[internal.CheckType]*internal.Checkpoint) (*consumeFacade, error) {
	// convert config.ConsumerConfig to config.MultiLevelConsumerConfig
	multiLevelConf := config.MultiLevelConsumerConfig{
		ConsumerConfig: &conf,
		Levels:         []internal.TopicLevel{conf.Level},
	}
	// forward
	return newMultiLevelConsumeFacade(c, multiLevelConf, handler, checkpointMap)
}

func newMultiLevelConsumeFacade(cli *client, conf config.MultiLevelConsumerConfig, handler PremiumHandler, checkpointMap map[internal.CheckType]*internal.Checkpoint) (*consumeFacade, error) {
	logTopic := conf.Topics[0]
	if len(conf.Topics) > 1 {
		logTopic = logTopic + "+" + strconv.Itoa(len(conf.Topics)-1)
	}
	facade := &consumeFacade{
		client:    cli,
		messageCh: make(chan ConsumerMessage, 10),
		logger:    cli.logger.SubLogger(log.Fields{"Topic": logTopic}),
	}
	// collect enables
	facade.enables = facade.collectEnables(conf.ConsumerConfig)
	// initialize checkers
	if checkers, err := newConsumeCheckers(facade.logger, facade.enables, checkpointMap); err != nil {
		return nil, err
	} else {
		facade.checkers = checkers
	}
	// initialize general handlers
	generalHdOptions := facade.formatGeneralHandlersOptions(conf.Topics[0], conf.ConsumerConfig)
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
		options := facade.formatLeveledHandlersOptions(conf.Topics[0]+suffix, conf.ConsumerConfig)
		if handlers, err := newLeveledConsumeHandlers(cli, options, facade.generalHandlers.deadHandler); err != nil {
			return nil, err
		} else {
			facade.levelHandlers[level] = handlers
		}
	}
	// initialize status multiStatusConsumer
	if len(conf.Levels) == 1 {
		level := conf.Levels[0]
		if consumer, err := newMultiStatusConsumer(facade.logger, cli, level, conf.ConsumerConfig, facade.messageCh, facade.levelHandlers[level]); err != nil {
			return nil, err
		} else {
			facade.logger.Info("newMultiStatusConsumer done", consumer)
		}
	} else {
		if consumer, err := newMultiLevelConsumer(facade.logger, cli, conf, facade.messageCh, facade.levelHandlers); err != nil {
			return nil, err
		} else {
			facade.logger.Info("newMultiLevelConsumer done", consumer)
		}
	}
	// start to consume
	go facade.startInParallel(handler, conf.Concurrency)
	return facade, nil
}

func (c *consumeFacade) startInParallel(handler PremiumHandler, concurrency uint) {
	defer c.Close()
	concurrencyChan := make(chan bool, 1)
	for msg := range c.messageCh {
		concurrencyChan <- true
		go func(msg ConsumerMessage) {
			c.consume(handler, msg)
			<-concurrencyChan
		}(msg)
	}
	fmt.Println("end multiStatusConsumeFacade start")
}

func (c *consumeFacade) consume(handler PremiumHandler, message ConsumerMessage) {
	// pre-check to route
	if routed := c.preCheckToHandleInTurn(message, checker.PreCheckTypes()...); routed {
		return
	}
	// do biz handle
	handleResult := handler(message)
	// post-check to route - for obvious goto action
	if handleResult.getGotoAction() != "" {
		if ok := c.handleMessageGotoAction(message, handleResult.getGotoAction()); ok {
			return
		}
	}
	// post-check to route - for obvious checkers or configured checkers
	postCheckTypesInTurn := checker.DefaultPostCheckTypes()
	if len(handleResult.getCheckTypes()) > 0 {
		postCheckTypesInTurn = handleResult.getCheckTypes()
	}
	if routed := c.postCheckToHandleInTurn(message, handleResult.getErr(), postCheckTypesInTurn...); routed {
		return
	}
	// here means to let application client Ack/Nack message
	return
}

func (c *consumeFacade) preCheckToHandleInTurn(message ConsumerMessage, checkTypes ...internal.CheckType) (handled bool) {
	l := message.Level()
	msg := message.ConsumerMessage
	for _, checkType := range checkTypes {
		switch checkType {
		case checker.CheckTypePreDiscard:
			if ok := c.internalCheckToHandleMsg(msg, c.enables.DiscardEnable, c.checkers.PreDiscardChecker, c.generalHandlers.discardHandler); ok {
				return true
			}
		case checker.CheckTypePreDead:
			if ok := c.internalCheckToHandleMsg(msg, c.enables.DeadEnable, c.checkers.PreDeadChecker, c.generalHandlers.deadHandler); ok {
				return true
			}
		case checker.CheckTypePreReroute:
			if ok := c.internalCheckToRerouteMsg(msg, c.enables.RerouteEnable, c.checkers.PreRerouteChecker, c.generalHandlers.rerouteHandler); ok {
				return true
			}
		case checker.CheckTypePreUpgrade:
			if ok := c.internalCheckToHandleMsg(msg, c.enables.UpgradeEnable, c.checkers.PreUpgradeChecker, c.levelHandlers[l].upgradeHandler); ok {
				return true
			}
		case checker.CheckTypePreDegrade:
			if ok := c.internalCheckToHandleMsg(msg, c.enables.DegradeEnable, c.checkers.PreDegradeChecker, c.levelHandlers[l].degradeHandler); ok {
				return true
			}
		case checker.CheckTypePreBlocking:
			if ok := c.internalCheckToHandleMsg(msg, c.enables.BlockingEnable, c.checkers.PreBlockingChecker, c.levelHandlers[l].blockingHandler); ok {
				return true
			}
		case checker.CheckTypePrePending:
			if ok := c.internalCheckToHandleMsg(msg, c.enables.PendingEnable, c.checkers.PrePendingChecker, c.levelHandlers[l].pendingHandler); ok {
				return true
			}
		case checker.CheckTypePreRetrying:
			if ok := c.internalCheckToHandleMsg(msg, c.enables.RetryingEnable, c.checkers.PreRetryingChecker, c.levelHandlers[l].retryingHandler); ok {
				return true
			}
		}
	}
	return false
}

func (c *consumeFacade) handleMessageGotoAction(consumerMessage ConsumerMessage, messageGoto internal.MessageGoto) (routed bool) {
	l := consumerMessage.Level()
	msg := consumerMessage.ConsumerMessage
	switch messageGoto {
	case message.GotoDone:
		return c.generalHandlers.doneHandler.Handle(msg)
	case message.GotoPending:
		return c.enables.PendingEnable && c.levelHandlers[l].pendingHandler.Handle(msg)
	case message.GotoBlocking:
		return c.enables.BlockingEnable && c.levelHandlers[l].blockingHandler.Handle(msg)
	case message.GotoRetrying:
		return c.enables.RetryingEnable && c.levelHandlers[l].retryingHandler.Handle(msg)
	case message.GotoDead:
		return c.enables.DeadEnable && c.generalHandlers.deadHandler.Handle(msg)
	case message.GotoDiscard:
		return c.enables.DiscardEnable && c.generalHandlers.discardHandler.Handle(msg)
	case message.GotoUpgrade:
		return c.enables.UpgradeEnable && c.levelHandlers[l].upgradeHandler.Handle(msg)
	case message.GotoDegrade:
		return c.enables.DegradeEnable && c.levelHandlers[l].degradeHandler.Handle(msg)
	default:
		c.logger.Warnf("invalid msg goto action: %v", messageGoto)
	}
	return false
}

func (c *consumeFacade) postCheckToHandleInTurn(message ConsumerMessage, err error, checkTypes ...internal.CheckType) (routed bool) {
	l := message.Level()
	msg := message.ConsumerMessage
	for _, checkType := range checkTypes {
		switch checkType {
		case checker.CheckTypePostDiscard:
			if ok := c.internalCheckToHandleMsgWithErr(msg, err, c.enables.DiscardEnable, c.checkers.PostDiscardChecker, c.generalHandlers.discardHandler); ok {
				return true
			}
		case checker.CheckTypePostDead:
			if ok := c.internalCheckToHandleMsgWithErr(msg, err, c.enables.DeadEnable, c.checkers.PostDeadChecker, c.generalHandlers.deadHandler); ok {
				return true
			}
		case checker.CheckTypePostReroute:
			if ok := c.internalCheckToRerouteMsgWithErr(msg, err, c.enables.RerouteEnable, c.checkers.PostReRouteChecker, c.generalHandlers.rerouteHandler); ok {
				return true
			}
		case checker.CheckTypePostUpgrade:
			if ok := c.internalCheckToHandleMsgWithErr(msg, err, c.enables.UpgradeEnable, c.checkers.PostUpgradeChecker, c.levelHandlers[l].upgradeHandler); ok {
				return true
			}
		case checker.CheckTypePostDegrade:
			if ok := c.internalCheckToHandleMsgWithErr(msg, err, c.enables.DegradeEnable, c.checkers.PostDegradeChecker, c.levelHandlers[l].degradeHandler); ok {
				return true
			}
		case checker.CheckTypePostBlocking:
			if ok := c.internalCheckToHandleMsgWithErr(msg, err, c.enables.BlockingEnable, c.checkers.PostBlockingChecker, c.levelHandlers[l].blockingHandler); ok {
				return true
			}
		case checker.CheckTypePostPending:
			if ok := c.internalCheckToHandleMsgWithErr(msg, err, c.enables.PendingEnable, c.checkers.PostPendingChecker, c.levelHandlers[l].pendingHandler); ok {
				return true
			}
		case checker.CheckTypePostRetrying:
			if ok := c.internalCheckToHandleMsgWithErr(msg, err, c.enables.RetryingEnable, c.checkers.PostRetryingChecker, c.levelHandlers[l].retryingHandler); ok {
				return true
			}
		}
	}
	return false
}

func (c *consumeFacade) internalHandleMsg(msg pulsar.ConsumerMessage, statusEnable bool, handler internal.Handler) (routed bool) {
	return statusEnable && handler.Handle(msg)
}

func (c *consumeFacade) internalCheckToHandleMsg(msg pulsar.ConsumerMessage, statusEnable bool,
	checker internal.PreStatusChecker, handler internal.Handler) (routed bool) {
	return statusEnable && checker(msg) && handler.Handle(msg)
}

func (c *consumeFacade) internalCheckToHandleMsgWithErr(msg pulsar.ConsumerMessage, err error, statusEnable bool,
	checker internal.PostStatusChecker, handler internal.Handler) (routed bool) {
	return statusEnable && checker(msg, err) && handler.Handle(msg)
}

func (c *consumeFacade) internalCheckToRerouteMsg(msg pulsar.ConsumerMessage, statusEnable bool,
	checker internal.PreRerouteChecker, handler internal.RerouteHandler) (routed bool) {
	if statusEnable {
		if tpc := checker(msg); tpc != "" {
			return handler.Handle(msg, tpc)
		}
	}
	return false
}

func (c *consumeFacade) internalCheckToRerouteMsgWithErr(msg pulsar.ConsumerMessage, err error, statusEnable bool,
	checker internal.PostRerouteChecker, handler internal.RerouteHandler) (routed bool) {
	if statusEnable {
		if tpc := checker(msg, err); tpc != "" {
			return handler.Handle(msg, tpc)
		}
	}
	return false
}
func (c *consumeFacade) Close() {
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
