package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/bmizerany/perks/quantile"
	"github.com/shenqianjin/soften-client-go/perf/internal"
	"github.com/shenqianjin/soften-client-go/soften"
	"github.com/shenqianjin/soften-client-go/soften/checker"
	"github.com/shenqianjin/soften-client-go/soften/config"
	"github.com/shenqianjin/soften-client-go/soften/message"
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

	handleDoneWeight     uint
	handleRetryingEnable bool
	handleRetryingWeight uint
	handlePendingEnable  bool
	handlePendingWeight  uint
	handleBlockingEnable bool
	handleBlockingWeight uint
	handleDeadEnable     bool
	handleDeadWeight     uint
	handleDiscardEnable  bool
	handleDiscardWeight  uint

	Concurrency         uint
	Limits              []uint64 // 每秒限制 //
	RadicalConcurrences []uint64 // 每秒限制 // radical粒度并发限制
}

type consumer struct {
	clientArgs   *clientArgs
	consumerArgs *consumeArgs

	choicePolicy internal.GotoPolicy
	costPolicy   internal.CostPolicy

	consumeStatCh chan *consumeStat

	//radicalLimiters     map[string]*rate.RateLimiter
	concurrencyLimiters map[string]internal.ConcurrencyLimiter
}

type consumeStat struct {
	bytes           int64
	receivedLatency float64
	finishedLatency float64
	handledLatency  float64

	radicalKey string
	//radicalFinishedLatencies map[string]float64
	//radicalHandledLatencies  map[string]float64
}

