package soam

import (
	"errors"
	"fmt"
	"log"

	"github.com/apache/pulsar-client-go/pulsar"
)

type client struct {
	pulsar.Client
	config ClientConfig
}

func NewClient(config ClientConfig) (*client, error) {
	clientOption := pulsar.ClientOptions{
		URL:               config.URL,
		ConnectionTimeout: config.ConnectionTimeout,
	}
	pulsarClient, err := pulsar.NewClient(clientOption)
	if err != nil {
		return nil, err
	}
	cli := &client{Client: pulsarClient, config: config}
	return cli, nil
}

func (c *client) subscribeByStatus(config ComsumerConfig, status messageStatus) (pulsar.Consumer, error) {
	suffix := ""
	if status == MessageStatusReady {
		suffix = ""
	} else if status == MessageStatusDead {
		suffix = "_DLQ"
	} else if status == MessageStatusPending || status == MessageStatusBlocking || status == MessageStatusRetrying {
		suffix = "_" + string(status)
	} else {
		return nil, errors.New(fmt.Sprintf("message status %s cannot be subsribed", status))
	}
	topic := config.Topic + suffix
	subscriptionName := config.SubscriptionName + suffix
	pulsarConsumer, err := c.Client.Subscribe(pulsar.ConsumerOptions{
		Topic:                       topic,
		SubscriptionName:            subscriptionName,
		Type:                        config.Type,
		SubscriptionInitialPosition: config.SubscriptionInitialPosition,
		NackBackoffPolicy:           config.NackBackoffPolicy,
		DLQ: &pulsar.DLQPolicy{
			MaxDeliveries:   config.DLQ.MaxDeliveries,
			DeadLetterTopic: config.Topic + "_DLQ",
		},
		MessageChannel: nil,
	})
	return pulsarConsumer, err
}

func (c *client) SubscribeInRegular(config ComsumerConfig, handler Handler, checkpoints ...*checkpoint) error {
	if handler == nil {
		panic("consumer handler cannot be empty")
	}
	handlerWithState := func(message pulsar.Message) handleResult {
		success, err := handler(message)
		if success {
			return HandledOk
		} else {
			return HandledFail.Err(err)
		}
	}
	return c.SubscribeInPremium(config, handlerWithState, checkpoints...)

}

