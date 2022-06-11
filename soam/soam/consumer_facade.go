package soam

import (
	"fmt"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
)

type consumer struct {
	//pulsar.Consumer
	client              *client
	config              ComsumerConfig
	logger              log.Logger
	messageCh           chan pulsar.ConsumerMessage // channel used to deliver message to application
	multiStatusConsumer *mixStatusConsumer
	checkers            *consumeCheckers
	handlers            *consumeHandlers
}

func newConsumer(c *client, config ComsumerConfig, handler PremiumHandler, checkpointMap map[CheckType]*checkpoint) (*consumer, error) {
	consumer := &consumer{
		client:    c,
		config:    config,
		messageCh: make(chan pulsar.ConsumerMessage, 10),
		logger:    c.logger,
	}
	// initialize status multiStatusConsumer
	if consumers, err := newMixStatusConsumer(consumer, config); err != nil {
		return nil, err
	} else {
		consumer.multiStatusConsumer = consumers
	}
	// initialize checkers
	if checkers, err := NewConsumeCheckers(consumer, config, checkpointMap); err != nil {
		return nil, err
	} else {
		consumer.checkers = checkers
	}
	// initialize handlers
	if handlers, err := newConsumeHandlers(consumer, config); err != nil {
		return nil, err
	} else {
		consumer.handlers = handlers
	}
	// start to consume
	go consumer.start(handler)
	return consumer, nil
}

func (c *consumer) start(handler PremiumHandler) {
	defer c.Close()
	for msg := range c.messageCh {
		c.consume(handler, msg)
	}
	fmt.Println("end consumer start")
}

func (c *consumer) consume(handler PremiumHandler, message pulsar.ConsumerMessage) {
	// pre-check to route
	if routed := c.preCheckToHandleInTurn(message, defaultPreCheckTypesInTurn...); routed {
		return
	}
	// do biz handle
	handleResult := handler(message.Message)
	// post-check to route - for obvious goto action
	if handleResult.getGotoAction() != "" {
		if ok := c.handleMessageGotoAction(message, handleResult.getGotoAction()); ok {
			return
		}
	}
	// post-check to route - for obvious checkers or configured checkers
	postCheckTypesInTurn := DefaultPostCheckTypesInTurn
	if len(handleResult.getCheckTypes()) > 0 {
		postCheckTypesInTurn = handleResult.getCheckTypes()
	}
	if routed := c.postCheckToHandleInTurn(message, handleResult.getErr(), postCheckTypesInTurn...); routed {
		return
	}
	// here means to let application client Ack/Nack message
	return
}

func (c *consumer) preCheckToHandleInTurn(msg pulsar.ConsumerMessage, checkTypes ...CheckType) (handled bool) {
	for _, checkType := range checkTypes {
		switch checkType {
		case CheckTypePreDiscard:
			if ok := c.internalCheckToHandleMsg(msg, c.config.DiscardEnable, c.checkers.PreDiscardChecker, c.handlers.discardHandler); ok {
				return true
			}
		case CheckTypePreDead:
			if ok := c.internalCheckToHandleMsg(msg, c.config.DeadEnable, c.checkers.PreDeadChecker, c.handlers.deadHandler); ok {
				return true
			}
		case CheckTypePreReroute:
			if ok := c.internalCheckToRerouteMsg(msg, c.config.RerouteEnable, c.checkers.PreRerouteChecker, c.handlers.rerouteHandler); ok {
				return true
			}
		case CheckTypePreUpgrade:
			if ok := c.internalCheckToHandleMsg(msg, c.config.UpgradeEnable, c.checkers.PreUpgradeChecker, c.handlers.upgradeHandler); ok {
				return true
			}
		case CheckTypePreDegrade:
			if ok := c.internalCheckToHandleMsg(msg, c.config.DegradeEnable, c.checkers.PreDegradeChecker, c.handlers.degradeHandler); ok {
				return true
			}
		case CheckTypePreBlocking:
			if ok := c.internalCheckToHandleMsg(msg, c.config.BlockingEnable, c.checkers.PreBlockingChecker, c.handlers.blockingHandler); ok {
				return true
			}
		case CheckTypePrePending:
			if ok := c.internalCheckToHandleMsg(msg, c.config.PendingEnable, c.checkers.PrePendingChecker, c.handlers.pendingHandler); ok {
				return true
			}
		case CheckTypePreRetrying:
			if ok := c.internalCheckToHandleMsg(msg, c.config.RetryingEnable, c.checkers.PreRetryingChecker, c.handlers.retryingHandler); ok {
				return true
			}
		}
	}
	return false
}

