package soften

import (
	"fmt"

	"github.com/shenqianjin/soften-client-go/soften/config"
	"github.com/shenqianjin/soften-client-go/soften/internal"

	"github.com/apache/pulsar-client-go/pulsar/log"
)

type multiLevelConsumer struct {
	logger         log.Logger
	messageCh      chan ConsumerMessage                        // channel used to deliver message to application
	levelStrategy  internal.BalanceStrategy                    // 消费策略
	levelPolicies  map[internal.TopicLevel]*config.LevelPolicy // 级别消费策略
	levelConsumers map[internal.TopicLevel]*multiStatusConsumer
}

func newMultiLevelConsumer(consumerLogger log.Logger, client *client, conf *config.ConsumerConfig, messageCh chan ConsumerMessage, levelHandlers map[internal.TopicLevel]*leveledConsumeHandlers) (*multiLevelConsumer, error) {
	consumer := &multiLevelConsumer{
		logger:        consumerLogger,
		levelStrategy: conf.LevelBalanceStrategy,
		messageCh:     messageCh,
	}
	consumer.levelConsumers = make(map[internal.TopicLevel]*multiStatusConsumer, len(conf.Levels))
	for _, level := range conf.Levels {
		levelConsumer, err := newMultiStatusConsumer(consumerLogger, client, level, conf, make(chan ConsumerMessage, 10), levelHandlers[level])
		if err != nil {
			return nil, fmt.Errorf("failed to new multi-status comsumer -> %v", err)
		}
		consumer.levelConsumers[level] = levelConsumer
	}
	// start to listen message from all status multiStatusConsumer
	go consumer.retrieveStatusMessages()
	return consumer, nil
}

func (c *multiLevelConsumer) retrieveStatusMessages() {
	chs := make([]<-chan ConsumerMessage, len(c.levelConsumers))
	weights := make([]uint, len(c.levelConsumers))
	for level, consumer := range c.levelConsumers {
		chs = append(chs, consumer.Chan())
		weights = append(weights, c.levelPolicies[level].ConsumeWeight)
	}
	balanceStrategy, err := config.BuildStrategy(c.levelStrategy, weights)
	if err != nil {
		panic(fmt.Errorf("failed to start retrieve: %v", err))
	}
	for {
		msg, ok := messageChSelector.receiveOneByWeight(chs, balanceStrategy, &[]int{})
		if !ok {
			c.logger.Warnf("status chan closed")
			break
		}
		// 获取到消息
		if msg.Message != nil && msg.Consumer != nil {
			fmt.Printf("received message  msgId: %v -- content: '%s'\n", msg.ID(), string(msg.Payload()))
			c.messageCh <- msg
		} else {
			panic(fmt.Sprintf("consumed an invalid message: %v", msg))
		}
	}
}
