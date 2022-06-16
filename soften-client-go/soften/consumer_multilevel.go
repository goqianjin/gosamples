package soften

import (
	"fmt"
	"time"

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

func newMultiLevelConsumer(consumerLogger log.Logger, client *client, conf config.MultiLevelConsumerConfig, messageCh chan ConsumerMessage, levelHandlers map[internal.TopicLevel]*leveledConsumeHandlers) (*multiLevelConsumer, error) {
	consumer := &multiLevelConsumer{
		logger:        consumerLogger,
		levelStrategy: conf.LevelBalanceStrategy,
		messageCh:     messageCh,
	}
	consumer.levelConsumers = make(map[internal.TopicLevel]*multiStatusConsumer, len(conf.Levels))
	for _, level := range conf.Levels {
		levelConsumer, err := newMultiStatusConsumer(consumerLogger, client, level, conf.ConsumerConfig, make(chan ConsumerMessage, 10), levelHandlers[level])
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
	multiConsumers := make([]*multiStatusConsumer, len(c.levelConsumers))
	weights := make([]uint, len(c.levelConsumers))
	for level, consumer := range c.levelConsumers {
		multiConsumers = append(multiConsumers, consumer)
		weights = append(weights, c.levelPolicies[level].ConsumeWeight)
	}
	consumeStrategy, err := config.BuildStrategy(c.levelStrategy, weights)
	if err != nil {
		panic(fmt.Errorf("failed to start retrieveStatusMessages: %v", err))
	}
	nilCount := 0 // TODO: use notify chan ?
	for {
		if nilCount == len(multiConsumers) {
			time.Sleep(50 * time.Millisecond)
			nilCount = 0
		}
		var message ConsumerMessage
		messageChan := multiConsumers[consumeStrategy.Next()].Chan()
		select {
		case msg, ok := <-messageChan:
			if !ok {
				c.logger.Info("multiLevelConsumeFacade chan is closed")
				break
			}
			message = msg
		default:
			nilCount++
			continue
		}
		// 获取到消息
		if message.Message != nil && message.Consumer != nil {
			fmt.Printf("received message  msgId: %v -- content: '%s'\n", message.ID(), string(message.Payload()))
			c.messageCh <- message
		} else {
			panic(fmt.Sprintf("consumed an invalid message: %v", message))
		}
	}
}