func (c *consumer) handleMessageGotoAction(message pulsar.ConsumerMessage, messageGoto messageGotoAction) (routed bool) {
	switch messageGoto {
	case MessageGotoDone:
		return c.handlers.doneHandler.Handle(message)
	case MessageGotoPending:
		return c.config.PendingEnable && c.handlers.pendingHandler.Handle(message)
	case MessageGotoBlocking:
		return c.config.BlockingEnable && c.handlers.blockingHandler.Handle(message)
	case MessageGotoRetrying:
		return c.config.RetryEnable && c.handlers.retryingHandler.Handle(message)
	case MessageGotoDead:
		return c.config.DeadEnable && c.handlers.deadHandler.Handle(message)
	case MessageGotoDiscard:
		return c.config.DiscardEnable && c.handlers.discardHandler.Handle(message)
	case MessageGotoUpgrade:
		return c.config.UpgradeEnable && c.handlers.upgradeHandler.Handle(message)
	case MessageGotoDegrade:
		return c.config.DegradeEnable && c.handlers.degradeHandler.Handle(message)
	default:
		c.logger.Warnf("invalid message goto action: %v", messageGoto)
	}
	return false
}

func (c *consumer) postCheckToHandleInTurn(msg pulsar.ConsumerMessage, err error, checkTypes ...CheckType) (routed bool) {
	for _, checkType := range checkTypes {
		switch checkType {
		case CheckTypePostDiscard:
			if ok := c.internalCheckToHandleMsgWithErr(msg, err, c.config.DiscardEnable, c.checkers.PostDiscardChecker, c.handlers.discardHandler); ok {
				return true
			}
		case CheckTypePostDead:
			if ok := c.internalCheckToHandleMsgWithErr(msg, err, c.config.DeadEnable, c.checkers.PostDeadChecker, c.handlers.deadHandler); ok {
				return true
			}
		case CheckTypePostReroute:
			if ok := c.internalCheckToRerouteMsgWithErr(msg, err, c.config.RerouteEnable, c.checkers.PostReRouteChecker, c.handlers.rerouteHandler); ok {
				return true
			}
		case CheckTypePostUpgrade:
			if ok := c.internalCheckToHandleMsgWithErr(msg, err, c.config.UpgradeEnable, c.checkers.PostUpgradeChecker, c.handlers.upgradeHandler); ok {
				return true
			}
		case CheckTypePostDegrade:
			if ok := c.internalCheckToHandleMsgWithErr(msg, err, c.config.DegradeEnable, c.checkers.PostDegradeChecker, c.handlers.degradeHandler); ok {
				return true
			}
		case CheckTypePostBlocking:
			if ok := c.internalCheckToHandleMsgWithErr(msg, err, c.config.BlockingEnable, c.checkers.PostBlockingChecker, c.handlers.blockingHandler); ok {
				return true
			}
		case CheckTypePostPending:
			if ok := c.internalCheckToHandleMsgWithErr(msg, err, c.config.PendingEnable, c.checkers.PostPendingChecker, c.handlers.pendingHandler); ok {
				return true
			}
		case CheckTypePostRetrying:
			if ok := c.internalCheckToHandleMsgWithErr(msg, err, c.config.RetryingEnable, c.checkers.PostRetryingChecker, c.handlers.retryingHandler); ok {
				return true
			}
		}
	}
	return false
}

func (c *consumer) internalHandleMsg(msg pulsar.ConsumerMessage, statusEnable bool, handler internalHandler) (routed bool) {
	return statusEnable && handler.Handle(msg)
}

func (c *consumer) internalCheckToHandleMsg(msg pulsar.ConsumerMessage, statusEnable bool,
	checker preStatusChecker, handler internalHandler) (routed bool) {
	return statusEnable && checker(msg) && handler.Handle(msg)
}

func (c *consumer) internalCheckToHandleMsgWithErr(msg pulsar.ConsumerMessage, err error, statusEnable bool,
	checker postStatusChecker, handler internalHandler) (routed bool) {
	return statusEnable && checker(msg, err) && handler.Handle(msg)
}

func (c *consumer) internalCheckToRerouteMsg(msg pulsar.ConsumerMessage, statusEnable bool,
	checker preRerouteChecker, handler internalRerouteHandler) (routed bool) {
	if statusEnable {
		if topic := checker(msg); topic != "" {
			return handler.Handle(msg, topic)
		}
	}
	return false
}

func (c *consumer) internalCheckToRerouteMsgWithErr(msg pulsar.ConsumerMessage, err error, statusEnable bool,
	checker postRerouteChecker, handler internalRerouteHandler) (routed bool) {
	if statusEnable {
		if topic := checker(msg, err); topic != "" {
			return handler.Handle(msg, topic)
		}
	}
	return false
}
func (c *consumer) Close() {
	/*c.closeOnce.Do(func() {
		var wg sync.WaitGroup
		wg.Add(len(c.multiStatusConsumer))
		for _, con := range c.multiStatusConsumer {
			go func(consumer Consumer) {
				defer wg.Done()
				consumer.Close()
			}(con)
		}
		wg.Wait()
		close(c.closeCh)
		c.client.handlers.Del(c)
		c.dlq.close()
		c.rlq.close()
	})*/
}
