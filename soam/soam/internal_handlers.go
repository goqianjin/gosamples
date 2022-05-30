package soam

type consumeHandlers struct {
	blockingHandler internalHandler        // 状态处理器
	pendingHandler  internalHandler        // 状态处理器
	retryingHandler internalHandler        // 状态处理器
	deadHandler     internalHandler        // 状态处理器
	doneHandler     internalHandler        // 状态处理器
	discardHandler  internalHandler        // 状态处理器
	upgradeReRouter internalHandler        // 状态处理器: 升级为NewReady
	degradeReRouter internalHandler        // 状态处理器: 升级为NewReady
	rerouteHandler  internalRerouteHandler // 重路由处理器: Reroute
}

func newConsumeHandlers(consumer *consumer, config ComsumerConfig) (*consumeHandlers, error) {
	checkers := &consumeHandlers{
		//consumer: consumer,
		//config:   config,
		//log:      consumer.log,
	}
	if config.ReRouteEnable {
		hd, err := newRerouteHandler(consumer.log, consumer)
		if err != nil {
			return nil, err
		}
		checkers.rerouteHandler = hd
	}
	if config.PendingEnable {
		hd, err := newStatusHandler(consumer.log, consumer, MessageStatusPending, config.Pending)
		if err != nil {
			return nil, err
		}
		checkers.pendingHandler = hd
	}
	if config.BlockingEnable {
		hd, err := newStatusHandler(consumer.log, consumer, MessageStatusBlocking, config.Blocking)
		if err != nil {
			return nil, err
		}
		checkers.pendingHandler = hd
	}
	if config.RetryEnable {
		hd, err := newStatusHandler(consumer.log, consumer, MessageStatusRetrying, config.Retrying)
		if err != nil {
			return nil, err
		}
		checkers.pendingHandler = hd
	}
	if config.DeadEnable {
		hd, err := newStatusHandler(consumer.log, consumer, MessageStatusDead, nil)
		if err != nil {
			return nil, err
		}
		checkers.deadHandler = hd
	}
	if config.UpgradeEnable {
		hd, err := newGradeHandler(consumer.log, consumer, config.UpgradeTopicLevel)
		if err != nil {
			return nil, err
		}
		checkers.upgradeReRouter = hd
	}
	if config.DegradeEnable {
		hd, err := newGradeHandler(consumer.log, consumer, config.UpgradeTopicLevel)
		if err != nil {
			return nil, err
		}
		checkers.upgradeReRouter = hd
	}
	return checkers, nil
}
