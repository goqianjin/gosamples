package soften

import (
	"fmt"

	"github.com/shenqianjin/soften-client-go/soften/message"

	"github.com/shenqianjin/soften-client-go/soften/config"

	"github.com/shenqianjin/soften-client-go/soften/internal"

	"github.com/apache/pulsar-client-go/pulsar/log"
)

type multiStatusConsumer struct {
	logger           log.Logger
	messageCh        chan ConsumerMessage // channel used to deliver message to application
	level            internal.TopicLevel
	statusStrategy   internal.BalanceStrategy // 消费策略
	mainConsumer     *statusConsumer
	pendingConsumer  *statusConsumer
	blockingConsumer *statusConsumer
	retryingConsumer *statusConsumer
}

func newMultiStatusConsumer(consumerLogger log.Logger, client *client, level internal.TopicLevel, conf *config.ConsumerConfig,
	messageCh chan ConsumerMessage, handlers *leveledConsumeHandlers) (*multiStatusConsumer, error) {
	cs := &multiStatusConsumer{logger: consumerLogger.SubLogger(log.Fields{"level": level}),
		level: level, messageCh: messageCh, statusStrategy: conf.BalanceStrategy}
	// create status multiStatusConsumer
	if mainConsumer, err := client.subscribeByStatus(conf, message.StatusReady); err != nil {
		return nil, err
	} else {
		cs.mainConsumer = newStatusConsumer(mainConsumer, message.StatusReady, conf.Ready, nil)
	}
	if conf.PendingEnable {
		pendingConsumer, err := client.subscribeByStatus(conf, message.StatusPending)
		if err != nil {
			return nil, err
		}
		cs.pendingConsumer = newStatusConsumer(pendingConsumer, message.StatusPending, conf.Pending, handlers.pendingHandler)
	}
	if conf.BlockingEnable {
		bc, err := client.subscribeByStatus(conf, message.StatusBlocking)
		if err != nil {
			return nil, err
		}
		cs.blockingConsumer = newStatusConsumer(bc, message.StatusBlocking, conf.Blocking, handlers.blockingHandler)
	}
	if conf.RetryingEnable {
		rc, err := client.subscribeByStatus(conf, message.StatusRetrying)
		if err != nil {
			return nil, err
		}
		cs.retryingConsumer = newStatusConsumer(rc, message.StatusRetrying, conf.Retrying, handlers.retryingHandler)
	}
	// start to listen message from all status multiStatusConsumer
	go cs.retrieveStatusMessages()
	return cs, nil
}

func (c *multiStatusConsumer) retrieveStatusMessages() {
	statusConsumers := []*statusConsumer{c.mainConsumer}
	if c.pendingConsumer != nil {
		statusConsumers = append(statusConsumers, c.pendingConsumer)
	}
	if c.blockingConsumer != nil {
		statusConsumers = append(statusConsumers, c.blockingConsumer)
	}
	if c.retryingConsumer != nil {
		statusConsumers = append(statusConsumers, c.retryingConsumer)
	}
	weights := make([]uint, len(statusConsumers))
	for index, c := range statusConsumers {
		weights[index] = c.policy.ConsumeWeight
	}
	consumeStrategy, err := config.BuildStrategy(c.statusStrategy, weights)
	if err != nil {
		panic(err)
	}
	for {
		var message ConsumerMessage
		messageChan := statusConsumers[consumeStrategy.Next()].StatusChan()
		select {
		case msg, ok := <-messageChan:
			if !ok {
				c.logger.Info("multiStatusConsumeFacade chan is closed")
				break
			}
			message = msg
		default:
			continue
		}
		// 获取到消息
		if message.Message != nil && message.Consumer != nil {
			fmt.Printf("received message  msgId: %v -- content: '%s'\n", message.ID(), string(message.Payload()))
			message.LeveledMessage = &leveledMessage{level: c.level}
			c.messageCh <- message
		} else {
			panic(fmt.Sprintf("consumed an invalid message: %v", message))
		}
	}
}

// Chan returns a channel to consume messages from
func (c *multiStatusConsumer) Chan() <-chan ConsumerMessage {
	return c.messageCh
}

func (c *multiStatusConsumer) Close() {
	c.mainConsumer.Close()
	close(c.messageCh)
	if c.blockingConsumer != nil {
		c.blockingConsumer.Close()
	}
	if c.pendingConsumer != nil {
		c.pendingConsumer.Close()
	}
	if c.retryingConsumer != nil {
		c.retryingConsumer.Close()
	}
}
