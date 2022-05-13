package lib_sbmc

import "time"

type ClientConfig struct {
	URL               string
	ConnectionTimeout time.Duration
}

type ProducerConfig struct {
	UidRouter    map[string]string // Uid级别路由静态配置; 优先级低于Itbl级别路由; 如果consul中存在相同的key,静态路由会被覆盖
	ItblRouter   map[string]string // Itbl级别路由静态配置; 优先级高于Uid级别路由; 如果consul中存在相同的key,静态路由会被覆盖
	ConsumeDelay int               // 消费延迟时间, 单位s;
}

type ComsumerConfig struct {
	SubscriptionName string
	Topic            string

	PendingEnable              bool     // Pending Topic开关, 默认false
	PendingNackBackoffPolicy   []string // default [1s 2s 3s 5s]
	PendingNackBackoffMaxTimes int      // default 4, 累积10s间隔
	PendingBackoffDelay        int      // default 50s
	PendingBackoffMaxTimes     int      // default 30, 累积延迟 30 * (10s+50s) = 30min
	//PendingTo                  string // default ${TOPIC}_PENDING, 固定后缀，不允许定制

	BlockingEnable              bool     // Blocking Topic开关, 默认false
	BlockingNackBackoffPolicy   []string // default []
	BlockingNackBackoffMaxTimes int      // default 0
	BlockingBackoffDelay        int      // default 1h
	BlockingBackoffMaxTimes     int      // default 24, 累积延迟 24 * (0+1h) = 24h
	//BlockingTo                  string // default ${TOPIC}_BLOCKING, 固定后缀，不允许定制

	RetryingEnable              bool     // Retrying Topic开关, 默认false
	RetryingNackBackoffPolicy   []string // default [1s 2s 3s 5s]
	RetryingNackBackoffMaxTimes int      // default 4, 累积10s间隔
	RetryingBackoffDelay        int      // default 50s
	RetryingBackoffMaxTimes     int      // default 30, 累积延迟 30 * (10s+50s) = 30min
	//RetryingTo                  string // default ${TOPIC}_RETRYING, 固定后缀，不允许定制

	DeadEnable bool // Dead Topic开关, 默认false
	//DeadLetterTopic string // default ${TOPIC}_DLQ, 固定后缀，不允许定制

	ReRouteEnable bool
	UidRerouter   map[string]string // Uid级别打散路由静态配置; 优先级低于Itbl级别路由; 如果consul中存在相同的key,静态路由会被覆盖
	ItblRerouter  map[string]string // Itbl级别打散路由静态配置; 优先级高于Uid级别路由; 如果consul中存在相同的key,静态路由会被覆盖

	Topics []string

	TopicsPattern string

	AutoDiscoveryPeriod time.Duration

	Properties map[string]string

	SubscriptionProperties map[string]string

	Type SubscriptionType

	SubscriptionInitialPosition

	DLQ *DLQPolicy

	KeySharedPolicy *KeySharedPolicy

	RetryEnable bool

	MessageChannel chan ConsumerMessage

	ReceiverQueueSize int

	NackRedeliveryDelay time.Duration

	Name string

	ReadCompacted bool

	ReplicateSubscriptionState bool

	Interceptors ConsumerInterceptors

	Schema Schema

	MaxReconnectToBroker *uint

	Decryption *MessageDecryptionInfo

	EnableDefaultNackBackoffPolicy bool

	NackBackoffPolicy NackBackoffPolicy
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
