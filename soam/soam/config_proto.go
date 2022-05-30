package soam

import (
	"github.com/apache/pulsar-client-go/pulsar"
)

type TopicLevel string

const (
	L3  = TopicLevel("L3")
	L2  = TopicLevel("L2")
	L1  = TopicLevel("L1")
	B1  = TopicLevel("B1")
	B2  = TopicLevel("B2")
	DLQ = TopicLevel("DLQ")
)

var topicLevelOrders = map[TopicLevel]int{
	L3:  3,
	L2:  2,
	L1:  1,
	B1:  -1,
	B2:  -2,
	DLQ: -1025,
}

var DelayUnitMap = map[string]int{
	"s": 1, "S": 1,
	"m": 60, "M": 60,
	"h": 60 * 60, "H": 60 * 60,
}

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
	BackoffPolicy     StatusBackoffPolicy // 补偿策略, 优先级高于 BackoffDelayss
	ReentrantDelay    uint                // 重入 补偿延迟, 单状态固定时间
	ReentrantMaxTimes int                 // 重入 补偿延迟最大次数
	CheckerMandatory  bool                // Checker 存在性检查标识
	checkToTopic      string              // default ${TOPIC}_RETRYING, 固定后缀，不允许定制
	//NackMaxDelay      int      // Nack 补偿策略最大延迟粒度
	//NackMaxTimes      int      // Nack 补偿延迟最大次数
	//BackoffPolicy     StatusBackoffPolicy // 补偿策略, 优先级高于 BackoffDelays
}

//type ReRouteMode string

type ReRoutePolicy struct {
	//ReRouteMode      ReRouteMode       // 重路由模式: local; config
	UidPreRouters    map[uint32]string // Uid级别打散路由静态配置; 优先级低于Bucket级别路由;
	UidParseFunc     func(message pulsar.Message) uint
	BucketPreRouters map[string]string // Bucket级别打散路由静态配置; 优先级高于Uid级别路由;
	BucketParseFunc  func(message pulsar.Message) string
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
