package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/shenqianjin/soften-client-go/perf/internal"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/beefsack/go-rate"
	"github.com/bmizerany/perks/quantile"
	log "github.com/sirupsen/logrus"
)

// produceArgs define the parameters required by perfProduce
type produceArgs struct {
	Topic              string
	Rate               int
	BatchingTimeMillis int
	BatchingMaxSize    int
	MessageSize        int
	ProducerQueueSize  int
	RadicalRate        float64 // 激进消息比例
}

type producer struct {
	clientArgs  *clientArgs
	produceArgs *produceArgs
}

func newProducer(cliArgs *clientArgs, pArgs *produceArgs) *producer {
	return &producer{clientArgs: cliArgs, produceArgs: pArgs}
}

func (p *producer) perfProduce(stopCh <-chan struct{}) {
	b, _ := json.MarshalIndent(p.clientArgs, "", "  ")
	log.Info("Client config: ", string(b))
	b, _ = json.MarshalIndent(p.produceArgs, "", "  ")
	log.Info("Producer config: ", string(b))
	// create client
	client, err := newClient(p.clientArgs)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	// create producer
	realProducer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic:                   p.produceArgs.Topic,
		MaxPendingMessages:      p.produceArgs.ProducerQueueSize,
		BatchingMaxPublishDelay: time.Millisecond * time.Duration(p.produceArgs.BatchingTimeMillis),
		BatchingMaxSize:         uint(p.produceArgs.BatchingMaxSize * 1024),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer realProducer.Close()

	ch := make(chan float64)
	// start monitoring: async
	go p.stats(stopCh, ch)
	// start perfProduce: sync to hang
	p.internalProduce(realProducer, stopCh, ch)
}

func (p *producer) internalProduce(realProducer pulsar.Producer, stopCh <-chan struct{}, ch chan<- float64) {
	ctx := context.Background()
	payload := make([]byte, p.produceArgs.MessageSize)
	var rateLimiter *rate.RateLimiter
	if p.produceArgs.Rate > 0 {
		rateLimiter = rate.New(p.produceArgs.Rate, time.Second)
	}
	radicalChoicePolicy := internal.NewRateChoicePolicy(p.produceArgs.RadicalRate)
	for {
		select {
		case <-stopCh:
			return
		default:
		}

		if rateLimiter != nil {
			rateLimiter.Wait()
		}
		choice := radicalChoicePolicy.Next()
		msg := &pulsar.ProducerMessage{
			Payload: payload,
		}
		if choice == 1 { // 激进模式
			msg.Properties = map[string]string{"RadicalFlag": "true"}
		}
		start := time.Now()

		realProducer.SendAsync(ctx, msg, func(msgID pulsar.MessageID, message *pulsar.ProducerMessage, e error) {
			if e != nil {
				log.WithError(e).Fatal("Failed to publish")
			}

			latency := time.Since(start).Seconds()
			ch <- latency
		})
	}
}

func (p *producer) stats(stop <-chan struct{}, latencyCh <-chan float64) {
	// Print stats of the publish rate and latencies
	tick := time.NewTicker(10 * time.Second)
	defer tick.Stop()
	q := quantile.NewTargeted(0.50, 0.95, 0.99, 0.999, 1.0)
	messagesPublished := 0

	for {
		select {
		case <-stop:
			return
		case <-tick.C:
			messageRate := float64(messagesPublished) / float64(10)
			log.Infof(`Stats - Publish rate: %6.1f msg/s - %6.1f Mbps - 
				Finished Latency ms: 50%% %5.1f -95%% %5.1f - 99%% %5.1f - 99.9%% %5.1f - max %6.1f`,
				messageRate,
				messageRate*float64(p.produceArgs.MessageSize)/1024/1024*8,
				q.Query(0.5)*1000,
				q.Query(0.95)*1000,
				q.Query(0.99)*1000,
				q.Query(0.999)*1000,
				q.Query(1.0)*1000,
			)

			q.Reset()
			messagesPublished = 0
		case latency := <-latencyCh:
			messagesPublished++
			q.Insert(latency)
		}
	}
}
