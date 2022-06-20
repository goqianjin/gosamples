package main

import (
	"encoding/json"
	"sync/atomic"
	"time"

	"github.com/shenqianjin/soften-client-go/soften/message"

	"github.com/shenqianjin/soften-client-go/soften/config"

	"github.com/shenqianjin/soften-client-go/soften"

	"github.com/shenqianjin/soften-client-go/perf/internal"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/bmizerany/perks/quantile"
	log "github.com/sirupsen/logrus"
)

// consumeArgs define the parameters required by perfConsume
type consumeArgs struct {
	Topic             string
	SubscriptionName  string
	ReceiverQueueSize int

	costAverageInMs    int64
	costPositiveJitter float64
	costNegativeJitter float64

	doneWeight     uint
	retryingEnable bool
	retryingWeight uint
	pendingEnable  bool
	pendingWeight  uint
	blockingEnable bool
	blockingWeight uint
	deadEnable     bool
	deadWeight     uint
	discardEnable  bool
	discardWeight  uint
}

type consumer struct {
	clientArgs   *clientArgs
	consumerArgs *consumeArgs

	choicePolicy internal.GotoPolicy
	costPolicy   internal.CostPolicy

	consumeStatCh chan *consumeStat
}

type consumeStat struct {
	bytes           int64
	receivedLatency float64
	finishedLatency float64
	consumedLatency float64
}

func newConsumer(clientArgs *clientArgs, consumerArgs *consumeArgs) *consumer {

	weightMap := map[string]uint{string(message.GotoDone): consumerArgs.doneWeight}
	if consumerArgs.retryingEnable && consumerArgs.retryingWeight > 0 {
		weightMap[string(message.GotoRetrying)] = consumerArgs.retryingWeight
	}
	if consumerArgs.pendingEnable && consumerArgs.pendingWeight > 0 {
		weightMap[string(message.GotoPending)] = consumerArgs.pendingWeight
	}
	if consumerArgs.blockingEnable && consumerArgs.blockingWeight > 0 {
		weightMap[string(message.GotoBlocking)] = consumerArgs.blockingWeight
	}
	if consumerArgs.deadEnable && consumerArgs.deadWeight > 0 {
		weightMap[string(message.GotoDead)] = consumerArgs.deadWeight
	}
	if consumerArgs.discardEnable && consumerArgs.discardWeight > 0 {
		weightMap[string(message.GotoDiscard)] = consumerArgs.discardWeight
	}
	consumeChoice := internal.NewRoundRandWeightGotoPolicy(weightMap)

	// Retry to create producer indefinitely
	costPolicy := internal.NewAvgCostPolicy(consumerArgs.costAverageInMs, consumerArgs.costPositiveJitter, consumerArgs.costNegativeJitter)

	return &consumer{
		clientArgs:    clientArgs,
		consumerArgs:  consumerArgs,
		choicePolicy:  consumeChoice,
		costPolicy:    costPolicy,
		consumeStatCh: make(chan *consumeStat),
	}
}

func (c *consumer) perfConsume(stop <-chan struct{}) {
	b, _ := json.MarshalIndent(c.clientArgs, "", "  ")
	log.Info("Client config: ", string(b))
	b, _ = json.MarshalIndent(c.consumerArgs, "", "  ")
	log.Info("Consumer config: ", string(b))
	// create client
	client, err := newClient(c.clientArgs)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// start monitoring: async
	go c.stats(stop, c.consumeStatCh)

	// create consumer
	realConsumer, err := client.SubscribePremium(config.ConsumerConfig{
		Topic:            c.consumerArgs.Topic,
		SubscriptionName: c.consumerArgs.SubscriptionName,
	}, c.internalHandle)
	if err != nil {
		log.Fatal(err)
	}
	defer realConsumer.Close()

	time.Sleep(2 * time.Minute)

	// start perfConsume: sync to hang
	//c.internalConsume(realConsumer, stop, c.consumeStatCh)
}

