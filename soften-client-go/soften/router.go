package soften

import (
	"errors"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
	"github.com/shenqianjin/soften-client-go/soften/internal/backoff"
)

// ------ router ------

type routerOptions struct {
	Topic               string
	connectInSyncEnable bool
}

type router struct {
	pulsar.Producer
	options routerOptions
	client  pulsar.Client
	logger  log.Logger
	ready   bool
}

func newRouter(logger log.Logger, client pulsar.Client, options routerOptions) (*router, error) {
	if options.Topic == "" {
		return nil, errors.New("routerOptions.Topic needs to be set to a valid topic name")
	}
	r := &router{
		client:  client,
		options: options,
	}
	r.logger = logger.SubLogger(log.Fields{"route-topic": options.Topic})
	// create real producer
	if options.connectInSyncEnable {
		// sync create
		r.getProducer()
		r.ready = true
	} else {
		// async create
		go r.getProducer()
	}
	return r, nil
}

func (r *router) getProducer() pulsar.Producer {
	if r.Producer != nil {
		// Producer was already initialized
		return r.Producer
	}

	// Retry to create producer indefinitely
	backoffPolicy := &backoff.Backoff{}
	for {
		producer, err := r.client.CreateProducer(pulsar.ProducerOptions{
			Topic:                   r.options.Topic,
			CompressionType:         pulsar.LZ4,
			BatchingMaxPublishDelay: 100 * time.Millisecond,
		})

		if err != nil {
			r.logger.WithError(err).Errorf("Failed to create producer for topic: %s", r.options.Topic)
			time.Sleep(backoffPolicy.Next())
			continue
		} else {
			r.Producer = producer
			r.ready = true
			return producer
		}
	}
}
