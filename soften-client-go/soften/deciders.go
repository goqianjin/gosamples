package soften

import (
	"context"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/shenqianjin/soften-client-go/soften/checker"
	"github.com/shenqianjin/soften-client-go/soften/config"
	"github.com/shenqianjin/soften-client-go/soften/internal"
	"github.com/shenqianjin/soften-client-go/soften/message"
)

type internalProduceDecider interface {
	Decide(ctx context.Context, msg *pulsar.ProducerMessage,
		checkStatus checker.CheckStatus) (mid pulsar.MessageID, err error, decided bool)
	DecideAsync(ctx context.Context, msg *pulsar.ProducerMessage, checkStatus checker.CheckStatus,
		callback func(pulsar.MessageID, *pulsar.ProducerMessage, error)) (decided bool)
	close()
}

type internalDecider interface {
	Decide(msg pulsar.ConsumerMessage, checkStatus checker.CheckStatus) (success bool)
	close()
}

type produceDecidersOptions struct {
	BlockingEnable bool
	PendingEnable  bool
	RetryingEnable bool
	UpgradeEnable  bool
	DegradeEnable  bool
	DeadEnable     bool
	DiscardEnable  bool
	RouteEnable    bool
}

type produceDeciders map[internal.MessageGoto]internalProduceDecider

func newProduceDeciders(producer *producer, conf produceDecidersOptions) (produceDeciders, error) {
	deciders := make(produceDeciders)
	deciderOpt := routeDeciderOptions{connectInSyncEnable: true}
	if conf.DiscardEnable {
		if err := deciders.tryLoadDecider(producer, message.GotoDiscard, deciderOpt); err != nil {
			return nil, err
		}
	}
	if conf.DeadEnable {
		if err := deciders.tryLoadDecider(producer, message.GotoDead, deciderOpt); err != nil {
			return nil, err
		}
	}
	if conf.BlockingEnable {
		if err := deciders.tryLoadDecider(producer, message.GotoBlocking, deciderOpt); err != nil {
			return nil, err
		}
	}
	if conf.PendingEnable {
		if err := deciders.tryLoadDecider(producer, message.GotoPending, deciderOpt); err != nil {
			return nil, err
		}
	}
	if conf.RetryingEnable {
		if err := deciders.tryLoadDecider(producer, message.GotoRetrying, deciderOpt); err != nil {
			return nil, err
		}
	}
	if conf.RouteEnable {
		if err := deciders.tryLoadDecider(producer, internalGotoRoute, deciderOpt); err != nil {
			return nil, err
		}
	}
	if conf.UpgradeEnable {
		deciderOpt.upgradeLevel = producer.upgradeLevel
		if err := deciders.tryLoadDecider(producer, message.GotoUpgrade, deciderOpt); err != nil {
			return nil, err
		}
	}
	if conf.DegradeEnable {
		deciderOpt.degradeLevel = producer.degradeLevel
		if err := deciders.tryLoadDecider(producer, message.GotoDegrade, deciderOpt); err != nil {
			return nil, err
		}
	}
	return deciders, nil
}

func (deciders *produceDeciders) tryLoadDecider(producer *producer, msgGoto internal.MessageGoto, options routeDeciderOptions) error {
	decider, err := newRouteDecider(producer, msgGoto, &options)
	if err != nil {
		return err
	}
	(*deciders)[msgGoto] = decider
	return nil
}

// ------ general consume handlers ------

type generalConsumeDeciders struct {
	rerouteDecider internalDecider // ??????????????????: Reroute
	deadDecider    internalDecider // ???????????????
	doneDecider    internalDecider // ???????????????
	discardDecider internalDecider // ???????????????
}

type generalConsumeDeciderOptions struct {
	Topic         string                // Business Topic
	DiscardEnable bool                  // Blocking ????????????
	DeadEnable    bool                  // Pending ????????????
	RerouteEnable bool                  // Retrying ??????????????????
	Reroute       *config.ReroutePolicy // Reroute Policy
}

func newGeneralConsumeDeciders(client *client, listener *consumeListener, conf generalConsumeDeciderOptions) (*generalConsumeDeciders, error) {
	handlers := &generalConsumeDeciders{}
	doneDecider, err := newFinalStatusDecider(client, listener, message.GotoDone)
	if err != nil {
		return nil, err
	}
	handlers.doneDecider = doneDecider
	if conf.DiscardEnable {
		decider, err := newFinalStatusDecider(client, listener, message.GotoDiscard)
		if err != nil {
			return nil, err
		}
		handlers.discardDecider = decider
	}
	if conf.DeadEnable {
		suffix := message.StatusDead.TopicSuffix()
		deadOptions := deadDecideOptions{topic: conf.Topic + suffix}
		decider, err := newDeadDecider(client, listener, deadOptions)
		if err != nil {
			return nil, err
		}
		handlers.deadDecider = decider
	}
	if conf.RerouteEnable {
		decider, err := newRerouteDecider(client, listener, conf.Reroute)
		if err != nil {
			return nil, err
		}
		handlers.rerouteDecider = decider
	}
	return handlers, nil
}