func (c *consumer) internalHandle(cm pulsar.Message) soften.HandleStatus {
	start := time.Now()
	stat := &consumeStat{
		bytes:           int64(len(cm.Payload())),
		receivedLatency: time.Since(cm.PublishTime()).Seconds(),
	}

	if c.consumerArgs.costAverageInMs != 0 {
		time.Sleep(c.costPolicy.Next()) // 模拟业务处理
	}
	result := soften.HandleStatusOk
	switch c.choicePolicy.Next() {
	case string(message.GotoRetrying):
		result = soften.HandleStatusDefault.GotoAction(message.GotoRetrying)
	case string(message.GotoPending):
		result = soften.HandleStatusDefault.GotoAction(message.GotoPending)
	case string(message.GotoBlocking):
		result = soften.HandleStatusDefault.GotoAction(message.GotoBlocking)
	case string(message.GotoDead):
		result = soften.HandleStatusDefault.GotoAction(message.GotoDead)
	case string(message.GotoDiscard):
		result = soften.HandleStatusDefault.GotoAction(message.GotoDiscard)
	default:
		stat.finishedLatency = time.Since(cm.PublishTime()).Seconds() // 从消息产生到处理完成的时间(中间状态不是完成状态)
		result = soften.HandleStatusOk
	}

	stat.consumedLatency = time.Since(start).Seconds()
	c.consumeStatCh <- stat
	return result
}

func (c *consumer) internalConsume(realConsumer pulsar.Consumer, stop <-chan struct{}, consumeStatCh chan<- *consumeStat) {

	for {
		select {
		case cm, ok := <-realConsumer.Chan():
			if !ok {
				return
			}
			start := time.Now()
			stat := &consumeStat{
				bytes:           int64(len(cm.Message.Payload())),
				receivedLatency: time.Since(cm.PublishTime()).Seconds(),
			}

			if c.consumerArgs.costAverageInMs != 0 {
				time.Sleep(c.costPolicy.Next()) // 模拟业务处理
			}

			realConsumer.Ack(cm.Message)

			stat.finishedLatency = time.Since(cm.PublishTime()).Seconds() // 从消息产生到处理完成的时间(中间状态不是完成状态)
			stat.consumedLatency = time.Since(start).Seconds()
			consumeStatCh <- stat
		case <-stop:
			return
		}
	}
}

func (c *consumer) stats(stop <-chan struct{}, consumeStatCh <-chan *consumeStat) {
	// Print stats of the perfConsume rate
	tick := time.NewTicker(10 * time.Second)
	defer tick.Stop()
	receivedQ := quantile.NewTargeted(0.50, 0.95, 0.99, 0.999, 1.0)
	finishedQ := quantile.NewTargeted(0.50, 0.95, 0.99, 0.999, 1.0)
	consumedQ := quantile.NewTargeted(0.50, 0.95, 0.99, 0.999, 1.0)
	msgReceived := int64(0)
	bytesReceived := int64(0)

	for {
		select {
		case <-stop:
			return
		case <-tick.C:
			currentMsgReceived := atomic.SwapInt64(&msgReceived, 0)
			currentBytesReceived := atomic.SwapInt64(&bytesReceived, 0)
			msgRate := float64(currentMsgReceived) / float64(10)
			bytesRate := float64(currentBytesReceived) / float64(10)

			log.Infof(`Stats - Consume rate: %6.1f msg/s - %6.1f Mbps - 
				Received Latency ms: 50%% %5.1f -95%% %5.1f - 99%% %5.1f - 99.9%% %5.1f - max %6.1f  
				Finished Latency ms: 50%% %5.1f -95%% %5.1f - 99%% %5.1f - 99.9%% %5.1f - max %6.1f
				Comsumed Latency ms: 50%% %5.1f -95%% %5.1f - 99%% %5.1f - 99.9%% %5.1f - max %6.1f`,
				msgRate, bytesRate*8/1024/1024,

				receivedQ.Query(0.5)*1000,
				receivedQ.Query(0.95)*1000,
				receivedQ.Query(0.99)*1000,
				receivedQ.Query(0.999)*1000,
				receivedQ.Query(1.0)*1000,

				finishedQ.Query(0.5)*1000,
				finishedQ.Query(0.95)*1000,
				finishedQ.Query(0.99)*1000,
				finishedQ.Query(0.999)*1000,
				finishedQ.Query(1.0)*1000,

				consumedQ.Query(0.5)*1000,
				consumedQ.Query(0.95)*1000,
				consumedQ.Query(0.99)*1000,
				consumedQ.Query(0.999)*1000,
				consumedQ.Query(1.0)*1000,
			)
			receivedQ.Reset()
			finishedQ.Reset()
			consumedQ.Reset()
			//messagesConsumed = 0
		case stat := <-consumeStatCh:
			msgReceived++
			bytesReceived += stat.bytes
			receivedQ.Insert(stat.receivedLatency)
			finishedQ.Insert(stat.finishedLatency)
			consumedQ.Insert(stat.consumedLatency)
		}
	}
}
