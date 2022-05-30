package soam

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
	UidTopicRouters    map[uint32]TopicLevel                    // Uid级别路由静态配置; 优先级低于Bucket级别路由;
	BucketTopicRouters map[string]TopicLevel                    // Bucket级别路由静态配置; 优先级高于Uid级别路由;
}

type ComsumerConfig struct {
	// extract from pulsar.ConsumerOptions
	//Topics                      []string
	Topic                       string
	SubscriptionName            string
	Type                        pulsar.SubscriptionType
	SubscriptionInitialPosition pulsar.SubscriptionInitialPosition
	NackBackoffPolicy           pulsar.NackBackoffPolicy

	Levels               []TopicLevel    // 默认L1, 且消费的Topic Level级别, len(Topics) == 1 or Topic存在的时候才生效
	LevelConsumeStrategy ConsumeStrategy // Topic级别消费策略

	ConsumeMaxTimes   int             // 最大消费次数
	ConsumeStrategy   ConsumeStrategy // 消费均衡策略
	Ready             *StatusPolicy   // Ready 主题检查策略
	BlockingEnable    bool            // Blocking 检查开关
	Blocking          *StatusPolicy   // Blocking 主题检查策略
	PendingEnable     bool            // Pending 检查开关
	Pending           *StatusPolicy   // Pending 主题检查策略
	RetryingEnable    bool            // Retrying 重试检查开关
	Retrying          *StatusPolicy   // Retrying 主题检查策略
	ReRouteEnable     bool            // PreReRoute 检查开关, 默认false
	reRouter          *ReRoutePolicy  // Handle失败时的动态重路由
	UpgradeEnable     bool            // 主动升级
	UpgradeTopicLevel TopicLevel      // 主动升级队列级别
	DegradeEnable     bool            // 主动降级
	DegradeTopicLevel TopicLevel      // 主动降级队列级别
	DeadEnable        bool            // 死信队列开关, 默认false; 如果所有校验器都没能校验通过, 应用代码需要自行Ack或者Nack
	DiscardEnable     bool            // 丢弃消息开关, 默认false

	// Dead Topic开关, 默认false; DeadToTopic string 默认 ${TOPIC}_DLQ, 固定后缀，不允许定制;
	// 需 pulsar.ConsumerOptions .DLQ 一起使用, 功能等同 pulsar.ConsumerOptions. RetryEnable;
	// 使用 RetryEnable, 行为跟pulsar client一致; 使用 DeadEnable 时，忽略DLQ中RetryLetterTopic,应该使用 Retrying *StatusPolicy

	// @Deprecated 不推荐使用
	RetryEnable bool
	DLQ         *DLQPolicy
}
