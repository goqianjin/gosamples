package soam

import (
	"fmt"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
)

type consumer struct {
	pulsar.Consumer
	client    *client
	config    ComsumerConfig
	log       log.Logger
	messageCh chan pulsar.ConsumerMessage // channel used to deliver message to application

	consumeStrategy  ConsumeStrategy // 消费策略
	mainConsumer     *statusConsumer
	pendingConsumer  *statusConsumer
	blockingConsumer *statusConsumer
	retryingConsumer *statusConsumer
	checkers         *consumeCheckers
	handlers         *consumeHandlers
}

func newConsumer(c *client, config ComsumerConfig, checkpointMap map[CheckType]*checkpoint) (*consumer, error) {
	consumer := &consumer{
		client: c,
		config: config,
	}
	// create status consumers
	if mainConsumer, err := c.subscribeByStatus(config, MessageStatusReady); err != nil {
		return nil, err
	} else {
		consumer.mainConsumer = newStatusConsumer(consumer, mainConsumer, MessageStatusReady, config.Ready)
	}
	if config.PendingEnable {
		pendingConsumer, err := c.subscribeByStatus(config, MessageStatusPending)
		if err != nil {
			return nil, err
		}
		consumer.pendingConsumer = newStatusConsumer(consumer, pendingConsumer, MessageStatusPending, config.Pending)
	}
	if config.BlockingEnable {
		bc, err := c.subscribeByStatus(config, MessageStatusBlocking)
		if err != nil {
			return nil, err
		}
		consumer.blockingConsumer = newStatusConsumer(consumer, bc, MessageStatusBlocking, config.Blocking)
	}
	if config.RetryEnable {
		rc, err := c.subscribeByStatus(config, MessageStatusRetrying)
		if err != nil {
			return nil, err
		}
		consumer.retryingConsumer = newStatusConsumer(consumer, rc, MessageStatusRetrying, config.Retrying)
	}
	// initialize checkers
	if checkers, err := newConsumeCheckers(consumer, config); err != nil {
		return nil, err
	} else {
		consumer.checkers = checkers
	}
	// initialize handlers
	if handlers, err := newConsumeHandlers(consumer, config); err != nil {
		return nil, err
	} else {
		consumer.handlers = handlers
	}

	return consumer, nil
}

func (c *consumer) start() {
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
	consumeStrategy := c.getConsumeStrategy(statusConsumers)
	for {
		var message pulsar.ConsumerMessage
		messageChan := statusConsumers[consumeStrategy.NextConsumer()].Chan()
		select {
		case msg := <-messageChan:
			message = msg
		default:
		}
		// 获取到消息
		if message.Message != nil && message.Consumer != nil {
			c.messageCh <- message
		} else {
			panic(fmt.Sprintf("consumed an invalid message: %v", message))
		}
	}
}

func (c *consumer) getConsumeStrategy(statusConsumers []*statusConsumer) consumeStrategy {
	if c.consumeStrategy == ConsumeStrategyRandEach {
		totalWeight := statusConsumers[0].policy.ConsumeWeight
		prefixSumWeights := []uint{totalWeight}
		for _, sc := range statusConsumers[1:] {
			totalWeight += sc.policy.ConsumeWeight
			prefixSumWeights = append(prefixSumWeights, totalWeight)
		}
		return &randEachConsumeStrategy{total: totalWeight, prefixSums: prefixSumWeights}
	} else { // c.consumeStrategy == ConsumeStrategyRandRound;
		totalWeight := 0
		count := 0
		indexConsumerMap := make(map[int]int)
		for index, sc := range statusConsumers {
			weight := int(sc.policy.ConsumeWeight)
			for i := 0; i < weight; i++ {
				indexConsumerMap[count] = index
				count++
			}
			totalWeight += weight
		}
		return &randRoundConsumeStrategy{total: totalWeight, indexConsumerMap: indexConsumerMap}
	}
}

func (c *consumer) consume(handler HandlerInPremium, message pulsar.ConsumerMessage) {
	// pre-check to route
	if routed := c.checkers.tryPreCheckInTurn(message, defaultPreCheckTypesInTurn...); routed {
		return
	}
	// do biz handle
	handleResult := handler(message)
	// consume successfully
	if handleResult.handledStatus == MessageStatusDone {
		message.Consumer.Ack(message.Message)
		return
	}
	// post-check to route - for obvious status
	//currentStatus := MessageParser.GetCurrentStatus(message)
	if handleResult.handledStatus != "" { //currentStatus {
		if checkType, ok := messageStatusToPostCheckTypeMap[handleResult.handledStatus]; ok {
			if ok2 := c.checkers.tryPostCheckInTurn(message, handleResult.err, checkType); ok2 {
				return
			}
		}
	}
	// post-check to route - for obvious checkers or configured checkers
	postCheckTypesInTurn := DefaultPostCheckTypesInTurn
	if len(handleResult.postCheckTypesInTurn) > 0 {
		postCheckTypesInTurn = handleResult.postCheckTypesInTurn
	}
	if routed := c.checkers.tryPostCheckInTurn(message, handleResult.err, postCheckTypesInTurn...); routed {
		return
	}
	// here means to let application client Ack/Nack message
	// metrics to record
	return
}
func (c *consumer) internalRouteByStatus(status messageStatus, message pulsar.ConsumerMessage) (routed bool) {
	switch status {
	case MessageStatusPending:
		return c.config.PendingEnable && c.handlers.pendingHandler.Handle(message)
	case MessageStatusBlocking:
		return c.config.BlockingEnable && c.handlers.blockingHandler.Handle(message)
	case MessageStatusRetrying:
		return c.config.RetryEnable && c.handlers.retryingHandler.Handle(message)
	case MessageStatusDead:
		return c.config.DeadEnable && c.handlers.deadHandler.Handle(message)
	}
	return false
}
