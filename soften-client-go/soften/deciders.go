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

type generalConsumeHandlers struct {
	rerouteHandler internalDecider // 重路由处理器: Reroute
	deadHandler    internalDecider // 状态处理器
	doneHandler    internalDecider // 状态处理器
	discardHandler internalDecider // 状态处理器
}

type generalConsumeHandlerOptions struct {
	Topic         string                // Business Topic
	DiscardEnable bool                  // Blocking 检查开关
	DeadEnable    bool                  // Pending 检查开关
	RerouteEnable bool                  // Retrying 重试检查开关
	Reroute       *config.ReroutePolicy // Reroute Policy
}

func newGeneralConsumeHandlers(client *client, conf generalConsumeHandlerOptions) (*generalConsumeHandlers, error) {
	handlers := &generalConsumeHandlers{}
	doneHandler, err := newFinalStatusHandler(client.logger, message.StatusDone)
	if err != nil {
		return nil, err
	}
	handlers.doneHandler = doneHandler
	if conf.DiscardEnable {
		hd, err := newFinalStatusHandler(client.logger, message.StatusDiscard)
		if err != nil {
			return nil, err
		}
		handlers.discardHandler = hd
	}
	if conf.DeadEnable {
		suffix, err := message.TopicSuffixOf(message.StatusDead)
		if err != nil {
			return nil, err
		}
		deadOptions := deadDecideOptions{topic: conf.Topic + suffix}
		hd, err := newDeadHandler(client, deadOptions)
		if err != nil {
			return nil, err
		}
		handlers.deadHandler = hd
	}
	if conf.RerouteEnable {
		hd, err := newRerouteHandler(client, conf.Reroute)
		if err != nil {
			return nil, err
		}
		handlers.rerouteHandler = hd
	}
	return handlers, nil
}

func (hds generalConsumeHandlers) Close() {
	if hds.rerouteHandler != nil {
		hds.rerouteHandler.close()
	}
	if hds.deadHandler != nil {
		hds.deadHandler.close()
	}
	if hds.doneHandler != nil {
		hds.doneHandler.close()
	}
	if hds.discardHandler != nil {
		hds.discardHandler.close()
	}
}

// ------ leveled consume handlers ------

type leveledConsumeHandlers struct {
	blockingHandler internalDecider // 状态处理器
	pendingHandler  internalDecider // 状态处理器
	retryingHandler internalDecider // 状态处理器
	upgradeHandler  internalDecider // 状态处理器: 升级为NewReady
	degradeHandler  internalDecider // 状态处理器: 升级为NewReady
}

type leveledConsumeHandlerOptions struct {
	Topic             string               // Business Topic
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

// newLeveledConsumeHandlers create handlers based on different levels.
// the topics[0], xxxEnable, xxxStatusPolicy and (topics[0] + Upgrade/DegradeLevel) parameters is used in this construction.
func newLeveledConsumeHandlers(client *client, options leveledConsumeHandlerOptions, deadHandler internalDecider) (*leveledConsumeHandlers, error) {
	handlers := &leveledConsumeHandlers{
		//multiStatusConsumeFacade: multiStatusConsumeFacade,
		//options:   options,
		//logger:      multiStatusConsumeFacade.logger,
	}
	if options.PendingEnable {
		suffix, err := message.TopicSuffixOf(message.StatusPending)
		if err != nil {
			return nil, err
		}
		hdOptions := statusHandleOptions{status: message.StatusPending,
			topic: options.Topic + suffix, deadHandler: deadHandler}
		hd, err := newStatusHandler(client, options.Pending, hdOptions)
		if err != nil {
			return nil, err
		}
		handlers.pendingHandler = hd
	}
	if options.BlockingEnable {
		suffix, err := message.TopicSuffixOf(message.StatusBlocking)
		if err != nil {
			return nil, err
		}
		hdOptions := statusHandleOptions{status: message.StatusBlocking,
			topic: options.Topic + suffix, deadHandler: deadHandler}
		hd, err := newStatusHandler(client, options.Blocking, hdOptions)
		if err != nil {
			return nil, err
		}
		handlers.pendingHandler = hd
	}
	if options.RetryingEnable {
		suffix, err := message.TopicSuffixOf(message.StatusBlocking)
		if err != nil {
			return nil, err
		}
		hdOptions := statusHandleOptions{status: message.StatusRetrying,
			topic: options.Topic + suffix, deadHandler: deadHandler}
		hd, err := newStatusHandler(client, options.Retrying, hdOptions)
		if err != nil {
			return nil, err
		}
		handlers.pendingHandler = hd
	}
	if options.UpgradeEnable {
		hd, err := newGradeHandler(client, options.Topic, options.UpgradeTopicLevel)
		if err != nil {
			return nil, err
		}
		handlers.upgradeHandler = hd
	}
	if options.DegradeEnable {
		hd, err := newGradeHandler(client, options.Topic, options.DegradeTopicLevel)
		if err != nil {
			return nil, err
		}
		handlers.degradeHandler = hd
	}
	return handlers, nil
}

func (hds leveledConsumeHandlers) Close() {
	if hds.blockingHandler != nil {
		hds.blockingHandler.close()
	}
	if hds.pendingHandler != nil {
		hds.pendingHandler.close()
	}
	if hds.retryingHandler != nil {
		hds.retryingHandler.close()
	}
	if hds.upgradeHandler != nil {
		hds.upgradeHandler.close()
	}
	if hds.degradeHandler != nil {
		hds.degradeHandler.close()
	}
}