func (hds generalConsumeDeciders) Close() {
	if hds.rerouteDecider != nil {
		hds.rerouteDecider.close()
	}
	if hds.deadDecider != nil {
		hds.deadDecider.close()
	}
	if hds.doneDecider != nil {
		hds.doneDecider.close()
	}
	if hds.discardDecider != nil {
		hds.discardDecider.close()
	}
}

// ------ leveled consume handlers ------

type leveledConsumeDeciders struct {
	blockingDecider internalDecider // ???????????????
	pendingDecider  internalDecider // ???????????????
	retryingDecider internalDecider // ???????????????
	upgradeDecider  internalDecider // ???????????????: ?????????NewReady
	degradeDecider  internalDecider // ???????????????: ?????????NewReady
}

type leveledConsumeDeciderOptions struct {
	Topic             string               // Business Topic
	Level             internal.TopicLevel  // level
	BlockingEnable    bool                 // Blocking ????????????
	Blocking          *config.StatusPolicy // Blocking ??????????????????
	PendingEnable     bool                 // Pending ????????????
	Pending           *config.StatusPolicy // Pending ??????????????????
	RetryingEnable    bool                 // Retrying ??????????????????
	Retrying          *config.StatusPolicy // Retrying ??????????????????
	UpgradeEnable     bool                 // ????????????
	UpgradeTopicLevel internal.TopicLevel  // ????????????????????????
	DegradeEnable     bool                 // ????????????
	DegradeTopicLevel internal.TopicLevel  // ????????????????????????
	//RerouteEnable     bool                  // PreReRoute ????????????, ??????false
	//Reroute           *config.ReroutePolicy // Handle???????????????????????????
}

// newLeveledConsumeDeciders create handlers based on different levels.
// the topics[0], xxxEnable, xxxStatusPolicy and (topics[0] + Upgrade/DegradeLevel) parameters is used in this construction.
func newLeveledConsumeDeciders(client *client, listener *consumeListener, options leveledConsumeDeciderOptions, deadHandler internalDecider) (*leveledConsumeDeciders, error) {
	deciders := &leveledConsumeDeciders{
		//multiStatusConsumeFacade: multiStatusConsumeFacade,
		//options:   options,
		//logger:      multiStatusConsumeFacade.logger,
	}
	if options.PendingEnable {
		suffix := message.StatusPending.TopicSuffix()
		hdOptions := statusDeciderOptions{status: message.StatusPending, msgGoto: message.GotoPending,
			topic: options.Topic + suffix, deaDecider: deadHandler, level: options.Level}
		decider, err := newStatusDecider(client, listener, options.Pending, hdOptions)
		if err != nil {
			return nil, err
		}
		deciders.pendingDecider = decider
	}
	if options.BlockingEnable {
		suffix := message.StatusBlocking.TopicSuffix()
		hdOptions := statusDeciderOptions{status: message.StatusBlocking, msgGoto: message.GotoBlocking,
			topic: options.Topic + suffix, deaDecider: deadHandler, level: options.Level}
		hd, err := newStatusDecider(client, listener, options.Blocking, hdOptions)
		if err != nil {
			return nil, err
		}
		deciders.blockingDecider = hd
	}
	if options.RetryingEnable {
		suffix := message.StatusRetrying.TopicSuffix()
		hdOptions := statusDeciderOptions{status: message.StatusRetrying, msgGoto: message.GotoRetrying,
			topic: options.Topic + suffix, deaDecider: deadHandler, level: options.Level}
		decider, err := newStatusDecider(client, listener, options.Retrying, hdOptions)
		if err != nil {
			return nil, err
		}
		deciders.retryingDecider = decider
	}
	if options.UpgradeEnable {
		gradeOpts := gradeOptions{topic: options.Topic, grade2Level: options.UpgradeTopicLevel,
			level: options.Level, msgGoto: message.GotoUpgrade}
		decider, err := newGradeDecider(client, listener, gradeOpts)
		if err != nil {
			return nil, err
		}
		deciders.upgradeDecider = decider
	}
	if options.DegradeEnable {
		gradeOpts := gradeOptions{topic: options.Topic, grade2Level: options.DegradeTopicLevel,
			level: options.Level, msgGoto: message.GotoDegrade}
		decider, err := newGradeDecider(client, listener, gradeOpts)
		if err != nil {
			return nil, err
		}
		deciders.degradeDecider = decider
	}
	return deciders, nil
}

func (hds leveledConsumeDeciders) Close() {
	if hds.blockingDecider != nil {
		hds.blockingDecider.close()
	}
	if hds.pendingDecider != nil {
		hds.pendingDecider.close()
	}
	if hds.retryingDecider != nil {
		hds.retryingDecider.close()
	}
	if hds.upgradeDecider != nil {
		hds.upgradeDecider.close()
	}
	if hds.degradeDecider != nil {
		hds.degradeDecider.close()
	}
}
