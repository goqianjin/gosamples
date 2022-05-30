package some

import (
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

type ClientConfig struct {
	URL               string
	ConnectionTimeout time.Duration
}

type ProducerConfig struct {
	// extract out from pulsar.ProducerOptions
	Topic       string
	SendTimeout time.Duration

	// custom
	TopicRouter        func(*pulsar.ProducerMessage) TopicLevel // 自定义Topic路由器
	UidTopicRouters    map[string]TopicLevel                    // Uid级别路由静态配置; 优先级低于Bucket级别路由;
	BucketTopicRouters map[string]TopicLevel                    // Bucket级别路由静态配置; 优先级高于Uid级别路由;
}

type TopicLevel string

const (
	L3  = TopicLevel("L3")
	L2  = TopicLevel("L2")
	L1  = TopicLevel("L1")
	B1  = TopicLevel("B1")
	B2  = TopicLevel("B2")
	DLQ = TopicLevel("DLQ")
)

var DelayUnitMap = map[string]time.Duration{
	"s": time.Second, "S": time.Second,
	"m": time.Minute, "M": time.Minute,
	"h": time.Hour, "H": time.Hour,
}

type ConsumerOptions struct {
	ComsumerConfig
}

type ComsumerConfig struct {
	// extract from pulsar.ConsumerOptions
	//Topics           []string
	Topic                       string
	SubscriptionName            string
	Type                        pulsar.SubscriptionType
	SubscriptionInitialPosition pulsar.SubscriptionInitialPosition
	NackBackoffPolicy           pulsar.NackBackoffPolicy

	// Pending Pending主题检查策略
	// Enable开关, 默认false; ConsumeWeight 消费权重默认为2; CheckToTopic default ${TOPIC}_PENDING, 固定后缀，不允许定制;
	// NackBackoffPolicy default [1s 2s 3s 4s]; NackBackoffMaxTimes default 4, 累积10s间隔;
	// BackoffDelay default 50s; BackoffMaxTimes default 30, 累积延迟 30 * (10s+50s) = 30min
	Pending *CheckPolicy

	// Blocking Blocking主题检查策略
	// Enable开关, 默认false; ConsumeWeight 消费权重默认为2; CheckToTopic default ${TOPIC}_BLOCKING, 固定后缀，不允许定制;
	// NackBackoffPolicy default []; NackBackoffMaxTimes default 0, 累积0s间隔;
	// BackoffDelay default 1h; BackoffMaxTimes default 6, 累积延迟 6 * (0+1h) = 6h
	Blocking *CheckPolicy

	// Blocking Retrying主题检查策略
	// Enable开关, 默认false; ConsumeWeight 消费权重默认为2; CheckToTopic default ${TOPIC}_BLOCKING, 固定后缀，不允许定制;
	// NackBackoffPolicy default [1s 2s 3s 4s]; NackBackoffMaxTimes default 4, 累积10s间隔;
	// BackoffDelay default 50s; BackoffMaxTimes default 30, 累积延迟 30 * (10s+50s) = 30min
	Retrying *CheckPolicy

	// Dead Topic开关, 默认false; DeadToTopic string 默认 ${TOPIC}_DLQ, 固定后缀，不允许定制;
	// 需 pulsar.ConsumerOptions .DLQ 一起使用, 功能等同 pulsar.ConsumerOptions. RetryEnable;
	// 使用 RetryEnable, 行为跟pulsar client一致; 使用 DeadEnable 时，忽略DLQ中RetryLetterTopic,应该使用 Retrying *CheckPolicy
	DeadEnable bool
	// @Deprecated 不推荐使用
	RetryEnable bool
	DLQ         *DLQPolicy

	// PreReRouter 处理消息的前置重路由
	PreReRouter *RoutePolicy
	// ReRouter Handle失败时的动态重路由
	ReRouter *RoutePolicy
}

type CheckPolicy struct {
	Enable              bool     // Retrying Topic开关, 默认false
	ConsumeWeight       uint     // 消费权重
	NackBackoffPolicy   []string // default [1s 2s 3s 4s]
	NackBackoffMaxTimes uint32   // default 4, 累积10s间隔
	BackoffDelay        int      // default 50s
	BackoffMaxTimes     uint32   // default 30, 累积延迟 30 * (10s+50s) = 30min
	CheckerMandatory    bool     // 强制需要checker开关: 默认false; 若为true, 订阅时如果没有checker会报错
	//CheckToTopic        string   // default ${TOPIC}_RETRYING, 固定后缀，不允许定制
}

type RoutePolicy struct {
	Enable           bool
	UidPreRouters    map[string]string // Uid级别打散路由静态配置; 优先级低于Bucket级别路由;
	BucketPreRouters map[string]string // Itbl级别打散路由静态配置; 优先级高于Uid级别路由;

}

// DLQPolicy represents the configuration for the Dead Letter Queue consumer policy.
type DLQPolicy struct {
	// MaxDeliveries specifies the maximum number of times that a message will be delivered before being
	// sent to the dead letter queue.
	MaxDeliveries uint32

	// DeadLetterTopic specifies the name of the topic where the failing messages will be sent.
	DeadLetterTopic string

	// RetryLetterTopic specifies the name of the topic where the retry messages will be sent.
	RetryLetterTopic string
}
