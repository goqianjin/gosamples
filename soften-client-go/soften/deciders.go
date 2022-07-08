package soften

import (
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/shenqianjin/soften-client-go/soften/checker"
	"github.com/shenqianjin/soften-client-go/soften/config"
	"github.com/shenqianjin/soften-client-go/soften/internal"
	"github.com/shenqianjin/soften-client-go/soften/message"
)

type internalDecider interface {
	Decide(msg pulsar.ConsumerMessage, checkStatus checker.CheckStatus) (success bool)
	close()
}

// ------ general consume handlers ------

type generalConsumeDeciders struct {
	rerouteDecider internalDecider // 重路由处理器: Reroute
	deadDecider    internalDecider // 状态处理器
	doneDecider    internalDecider // 状态处理器
	discardDecider internalDecider // 状态处理器
}

type generalConsumeDeciderOptions struct {
	Topic         string                // Business Topic
	DiscardEnable bool                  // Blocking 检查开关
	DeadEnable    bool                  // Pending 检查开关
	RerouteEnable bool                  // Retrying 重试检查开关
	Reroute       *config.ReroutePolicy // Reroute Policy
}

func newGeneralConsumeDeciders(client *client, listener *consumeListener, conf generalConsumeDeciderOptions) (*generalConsumeDeciders, error) {
	handlers := &generalConsumeDeciders{}
	doneHandler, err := newFinalStatusHandler(client, listener, message.GotoDone)
	if err != nil {
		return nil, err
	}
	handlers.doneDecider = doneHandler
	if conf.DiscardEnable {
		hd, err := newFinalStatusHandler(client, listener, message.GotoDiscard)
		if err != nil {
			return nil, err
		}
		handlers.discardDecider = hd
	}
	if conf.DeadEnable {
		suffix, err := message.TopicSuffixOf(message.StatusDead)
		if err != nil {
			return nil, err
		}
		deadOptions := deadDecideOptions{topic: conf.Topic + suffix}
		hd, err := newDeadHandler(client, listener, deadOptions)
		if err != nil {
			return nil, err
		}
		handlers.deadDecider = hd
	}
	if conf.RerouteEnable {
		hd, err := newRerouteHandler(client, listener, conf.Reroute)
		if err != nil {
			return nil, err
		}
		handlers.rerouteDecider = hd
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
	blockingDecider internalDecider // 状态处理器
	pendingDecider  internalDecider // 状态处理器
	retryingDecider internalDecider // 状态处理器
	upgradeDecider  internalDecider // 状态处理器: 升级为NewReady
	degradeDecider  internalDecider // 状态处理器: 升级为NewReady
}

type leveledConsumeDeciderOptions struct {
	Topic             string               // Business Topic
	Level             internal.TopicLevel  // level
	BlockingEnable    bool                 // Blocking 检查开关
	Blocking          *config.StatusPolicy // Blocking 主题检查策略
	PendingEnable     bool                 // Pending 检查开关
	Pending           *config.StatusPolicy // Pending 主题检查策略
	RetryingEnable    bool                 // Retrying 重试检查开关
	Retrying          *config.StatusPolicy // Retrying 主题检查策略
	UpgradeEnable     bool                 // 主动升级
	UpgradeTopicLevel internal.TopicLevel  // 主动升级队列级别
	DegradeEnable     bool                 // 主动降级
	DegradeTopicLevel internal.TopicLevel  // 主动升级队列级别
	//RerouteEnable     bool                  // PreReRoute 检查开关, 默认false
	//Reroute           *config.ReroutePolicy // Handle失败时的动态重路由
}

// newLeveledConsumeDeciders create handlers based on different levels.
// the topics[0], xxxEnable, xxxStatusPolicy and (topics[0] + Upgrade/DegradeLevel) parameters is used in this construction.
func newLeveledConsumeDeciders(client *client, listener *consumeListener, options leveledConsumeDeciderOptions, deadHandler internalDecider) (*leveledConsumeDeciders, error) {
	handlers := &leveledConsumeDeciders{
		//multiStatusConsumeFacade: multiStatusConsumeFacade,
		//options:   options,
		//logger:      multiStatusConsumeFacade.logger,
	}
	if options.PendingEnable {
		suffix, err := message.TopicSuffixOf(message.StatusPending)
		if err != nil {
			return nil, err
		}
		hdOptions := statusHandleOptions{status: message.StatusPending, msgGoto: message.GotoPending,
			topic: options.Topic + suffix, deadHandler: deadHandler, level: options.Level}
		hd, err := newStatusHandler(client, listener, options.Pending, hdOptions)
		if err != nil {
			return nil, err
		}
		handlers.pendingDecider = hd
	}
	if options.BlockingEnable {
		suffix, err := message.TopicSuffixOf(message.StatusBlocking)
		if err != nil {
			return nil, err
		}
		hdOptions := statusHandleOptions{status: message.StatusBlocking, msgGoto: message.GotoBlocking,
			topic: options.Topic + suffix, deadHandler: deadHandler, level: options.Level}
		hd, err := newStatusHandler(client, listener, options.Blocking, hdOptions)
		if err != nil {
			return nil, err
		}
		handlers.pendingDecider = hd
	}
	if options.RetryingEnable {
		suffix, err := message.TopicSuffixOf(message.StatusRetrying)
		if err != nil {
			return nil, err
		}
		hdOptions := statusHandleOptions{status: message.StatusRetrying, msgGoto: message.GotoRetrying,
			topic: options.Topic + suffix, deadHandler: deadHandler, level: options.Level}
		hd, err := newStatusHandler(client, listener, options.Retrying, hdOptions)
		if err != nil {
			return nil, err
		}
		handlers.pendingDecider = hd
	}
	if options.UpgradeEnable {
		gradeOpts := gradeOptions{topic: options.Topic, grade2Level: options.UpgradeTopicLevel,
			level: options.Level, msgGoto: message.GotoUpgrade}
		hd, err := newGradeHandler(client, listener, gradeOpts)
		if err != nil {
			return nil, err
		}
		handlers.upgradeDecider = hd
	}
	if options.DegradeEnable {
		gradeOpts := gradeOptions{topic: options.Topic, grade2Level: options.DegradeTopicLevel,
			level: options.Level, msgGoto: message.GotoDegrade}
		hd, err := newGradeHandler(client, listener, gradeOpts)
		if err != nil {
			return nil, err
		}
		handlers.degradeDecider = hd
	}
	return handlers, nil
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
