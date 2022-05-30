package soam

const (
	ConsumeWeightMain     = uint(10) // 业务Main  队列: 50% 权重
	ConsumeWeightRetrying = uint(6)  // Retrying 队列: 30% 权重
	ConsumeWeightPending  = uint(3)  // Pending  队列: 15% 权重
	ConsumeWeightBlocking = uint(1)  // Blocking 队列:  5% 权重
)

var (
	DefaultConsumeMaxTimes     uint32 = 30                                              // 一个消息整个生命周期中的最大消费次数
	DefaultStatusBackoffDelays        = []string{"3s", "5s", "8s", "14s", "30s", "60s"} // 消息默认Nack重试间隔策略策略
	DefaultNackMaxDelay               = 300                                             // 最大Nack延迟，默认5分钟

	// DefaultCheckPolicyReady 默认pending状态的校验策略。CheckToTopic default ${TOPIC}, 固定后缀，不允许定制;
	// 默认开启; 默认权重 10; CheckerMandatory 默认false;
	// BackoffPolicy: [3s 5s 8s 14s 30s 60s]; NackBackoffMaxTimes: 6, 累积120s间隔;
	// ReentrantDelay: 0; ReentrantMaxTimes: 0, 累积延迟 0*(0s+120s) + 120s = 2min 内处理不成功到Retrying 或 Dead 或 Discard。
	DefaultCheckPolicyReady = &StatusPolicy{
		ConsumeWeight:     ConsumeWeightMain,
		ConsumeMaxTimes:   DefaultConsumeMaxTimes,
		BackoffDelays:     DefaultStatusBackoffDelays,
		BackoffPolicy:     nil,
		ReentrantDelay:    0, // 不需要
		ReentrantMaxTimes: 0, // 不需要
		CheckerMandatory:  false,
	}

	// DefaultCheckPolicyRetrying 默认pending状态的校验策略。CheckToTopic default ${TOPIC}_PENDING, 固定后缀，不允许定制;
	// 默认开启; 默认权重 10; CheckerMandatory 默认false;
	// BackoffPolicy: [3s 5s 8s 14s 30s 60s]; NackBackoffMaxTimes: 6, 累积120s间隔;
	// ReentrantDelay: 60s; ReentrantMaxTimes: 15, 累积延迟 15 * (60s+60s) = 30min 内处理不成功到DLQ 或 Discard。
	DefaultCheckPolicyRetrying = &StatusPolicy{
		ConsumeWeight:     ConsumeWeightRetrying,
		ConsumeMaxTimes:   DefaultConsumeMaxTimes,
		BackoffDelays:     DefaultStatusBackoffDelays,
		BackoffPolicy:     nil,
		ReentrantDelay:    120,
		ReentrantMaxTimes: 15,
		CheckerMandatory:  false,
	}

	// DefaultCheckPolicyPending 同 DefaultCheckPolicyRetrying。
	DefaultCheckPolicyPending = &StatusPolicy{
		ConsumeWeight:     ConsumeWeightPending,
		ConsumeMaxTimes:   DefaultConsumeMaxTimes,
		BackoffDelays:     DefaultStatusBackoffDelays,
		BackoffPolicy:     nil,
		ReentrantDelay:    120,
		ReentrantMaxTimes: 15,
		CheckerMandatory:  false,
	}

	// DefaultCheckPolicyBlocking 默认pending状态的校验策略。CheckToTopic default ${TOPIC}_PENDING, 固定后缀，不允许定制;
	// 默认开启; 默认权重 10; CheckerMandatory 默认false;
	// BackoffPolicy: []; NackBackoffMaxTimes: 0, 累积0s间隔;
	// ReentrantDelay: 1800s; ReentrantMaxTimes: 12, 累积延迟 12 * (1800s+0s) = 6h 内处理不成功到DLQ 或 Discard。
	DefaultCheckPolicyBlocking = &StatusPolicy{
		ConsumeWeight:     ConsumeWeightBlocking,
		ConsumeMaxTimes:   DefaultConsumeMaxTimes,
		BackoffDelays:     nil,
		BackoffPolicy:     nil,
		ReentrantDelay:    1800,
		ReentrantMaxTimes: 12,
		CheckerMandatory:  false,
	}
)
