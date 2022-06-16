package config

import (
	"time"

	"github.com/apache/pulsar-client-go/pulsar/log"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/shenqianjin/soften-client-go/soften/internal"
)

// ------ client configuration ------

type ClientConfig struct {
	URL                     string
	ConnectionTimeout       uint // Optional: default 5s
	OperationTimeout        uint // Optional: default 30s
	MaxConnectionsPerBroker uint // Optional: default 1

	Logger log.Logger `json: -`
}

// ------ producer configuration ------

type ProducerConfig struct {
	// extract out from pulsar.ProducerOptions
	Topic       string        // Required:
	SendTimeout time.Duration // Optional: s

	// custom
	TopicRouter        func(*pulsar.ProducerMessage) internal.TopicLevel // 自定义Topic路由器
	UidTopicRouters    map[uint32]internal.TopicLevel                    // Uid级别路由静态配置; 优先级低于Bucket级别路由;
	BucketTopicRouters map[string]internal.TopicLevel                    // Bucket级别路由静态配置; 优先级高于Uid级别路由;
}

// ------ consumer configuration (multi-level) ------

type MultiLevelConsumerConfig struct {
	*ConsumerConfig
	Levels               []internal.TopicLevel                // Required: 默认L1, 且消费的Topic Level级别, len(Topics) == 1 or Topic存在的时候才生效
	LevelBalanceStrategy internal.BalanceStrategy             // Optional: Topic级别消费策略
	LevelPolicies        map[internal.TopicLevel]*LevelPolicy // Optional: 级别消费策略
}

// ------ consumer configuration (multi-status) ------

type ConsumerConfig struct {
	// extract from pulsar.ConsumerOptions
	Concurrency                 uint                               // Optional: 并发控制
	Topics                      []string                           // Alternative with Topic: 如果有值, Topic 配置将被忽略
	Topic                       string                             // Alternative with Topics: Topics缺失的情况下，该值生效
	SubscriptionName            string                             //
	Level                       internal.TopicLevel                // Optional:
	Type                        pulsar.SubscriptionType            //
	SubscriptionInitialPosition pulsar.SubscriptionInitialPosition //
	NackBackoffPolicy           pulsar.NackBackoffPolicy           // Optional: Unrecommended, compatible with origin pulsar client
	NackRedeliveryDelay         time.Duration                      // Optional: Unrecommended, compatible with origin pulsar client
	RetryEnable                 bool                               // Optional: Unrecommended, compatible with origin pulsar client
	DLQ                         *DLQPolicy                         // Optional: Unrecommended, compatible with origin pulsar client
	ConsumeMaxTimes             int                                // Optional: 最大消费次数
	BalanceStrategy             internal.BalanceStrategy           // Optional: 消费均衡策略
	Ready                       *StatusPolicy                      // Optional: Ready 主题检查策略
	BlockingEnable              bool                               // Optional: Blocking 检查开关
	Blocking                    *StatusPolicy                      // Optional: Blocking 主题检查策略
	PendingEnable               bool                               // Optional: Pending 检查开关
	Pending                     *StatusPolicy                      // Optional: Pending 主题检查策略
	RetryingEnable              bool                               // Optional: Retrying 重试检查开关
	Retrying                    *StatusPolicy                      // Optional: Retrying 主题检查策略
	RerouteEnable               bool                               // Optional: PreReRoute 检查开关, 默认false
	Reroute                     *ReroutePolicy                     // Optional: Handle失败时的动态重路由
	UpgradeEnable               bool                               // Optional: 主动升级
	UpgradeTopicLevel           internal.TopicLevel                // Optional: 主动升级队列级别
	DegradeEnable               bool                               // Optional: 主动降级
	DegradeTopicLevel           internal.TopicLevel                // Optional: 主动降级队列级别
	DeadEnable                  bool                               // Optional: 死信队列开关, 默认false; 如果所有校验器都没能校验通过, 应用代码需要自行Ack或者Nack
	Dead                        *StatusPolicy                      // Optional: Dead 主题检查策略
	DiscardEnable               bool                               // Optional: 丢弃消息开关, 默认false

}

// ------ helper structs ------

// StatusPolicy 定义单状态的消费重入策略。
// 消费权重: 按整形值记录。
// 补偿策略:
// (1) 补偿延迟小于等于 NackBackoffMaxDelay(默认1min)时, 优先选择 Nack方式进行补偿;
// (2) 借助于 Reentrant 进行补偿, 每次重入代表额外增加一个 ReentrantDelay 延迟;
// (3) 如果 补偿延迟 - ReentrantDelay 仍然大于 NackBackoffMaxDelay, 那么会发生多次重入。
type StatusPolicy struct {
	ConsumeWeight     uint                // 消费权重
	ConsumeMaxTimes   int                 // 最大消费次数
	BackoffDelays     []string            // 补偿延迟 e.g: [5s, 2m, 1h], 如果值大于 ReentrantDelay 时，自动取整为 ReentrantDelay 的整数倍 (默认向下取整)
	BackoffPolicy     StatusBackoffPolicy // 补偿策略, 优先级高于 BackoffDelays
	ReentrantDelay    uint                // 重入 补偿延迟, 单状态固定时间
	ReentrantMaxTimes int                 // 重入 补偿延迟最大次数
	CheckerMandatory  bool                // Checker 存在性检查标识
	//NackMaxDelay      int      // Nack 补偿策略最大延迟粒度
	//NackMaxTimes      int      // Nack 补偿延迟最大次数
	//BackoffPolicy     StatusBackoffPolicy // 补偿策略, 优先级高于 BackoffDelays
}

type LevelPolicy struct {
	ConsumeWeight uint                // 消费权重
	UpgradeLevel  internal.TopicLevel // 升级级别
	DegradeLevel  internal.TopicLevel // 降级级别
}

type ReroutePolicy struct {
	//ReRouteMode      ReRouteMode       // 重路由模式: local; config
	UidPreRouters    map[uint32]string // Uid级别打散路由静态配置; 优先级低于Bucket级别路由;
	UidParseFunc     func(message pulsar.Message) uint
	BucketPreRouters map[string]string // Bucket级别打散路由静态配置; 优先级高于Uid级别路由;
	BucketParseFunc  func(message pulsar.Message) string
}

// DLQPolicy represents the configuration for the Dead Letter Queue multiStatusConsumeFacade policy.
type DLQPolicy struct {
	// MaxDeliveries specifies the maximum number of times that a message will be delivered before being
	// sent to the dead letter queue.
	MaxDeliveries uint32

	// DeadLetterTopic specifies the name of the topic where the failing messages will be sent.
	DeadLetterTopic string

	// RetryLetterTopic specifies the name of the topic where the retry messages will be sent.
	RetryLetterTopic string
}