func (c *client) SubscribeInPremium(config ComsumerConfig, handler HandlerInPremium, checkpoints ...*checkpoint) error {
	// validate and default config
	if handler == nil {
		panic("consumer handler cannot be empty")
	}
	if len(config.Levels) <= 0 {
		config.Levels = []TopicLevel{L1}
	}
	if config.ConsumeStrategy == "" {
		config.ConsumeStrategy = ConsumeStrategyRandRound
	}
	//
	optMap := make(map[CheckType]*checkpoint)
	for _, opt := range checkpoints {
		optMap[opt.checkType] = opt
	}
	// 校验配置policy
	if config.Ready == nil {
		config.Ready = DefaultCheckPolicyReady
	} else {
		if err := c.validateAndDefaultCheckPolicy(config.Ready, DefaultCheckPolicyReady); err != nil {
			return err
		}
	}
	if config.PendingEnable {
		if err := c.validateAndDefaultCheckPolicy(config.Pending, DefaultCheckPolicyPending); err != nil {
			return err
		}
		if config.Pending.CheckerMandatory && c.findCheckpointByType(checkpoints, CheckTypePrePending, CheckTypePostPending) == nil {
			return errors.New(fmt.Sprintf("[%s] checkOption is missing. please add one or disable the mandatory if necessary", MessageStatusPending))
		}
	}
	if config.BlockingEnable {
		if err := c.validateAndDefaultCheckPolicy(config.Blocking, DefaultCheckPolicyBlocking); err != nil {
			return err
		}
		if config.Pending.CheckerMandatory && c.findCheckpointByType(checkpoints, CheckTypePreBlocking, CheckTypePostBlocking) == nil {
			return errors.New(fmt.Sprintf("[%s] checkOption is missing. please add one or disable the mandatory if necessary", MessageStatusBlocking))
		}
	}
	if config.RetryingEnable {
		if err := c.validateAndDefaultCheckPolicy(config.Retrying, DefaultCheckPolicyRetrying); err != nil {
			return err
		}
		if config.Pending.CheckerMandatory && c.findCheckpointByType(checkpoints, CheckTypePreRetrying, CheckTypePostRetrying) == nil {
			return errors.New(fmt.Sprintf("[%s] checkOption is missing. please add one or disable the mandatory if necessary", MessageStatusRetrying))
		}
	}
	maxLevel := config.Levels[0]
	minLevel := config.Levels[0]
	/*for _, level := range config.Levels {
		if topicLevelOrders[level] < topicLevelOrders[minLevel] {
			minLevel = level
		}
		if topicLevelOrders[level] > topicLevelOrders[maxLevel] {
			maxLevel = level
		}
	}*/
	// 通过订阅option指定升级的TopicLevel可能不存在，对于升级要求配置文件必须写明升级到哪里
	if config.UpgradeEnable {
		if config.UpgradeTopicLevel == "" {
			return errors.New("missing UpgradeTopicLevel configuration")
		}
		if topicLevelOrders[config.UpgradeTopicLevel] <= topicLevelOrders[maxLevel] {
			return errors.New("upgrade topic level show be higher than the maximum of consume levels")
		}
	}
	if config.DegradeEnable {
		if config.DegradeTopicLevel == "" {
			return errors.New("missing DegradeTopicLevel configuration")
		}
		if topicLevelOrders[config.UpgradeTopicLevel] >= topicLevelOrders[minLevel] {
			return errors.New("degrade topic level show be lower than the minimum of consume levels")
		}
	}
	// 校验checker: checker可以在对应配置enable=false的情况下存在
	checkpointMap := make(map[CheckType]*checkpoint)
	for _, checkOpt := range checkpoints {
		if checkOpt.checkType == "" {
			return errors.New("checkType can not be empty")
		} else if checkOpt.preStatusChecker == nil {
			return errors.New(fmt.Sprintf("preStatusChecker can not be nil for input checkOption: %s", checkOpt.checkType))
		}
		checkpointMap[checkOpt.checkType] = checkOpt
	}

	// create consumer
	consumer, err := newConsumer(c, config)
	if err != nil {
		return err
	}

	// listen messages
	go consumer.start()
	// consume
	for message := range consumer.messageCh {
		consumer.consume(handler, message)
	}
	// cleanup
	defer consumer.Close()
	// unsubscribe if any error happens
	if err := consumer.Unsubscribe(); err != nil {
		log.Fatal(err)
	}
	return err
}

func (c *client) validateAndDefaultCheckPolicy(configuredPolicy *StatusPolicy, defaultPolicy *StatusPolicy) error {
	if configuredPolicy == nil {
		configuredPolicy = defaultPolicy
		return nil
	}
	if configuredPolicy.ConsumeWeight == 0 {
		configuredPolicy.ConsumeWeight = defaultPolicy.ConsumeWeight
	}
	if configuredPolicy.ConsumeMaxTimes == 0 {
		configuredPolicy.ConsumeMaxTimes = defaultPolicy.ConsumeMaxTimes
	}
	if configuredPolicy.BackoffDelays == nil && configuredPolicy.BackoffPolicy == nil {
		configuredPolicy.BackoffDelays = defaultPolicy.BackoffDelays
		configuredPolicy.BackoffPolicy = defaultPolicy.BackoffPolicy
	}
	if configuredPolicy.ReentrantDelay == 0 {
		configuredPolicy.ReentrantDelay = defaultPolicy.ReentrantDelay
	}
	if configuredPolicy.ReentrantMaxTimes == 0 {
		configuredPolicy.ReentrantMaxTimes = defaultPolicy.ReentrantMaxTimes
	}
	// default policy
	if configuredPolicy.BackoffPolicy == nil && configuredPolicy.BackoffDelays != nil {
		if backoffPolicy, err := newAbbrStatusBackoffPolicy(configuredPolicy.BackoffDelays); err != nil {
			return err
		} else {
			configuredPolicy.BackoffDelays = nil // release unnecessary reference
			configuredPolicy.BackoffPolicy = backoffPolicy
		}
	}
	return nil
}

func (c *client) findCheckpointByType(checkpoints []*checkpoint, checkTypes ...CheckType) *checkpoint {
	for _, opt := range checkpoints {
		for _, checkType := range checkTypes {
			if opt.checkType == checkType {
				return opt
			}
		}
	}
	return nil
}
