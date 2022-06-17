package soften

import (
	"context"
	"errors"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
	"github.com/shenqianjin/soften-client-go/soften/internal/backoff"
)

// ------ base reRouter ------

type reRouterOptions struct {
	Topic               string
	connectInSyncEnable bool
	connectMaxRetries   uint
	//MaxDeliveries uint
}

type reRouter struct {
	client    pulsar.Client
	producer  pulsar.Producer
	option    reRouterOptions
	messageCh chan *RerouteMessage
	closeCh   chan interface{}
	logger    log.Logger
	ready     bool
	readyCh   chan struct{}
}

func newReRouter(logger log.Logger, client pulsar.Client, options reRouterOptions) (*reRouter, error) {
	r := &reRouter{
		client: client,
		option: options,
		logger: logger.SubLogger(log.Fields{"reroute-topic": options.Topic}),
	}

	if options.connectMaxRetries <= 0 {
		return nil, errors.New("reRouterOptions.connectMaxRetries needs to be > 0")
	}

	if options.Topic == "" {
		return nil, errors.New("reRouterOptions.Topic needs to be set to a valid topic name")
	}

	r.messageCh = make(chan *RerouteMessage)
	r.closeCh = make(chan interface{}, 1)
	if options.connectInSyncEnable {
		r.getProducer()
	}
	go r.run()
	return r, nil
}

func (r *reRouter) Chan() chan *RerouteMessage {
	return r.messageCh
}

func (r *reRouter) run() {
	for {
		select {
		case rm := <-r.messageCh:
			r.logger.WithField("msgID", rm.consumerMsg.ID()).Debugf("Got message for topic: %s", r.option.Topic)
			producer := r.getProducer()

			msgID := rm.consumerMsg.ID()
			producer.SendAsync(context.Background(), &rm.producerMsg, func(messageID pulsar.MessageID,
				producerMessage *pulsar.ProducerMessage, err error) {
				if err != nil {
					r.logger.WithError(err).WithField("msgID", msgID).Errorf("Failed to send message to topic: %s", r.option.Topic)
					rm.consumerMsg.Consumer.Nack(rm.consumerMsg)
				} else {
					r.logger.WithField("msgID", msgID).Debugf("Succeed to send message to topic: %s", r.option.Topic)
					rm.consumerMsg.Consumer.AckID(msgID)
				}
			})

		case <-r.closeCh:
			if r.producer != nil {
				r.producer.Close()
			}
			r.logger.Debugf("Closed reRouter for topic: %s", r.option.Topic)
			return
		}
	}
}

func (r *reRouter) close() {
	// Attempt to write on the close channel, without blocking
	select {
	case r.closeCh <- nil:
	default:
	}
}

func (r *reRouter) getProducer() pulsar.Producer {
	if r.producer != nil {
		// Producer was already initialized
		return r.producer
	}

	// Retry to create producer indefinitely
	backoffPolicy := &backoff.Backoff{}
	for {
		producer, err := r.client.CreateProducer(pulsar.ProducerOptions{
			Topic:                   r.option.Topic,
			CompressionType:         pulsar.LZ4,
			BatchingMaxPublishDelay: 100 * time.Millisecond,
		})

		if err != nil {
			r.logger.WithError(err).Errorf("Failed to create producer for topic: %s", r.option.Topic)
			time.Sleep(backoffPolicy.Next())
			continue
		} else {
			r.producer = producer
			r.ready = true
			r.readyCh <- struct{}{}
			close(r.readyCh)
			return producer
		}
	}
}
