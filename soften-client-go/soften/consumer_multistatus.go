package soften

import (
	"fmt"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
	"github.com/shenqianjin/soften-client-go/soften/config"
	"github.com/shenqianjin/soften-client-go/soften/internal"
	"github.com/shenqianjin/soften-client-go/soften/message"
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
	if mainConsumer, err := cs.internalSubscribe(client, conf, message.StatusReady); err != nil {
		return nil, err
	} else {
		cs.mainConsumer = newStatusConsumer(mainConsumer, message.StatusReady, conf.Ready, nil)
	}
	if conf.PendingEnable {
		if pendingConsumer, err := cs.internalSubscribe(client, conf, message.StatusPending); err != nil {
			return nil, err
		} else {
			cs.pendingConsumer = newStatusConsumer(pendingConsumer, message.StatusPending, conf.Pending, handlers.pendingHandler)
		}
	}
	if conf.BlockingEnable {
		if blockingConsumer, err := cs.internalSubscribe(client, conf, message.StatusBlocking); err != nil {
			return nil, err
		} else {
			cs.blockingConsumer = newStatusConsumer(blockingConsumer, message.StatusBlocking, conf.Blocking, handlers.blockingHandler)
		}
	}
	if conf.RetryingEnable {
		if retryingConsumer, err := cs.internalSubscribe(client, conf, message.StatusRetrying); err != nil {
			return nil, err
		} else {
			cs.retryingConsumer = newStatusConsumer(retryingConsumer, message.StatusRetrying, conf.Retrying, handlers.retryingHandler)
		}
	}
	// start to listen message from all status multiStatusConsumer
	go cs.retrieveStatusMessages()
	return cs, nil
}

func (msc *multiStatusConsumer) retrieveStatusMessages() {
	statusConsumers := []*statusConsumer{msc.mainConsumer}
	if msc.pendingConsumer != nil {
		statusConsumers = append(statusConsumers, msc.pendingConsumer)
	}
	if msc.blockingConsumer != nil {
		statusConsumers = append(statusConsumers, msc.blockingConsumer)
	}
	if msc.retryingConsumer != nil {
		statusConsumers = append(statusConsumers, msc.retryingConsumer)
	}
	weights := make([]uint, len(statusConsumers))
	for index, c := range statusConsumers {
		weights[index] = c.policy.ConsumeWeight
	}
	consumeStrategy, err := config.BuildStrategy(msc.statusStrategy, weights)
	if err != nil {
		panic(err)
	}
	for {
		var msg ConsumerMessage
		messageChan := statusConsumers[consumeStrategy.Next()].StatusChan()
		select {
		case statusMsg, ok := <-messageChan:
			if !ok {
				msc.logger.Info("multiStatusConsumeFacade chan is closed")
				break
			}
			msg = statusMsg
		default:
			continue
		}
		// 获取到消息
		if msg.Message != nil && msg.Consumer != nil {
			fmt.Printf("received msg  msgId: %v -- content: '%s'\n", msg.ID(), string(msg.Payload()))
			msg.LeveledMessage = &leveledMessage{level: msc.level}
			msc.messageCh <- msg
		} else {
			panic(fmt.Sprintf("consumed an invalid msg: %v", msg))
		}
	}
}

// Chan returns a channel to consume messages from
func (msc *multiStatusConsumer) Chan() <-chan ConsumerMessage {
	return msc.messageCh
}

func (msc *multiStatusConsumer) Close() {
	msc.mainConsumer.Close()
	close(msc.messageCh)
	if msc.blockingConsumer != nil {
		msc.blockingConsumer.Close()
	}
	if msc.pendingConsumer != nil {
		msc.pendingConsumer.Close()
	}
	if msc.retryingConsumer != nil {
		msc.retryingConsumer.Close()
	}
}

func (msc *multiStatusConsumer) internalSubscribe(cli *client, conf *config.ConsumerConfig, status internal.MessageStatus) (pulsar.Consumer, error) {
	suffix, err := message.TopicSuffixOf(status)
	if err != nil {
		return nil, err
	}
	topic := conf.Topics[0] + suffix
	subscriptionName := conf.SubscriptionName + suffix
	consumerOption := pulsar.ConsumerOptions{
		Topic:                       topic,
		SubscriptionName:            subscriptionName,
		Type:                        conf.Type,
		SubscriptionInitialPosition: conf.SubscriptionInitialPosition,
		NackRedeliveryDelay:         conf.NackRedeliveryDelay,
		NackBackoffPolicy:           conf.NackBackoffPolicy,
		MessageChannel:              nil,
	}
	if conf.DLQ != nil {
		consumerOption.DLQ = &pulsar.DLQPolicy{
			MaxDeliveries:    conf.DLQ.MaxDeliveries,
			RetryLetterTopic: conf.DLQ.RetryLetterTopic,
			DeadLetterTopic:  conf.DLQ.DeadLetterTopic,
		}
	}
	// only main status need compatible with pulsar retry enable and multi-topics
	if status == message.StatusReady {
		consumerOption.Topics = conf.Topics
		consumerOption.RetryEnable = conf.RetryEnable
	} else {
		consumerOption.Topics = nil
		consumerOption.RetryEnable = false
	}
	pulsarConsumer, err := cli.Client.Subscribe(consumerOption)
	return pulsarConsumer, err
}