func newConsumer(clientArgs *clientArgs, consumerArgs *consumeArgs) *consumer {

	weightMap := map[string]uint{string(message.GotoDone): consumerArgs.handleDoneWeight}
	if consumerArgs.handleRetryingEnable && consumerArgs.handleRetryingWeight > 0 {
		weightMap[string(message.GotoRetrying)] = consumerArgs.handleRetryingWeight
	}
	if consumerArgs.handlePendingEnable && consumerArgs.handlePendingWeight > 0 {
		weightMap[string(message.GotoPending)] = consumerArgs.handlePendingWeight
	}
	if consumerArgs.handleBlockingEnable && consumerArgs.handleBlockingWeight > 0 {
		weightMap[string(message.GotoBlocking)] = consumerArgs.handleBlockingWeight
	}
	if consumerArgs.handleDeadEnable && consumerArgs.handleDeadWeight > 0 {
		weightMap[string(message.GotoDead)] = consumerArgs.handleDeadWeight
	}
	if consumerArgs.handleDiscardEnable && consumerArgs.handleDiscardWeight > 0 {
		weightMap[string(message.GotoDiscard)] = consumerArgs.handleDiscardWeight
	}
	consumeChoice := internal.NewRoundRandWeightGotoPolicy(weightMap)

	// Retry to create producer indefinitely
	c := &consumer{
		clientArgs:          clientArgs,
		consumerArgs:        consumerArgs,
		choicePolicy:        consumeChoice,
		consumeStatCh:       make(chan *consumeStat),
		concurrencyLimiters: make(map[string]internal.ConcurrencyLimiter),
		//radicalLimiters:     make(map[string]*rate.RateLimiter),
	}
	// initialize cost policy
	if consumerArgs.costAverageInMs > 0 {
		c.costPolicy = internal.NewAvgCostPolicy(consumerArgs.costAverageInMs, consumerArgs.costPositiveJitter, consumerArgs.costNegativeJitter)
	}

	/*if len(consumerArgs.Limits) > 0 {
		for index, li := range consumerArgs.Limits {
			if li > 0 {
				c.radicalLimiters[fmt.Sprintf("Radical-%d", index)] = rate.New(int(li), time.Second)
			} else {
				c.radicalLimiters[fmt.Sprintf("Radical-%d", index)] = nil
			}
		}
	}*/

	if len(consumerArgs.RadicalConcurrences) > 0 {
		for index, con := range consumerArgs.RadicalConcurrences {
			if con > 0 {
				c.concurrencyLimiters[fmt.Sprintf("Radical-%d", index)] = internal.NewConcurrencyLimiter(int(con))
			} else {
				c.concurrencyLimiters[fmt.Sprintf("Radical-%d", index)] = nil
			}
		}
	}

	return c
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

	// create consumer
	listener, err := client.CreateListener(config.ConsumerConfig{
		Topic:            c.consumerArgs.Topic,
		SubscriptionName: c.consumerArgs.SubscriptionName,
		Concurrency:      &config.ConcurrencyPolicy{CorePoolSize: c.consumerArgs.Concurrency},
		PendingEnable:    true,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	// start monitoring: async
	go c.stats(stop, c.consumeStatCh)

	// start message listener
	err = listener.StartPremium(context.Background(), c.internalHandle, checker.PrePendingChecker(c.internalPrePendingChecker))
	if err != nil {
		log.Fatal(err)
	}
}

func (c *consumer) internalPrePendingChecker(cm pulsar.Message) checker.CheckStatus {
	if radicalKey, ok := cm.Properties()["Radical"]; ok {
		if radicalKey == "Radical-0" {
			return checker.CheckStatusRejected
		}
		if limiter, ok2 := c.concurrencyLimiters[radicalKey]; ok2 && limiter != nil {
			if !limiter.TryAcquire() {
				return checker.CheckStatusPassed
			} else {
				//time := time.Now().Format(time.RFC3339Nano)
				//log.Info("acquire .................", time)
				//defer limiter.Release()
				return checker.CheckStatusRejected.WithHandledDefer(func() {
					//log.Info("release .................", time)
					limiter.Release()
				})
			}
		}
	}
	return checker.CheckStatusRejected
}

var RFC3339TimeInSecondPattern = "20060102150405.999"

func (c *consumer) internalHandle(cm pulsar.Message) soften.HandleStatus {
	originPublishTime := cm.PublishTime()
	if originalPublishTime, ok := cm.Properties()[message.XPropertyOriginPublishTime]; ok {
		if parsedTime, err := time.Parse(RFC3339TimeInSecondPattern, originalPublishTime); err == nil {
			//now := time.Now()
			//log.Info("handle pending message..................", now)
			originPublishTime = parsedTime
		}
	}
	start := time.Now()
	stat := &consumeStat{
		bytes:           int64(len(cm.Payload())),
		receivedLatency: time.Since(cm.PublishTime()).Seconds(),
		//radicalHandledLatencies:  make(map[string]float64, len(c.consumerArgs.RadicalConcurrences)),
		//radicalFinishedLatencies: make(map[string]float64, len(c.consumerArgs.RadicalConcurrences)),
	}

	/*if radicalKey, ok := cm.Properties()["Radical"]; ok {
		if limiter, ok2 := c.radicalLimiters[radicalKey]; ok2 && limiter != nil {
			limiter.Wait()
		}
	}*/

	if c.consumerArgs.costAverageInMs > 0 {
		time.Sleep(c.costPolicy.Next()) // 模拟业务处理
	}
	result := soften.HandleStatusOk
	//n := c.choicePolicy.Next()
	n := string(message.GotoDone)
	switch n {
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
		result = soften.HandleStatusOk
		stat.finishedLatency = time.Since(originPublishTime).Seconds() // 从消息产生到处理完成的时间(中间状态不是完成状态)
		if radicalKey, ok := cm.Properties()["Radical"]; ok {
			stat.radicalKey = radicalKey
			//stat.radicalFinishedLatencies[radicalKey] = stat.finishedLatency
		}
	}

	stat.handledLatency = time.Since(start).Seconds()
	if radicalKey, ok := cm.Properties()["Radical"]; ok {
		stat.radicalKey = radicalKey
		//stat.radicalHandledLatencies[radicalKey] = stat.handledLatency
	}
	c.consumeStatCh <- stat
	return result
}

/*func (c *consumer) internalConsume(realConsumer pulsar.Consumer, stop <-chan struct{}, consumeStatCh chan<- *consumeStat) {

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

			if c.consumerArgs.costAverageInMs > 0 {
				time.Sleep(c.costPolicy.Next()) // 模拟业务处理
			}

			realConsumer.Ack(cm.Message)

			stat.finishedLatency = time.Since(cm.PublishTime()).Seconds() // 从消息产生到处理完成的时间(中间状态不是完成状态)
			stat.handledLatency = time.Since(start).Seconds()
			consumeStatCh <- stat
		case <-stop:
			return
		}
	}
}*/

func (c *consumer) stats(stop <-chan struct{}, consumeStatCh <-chan *consumeStat) {
	// Print stats of the perfConsume rate
	tick := time.NewTicker(10 * time.Second)
	defer tick.Stop()
	receivedQ := quantile.NewTargeted(0.50, 0.95, 0.99, 0.999, 1.0)
	finishedQ := quantile.NewTargeted(0.50, 0.95, 0.99, 0.999, 1.0)
	handledQ := quantile.NewTargeted(0.50, 0.95, 0.99, 0.999, 1.0)
	msgHandled := int64(0)
	bytesHandled := int64(0)
	radicalHandleMsg := make(map[string]int64, len(c.consumerArgs.RadicalConcurrences))
	radicalHandleQ := make(map[string]*quantile.Stream, len(c.consumerArgs.RadicalConcurrences))
	radicalFinishedQ := make(map[string]*quantile.Stream, len(c.consumerArgs.RadicalConcurrences))

	for {
		select {
		case <-stop:
			log.Infof("Closing consume stats printer")
			return
		case <-tick.C:
			currentMsgReceived := atomic.SwapInt64(&msgHandled, 0)
			currentBytesReceived := atomic.SwapInt64(&bytesHandled, 0)
			msgRate := float64(currentMsgReceived) / float64(10)
			bytesRate := float64(currentBytesReceived) / float64(10)

			statB := &bytes.Buffer{}
			fmt.Fprintf(statB, `<<<<<<<<<<
			Summary - Consume rate: %6.1f msg/s - %6.1f Mbps - 
				Received Latency ms: 50%% %5.1f - 95%% %5.1f - 99%% %5.1f - 99.9%% %5.1f - max %6.1f  
				Finished Latency ms: 50%% %5.1f - 95%% %5.1f - 99%% %5.1f - 99.9%% %5.1f - max %6.1f
				Handled  Latency ms: 50%% %5.1f - 95%% %5.1f - 99%% %5.1f - 99.9%% %5.1f - max %6.1f`,
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

				handledQ.Query(0.5)*1000,
				handledQ.Query(0.95)*1000,
				handledQ.Query(0.99)*1000,
				handledQ.Query(0.999)*1000,
				handledQ.Query(1.0)*1000,
			)
			if len(radicalHandleMsg) > 0 {
				fmt.Fprintf(statB, `
			Detail >> `)
			}
			for key, v := range radicalHandleMsg {
				fmt.Fprintf(statB, "%s rate: %6.1f msg/s - ", key, float64(v)/float64(10))
				radicalHandleMsg[key] = 0
			}
			for key, q := range radicalFinishedQ {
				fmt.Fprintf(statB, `
				  %s Finished Latency ms: 50%% %5.1f - 95%% %5.1f - 99%% %5.1f - 99.9%% %5.1f - max %6.1f`, key,
					q.Query(0.5)*1000, q.Query(0.95)*1000, q.Query(0.99)*1000, q.Query(0.999)*1000, q.Query(1.0)*1000)
			}
			for key, q := range radicalHandleQ {
				fmt.Fprintf(statB, `
				  %s Handled  Latency ms: 50%% %5.1f - 95%% %5.1f - 99%% %5.1f - 99.9%% %5.1f - max %6.1f`, key,
					q.Query(0.5)*1000, q.Query(0.95)*1000, q.Query(0.99)*1000, q.Query(0.999)*1000, q.Query(1.0)*1000)
			}
			log.Info(statB.String())

			receivedQ.Reset()
			finishedQ.Reset()
			handledQ.Reset()
			for _, q := range radicalHandleQ {
				q.Reset()
			}
			for _, q := range radicalFinishedQ {
				q.Reset()
			}
			//messagesConsumed = 0
		case stat := <-consumeStatCh:
			msgHandled++
			bytesHandled += stat.bytes
			receivedQ.Insert(stat.receivedLatency)
			finishedQ.Insert(stat.finishedLatency)
			handledQ.Insert(stat.handledLatency)
			radicalHandleMsg[stat.radicalKey]++
			// handle
			if _, ok := radicalHandleQ[stat.radicalKey]; !ok {
				radicalHandleQ[stat.radicalKey] = quantile.NewTargeted(0.50, 0.95, 0.99, 0.999, 1.0)
			}
			radicalHandleQ[stat.radicalKey].Insert(stat.handledLatency)
			// finish
			if _, ok := radicalFinishedQ[stat.radicalKey]; !ok {
				radicalFinishedQ[stat.radicalKey] = quantile.NewTargeted(0.50, 0.95, 0.99, 0.999, 1.0)
			}
			radicalFinishedQ[stat.radicalKey].Insert(stat.finishedLatency)
		}
	}
}
