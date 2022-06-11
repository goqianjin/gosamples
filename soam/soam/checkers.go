package soam

import (
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
)

type consumeCheckers struct {
	consumer *consumer
	config   ComsumerConfig
	log      log.Logger

	preBlockingChecker  preStatusChecker  //
	postBlockingChecker postStatusChecker //
	prePendingChecker   preStatusChecker  //
	postPendingChecker  postStatusChecker //
	preRetryingChecker  preStatusChecker  //
	postRetryingChecker postStatusChecker //
	preDeadChecker      preStatusChecker  //
	postDeadChecker     postStatusChecker //
	preDiscardChecker   preStatusChecker  //
	postDiscardChecker  postStatusChecker //

	preUpgradeChecker  preStatusChecker  //
	postUpgradeChecker postStatusChecker //
	preDegradeChecker  preStatusChecker  //
	postDegradeChecker postStatusChecker //

	preReRouteChecker  preReRouterChecker  //
	postReRouteChecker postReRouterChecker //
}

func newConsumeCheckers(consumer *consumer, config ComsumerConfig, checkpointMap map[CheckType]*checkpoint) (*consumeCheckers, error) {
	checkers := &consumeCheckers{
		consumer: consumer,
		config:   config,
		log:      consumer.log,
	}
	if config.ReRouteEnable {
		checkers.preReRouteChecker = checkers.parsePreRerouteChecker(CheckTypePreReRoute, checkpointMap)
		checkers.postReRouteChecker = checkers.parsePostRerouteChecker(CheckTypePostReRoute, checkpointMap)
	}
	if config.PendingEnable {
		checkers.prePendingChecker = checkers.parsePreStatusChecker(CheckTypePrePending, checkpointMap)
		checkers.postPendingChecker = checkers.parsePostStatusChecker(CheckTypePostPending, checkpointMap)
	}
	if config.BlockingEnable {
		checkers.preBlockingChecker = checkers.parsePreStatusChecker(CheckTypePreBlocking, checkpointMap)
		checkers.postBlockingChecker = checkers.parsePostStatusChecker(CheckTypePostBlocking, checkpointMap)
	}
	if config.RetryingEnable {
		checkers.preRetryingChecker = checkers.parsePreStatusChecker(CheckTypePreRetrying, checkpointMap)
		checkers.postRetryingChecker = checkers.parsePostStatusChecker(CheckTypePostRetrying, checkpointMap)
	}
	if config.DeadEnable {
		checkers.preDeadChecker = checkers.parsePreStatusChecker(CheckTypePreDead, checkpointMap)
		checkers.postDeadChecker = checkers.parsePostStatusChecker(CheckTypePostDead, checkpointMap)
	}
	if config.DiscardEnable {
		checkers.preDiscardChecker = checkers.parsePreStatusChecker(CheckTypePreDiscard, checkpointMap)
		checkers.postDiscardChecker = checkers.parsePostStatusChecker(CheckTypePostDiscard, checkpointMap)
	}
	if config.UpgradeEnable {
		checkers.preUpgradeChecker = checkers.parsePreStatusChecker(CheckTypePreUpgrade, checkpointMap)
		checkers.postUpgradeChecker = checkers.parsePostStatusChecker(CheckTypePostUpgrade, checkpointMap)
	}
	if config.DegradeEnable {
		checkers.preDegradeChecker = checkers.parsePreStatusChecker(CheckTypePreDegrade, checkpointMap)
		checkers.postDegradeChecker = checkers.parsePostStatusChecker(CheckTypePostDegrade, checkpointMap)
	}
	return checkers, nil
}

func (ch *consumeCheckers) parsePreStatusChecker(checkType CheckType, checkpointMap map[CheckType]*checkpoint) preStatusChecker {
	if ckp, ok := checkpointMap[checkType]; ok {
		if ckp.preStatusChecker != nil {
			return ckp.preStatusChecker
		}
	}
	return nilPreStatusChecker
}

func (ch *consumeCheckers) parsePostStatusChecker(checkType CheckType, checkpointMap map[CheckType]*checkpoint) postStatusChecker {
	if ckp, ok := checkpointMap[checkType]; ok {
		if ckp.postStatusChecker != nil {
			return ckp.postStatusChecker
		}
	}
	return nilPostStatusChecker
}

func (ch *consumeCheckers) parsePreRerouteChecker(checkType CheckType, checkpointMap map[CheckType]*checkpoint) preReRouterChecker {
	if ckp, ok := checkpointMap[checkType]; ok {
		if ckp.preReRouterChecker != nil {
			return ckp.preReRouterChecker
		}
	}
	return nilPreReRouterChecker
}

func (ch *consumeCheckers) parsePostRerouteChecker(checkType CheckType, checkpointMap map[CheckType]*checkpoint) postReRouterChecker {
	if ckp, ok := checkpointMap[checkType]; ok {
		if ckp.postReRouterChecker != nil {
			return ckp.postReRouterChecker
		}
	}
	return nilPostReRouterChecker
}

