package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/ratelimit"

	"github.com/shenqianjin/soften-client-go/soften/config"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/bmizerany/perks/quantile"
	log "github.com/sirupsen/logrus"
)

// produceArgs define the parameters required by perfProduce
type produceArgs struct {
	Topic              string
	Rates              []uint64 // ordered produce rates: [normal, radical 1, radical 2, ..., radical N]
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
	realProducer, err := client.CreateProducer(config.ProducerConfig{
		Topic:              p.produceArgs.Topic,
		MaxPendingMessages: p.produceArgs.ProducerQueueSize,
		//BatchingMaxPublishDelay: time.Millisecond * time.Duration(p.produceArgs.BatchingTimeMillis),
		BatchingMaxSize: uint(p.produceArgs.BatchingMaxSize * 1024),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer realProducer.Close()

	ch := make(chan *produceStat)
	// start monitoring: async
	go p.stats(stopCh, ch)
	// start perfProduce: sync to hang
	// radical
	if len(p.produceArgs.Rates) > 1 {
		for index, r := range p.produceArgs.Rates[1:] {
			go p.internalProduce(realProducer, r, fmt.Sprintf("Radical-%d", index+1), stopCh, ch)
		}
	}
	normalRate := uint64(0)
	if len(p.produceArgs.Rates) > 0 {
		normalRate = p.produceArgs.Rates[0]

	}
	// normal
	p.internalProduce(realProducer, normalRate, "Radical-0", stopCh, ch)
}

func (p *producer) internalProduce(realProducer pulsar.Producer, r uint64, radical string, stopCh <-chan struct{}, ch chan<- *produceStat) {
	ctx := context.Background()
	payload := make([]byte, p.produceArgs.MessageSize)
	//var rateLimiter *rate.RateLimiter
	var rateLimiter ratelimit.Limiter

	if r > 0 {
		//rateLimiter = rate.New(int(r), time.Second)
		//rateLimiter = rate.New(int(r)/10, 100*time.Millisecond)
		rateLimiter = ratelimit.New(int(r), ratelimit.Per(time.Second))
	}
	/*var radicalChoicePolicy internal.ChoicePolicy
	if p.produceArgs.RadicalRate > 0 {
		radicalChoicePolicy = internal.NewRateChoicePolicy(p.produceArgs.RadicalRate)
	}*/
	for {
		select {
		case <-stopCh:
			log.Infof("Closing produce stats printer")
			return
		default:
		}

		if rateLimiter != nil {
			//rateLimiter.Wait()
			rateLimiter.Take()
		}

		msg := &pulsar.ProducerMessage{
			Payload: payload,
		}

		if radical != "" {
			msg.Properties = map[string]string{"Radical": radical}
		} /* else if radicalChoicePolicy != nil {
			choice := radicalChoicePolicy.Next()
			if choice == 1 { // 激进模式
				msg.Properties = map[string]string{"Radical": "Radical-1"}
			}
		}*/

		start := time.Now()
		realProducer.SendAsync(ctx, msg, func(msgID pulsar.MessageID, message *pulsar.ProducerMessage, e error) {
			if e != nil {
				log.WithError(e).Fatal("Failed to publish")
			}

			latency := time.Since(start).Seconds()
			stat := &produceStat{latency: latency}
			if radicalKey, ok := msg.Properties["Radical"]; ok {
				stat.radicalKey = radicalKey
			}
			ch <- stat
		})
	}
}

func (p *producer) stats(stop <-chan struct{}, statCh <-chan *produceStat) {
	// Print stats of the publish rate and latencies
	tick := time.NewTicker(10 * time.Second)
	defer tick.Stop()
	q := quantile.NewTargeted(0.50, 0.95, 0.99, 0.999, 1.0)
	messagesPublished := 0
	radicalMsgPublished := make(map[string]int64, len(p.produceArgs.Rates))

	for {
		select {
		case <-stop:
			return
		case <-tick.C:
			messageRate := float64(messagesPublished) / float64(10)

			statB := &bytes.Buffer{}
			fmt.Fprintf(statB, `>>>>>>>>>>
		Stats - Publish rate: %6.1f msg/s - %6.1f Mbps - 
				Finished Latency ms: 50%% %5.1f - 95%% %5.1f - 99%% %5.1f - 99.9%% %5.1f - max %6.1f`,
				messageRate,
				messageRate*float64(p.produceArgs.MessageSize)/1024/1024*8,
				q.Query(0.5)*1000,
				q.Query(0.95)*1000,
				q.Query(0.99)*1000,
				q.Query(0.999)*1000,
				q.Query(1.0)*1000,
			)
			if len(radicalMsgPublished) > 0 {
				fmt.Fprintf(statB, `
			Detail >> `)
			}
			for key, v := range radicalMsgPublished {
				fmt.Fprintf(statB, "%s rate: %6.1f msg/s - ", key, float64(v)/float64(10))
				radicalMsgPublished[key] = 0
			}
			log.Info(statB.String())
			q.Reset()
			messagesPublished = 0
		case stat := <-statCh:
			messagesPublished++
			radicalMsgPublished[stat.radicalKey]++
			q.Insert(stat.latency)
		}
	}
}

type produceStat struct {
	latency float64

	radicalKey string
}