func (ch *consumeCheckers) tryPreCheckToHandleInTurn(msg pulsar.ConsumerMessage, checkTypes ...CheckType) (handled bool) {
	for _, checkType := range checkTypes {
		switch checkType {
		case CheckTypePreDiscard:
			if ok := ch.internalCheckToRouteMsg(msg, ch.config.DiscardEnable, ch.preDiscardChecker, ch.consumer.handlers.discardHandler); ok {
				return true
			}
		case CheckTypePreDead:
			if ok := ch.internalCheckToRouteMsg(msg, ch.config.DeadEnable, ch.preDeadChecker, ch.consumer.handlers.deadHandler); ok {
				return true
			}
		case CheckTypePreReRoute:
			if ok := ch.internalCheckToCustomRouteMsg(msg, ch.config.ReRouteEnable, ch.preReRouteChecker, ch.consumer.handlers.rerouteHandler); ok {
				return true
			}
		case CheckTypePreUpgrade:
			if ok := ch.internalCheckToRouteMsg(msg, ch.config.UpgradeEnable, ch.preUpgradeChecker, ch.consumer.handlers.upgradeHandler); ok {
				return true
			}
		case CheckTypePreDegrade:
			if ok := ch.internalCheckToRouteMsg(msg, ch.config.DegradeEnable, ch.preDegradeChecker, ch.consumer.handlers.degradeHandler); ok {
				return true
			}
		case CheckTypePreBlocking:
			if ok := ch.internalCheckToRouteMsg(msg, ch.config.BlockingEnable, ch.preBlockingChecker, ch.consumer.handlers.blockingHandler); ok {
				return true
			}
		case CheckTypePrePending:
			if ok := ch.internalCheckToRouteMsg(msg, ch.config.PendingEnable, ch.prePendingChecker, ch.consumer.handlers.pendingHandler); ok {
				return true
			}
		case CheckTypePreRetrying:
			if ok := ch.internalCheckToRouteMsg(msg, ch.config.RetryingEnable, ch.preRetryingChecker, ch.consumer.handlers.retryingHandler); ok {
				return true
			}
		}
	}
	return false
}

func (ch *consumeCheckers) tryPostCheckInTurn(msg pulsar.ConsumerMessage, err error, checkTypes ...CheckType) (routed bool) {
	for _, checkType := range checkTypes {
		switch checkType {
		case CheckTypePostDiscard:
			if ok := ch.internalCheckToRouteMsgWithErr(msg, err, ch.config.DiscardEnable, ch.postDiscardChecker, ch.consumer.handlers.discardHandler); ok {
				return true
			}
		case CheckTypePostDead:
			if ok := ch.internalCheckToRouteMsgWithErr(msg, err, ch.config.DeadEnable, ch.postDeadChecker, ch.consumer.handlers.deadHandler); ok {
				return true
			}
		case CheckTypePostReRoute:
			if ok := ch.internalCheckToCustomRouteMsgWithErr(msg, err, ch.config.ReRouteEnable, ch.postReRouteChecker, ch.consumer.handlers.rerouteHandler); ok {
				return true
			}
		case CheckTypePostUpgrade:
			if ok := ch.internalCheckToRouteMsgWithErr(msg, err, ch.config.UpgradeEnable, ch.postUpgradeChecker, ch.consumer.handlers.upgradeHandler); ok {
				return true
			}
		case CheckTypePostDegrade:
			if ok := ch.internalCheckToRouteMsgWithErr(msg, err, ch.config.DegradeEnable, ch.postDegradeChecker, ch.consumer.handlers.degradeHandler); ok {
				return true
			}
		case CheckTypePostBlocking:
			if ok := ch.internalCheckToRouteMsgWithErr(msg, err, ch.config.BlockingEnable, ch.postBlockingChecker, ch.consumer.handlers.blockingHandler); ok {
				return true
			}
		case CheckTypePostPending:
			if ok := ch.internalCheckToRouteMsgWithErr(msg, err, ch.config.PendingEnable, ch.postPendingChecker, ch.consumer.handlers.pendingHandler); ok {
				return true
			}
		case CheckTypePostRetrying:
			if ok := ch.internalCheckToRouteMsgWithErr(msg, err, ch.config.RetryingEnable, ch.postRetryingChecker, ch.consumer.handlers.retryingHandler); ok {
				return true
			}
		}
	}
	return false
}

func (ch *consumeCheckers) internalCheckToRouteMsg(msg pulsar.ConsumerMessage, statusEnable bool,
	checker preStatusChecker, handler internalHandler) (routed bool) {
	return statusEnable && checker(msg) && handler.Handle(msg)
}

func (ch *consumeCheckers) internalCheckToRouteMsgWithErr(msg pulsar.ConsumerMessage, err error, statusEnable bool,
	checker postStatusChecker, handler internalHandler) (routed bool) {
	return statusEnable && checker(msg, err) && handler.Handle(msg)
}

func (ch *consumeCheckers) internalCheckToCustomRouteMsg(msg pulsar.ConsumerMessage, statusEnable bool,
	checker preReRouterChecker, handler internalRerouteHandler) (routed bool) {
	if statusEnable {
		if topic := checker(msg); topic != "" {
			return handler.Handle(msg, topic)
		}
	}
	return false
}

func (ch *consumeCheckers) internalCheckToCustomRouteMsgWithErr(msg pulsar.ConsumerMessage, err error, statusEnable bool,
	checker postReRouterChecker, handler internalRerouteHandler) (routed bool) {
	if statusEnable {
		if topic := checker(msg, err); topic != "" {
			return handler.Handle(msg, topic)
		}
	}
	return false
}
