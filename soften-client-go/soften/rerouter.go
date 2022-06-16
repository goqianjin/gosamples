package soften

import (
	"context"
	"errors"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
	"github.com/shenqianjin/soften-client-go/soften/internal/backoff"
)

// ------ base router ------

type routerOption struct {
	Topic         string
	Enable        bool
	MaxDeliveries uint
}

type router struct {
	client    pulsar.Client
	producer  pulsar.Producer
	option    routerOption
	messageCh chan *RerouteMessage
	closeCh   chan interface{}
	log       log.Logger
}

func newRouter(logger log.Logger, client pulsar.Client, option routerOption) (*router, error) {
	r := &router{
		client: client,
		option: option,
		log:    logger,
	}

	if option.Enable {
		if option.MaxDeliveries <= 0 {
			return nil, errors.New("routerOption.MaxDeliveries needs to be > 0")
		}

		if option.Topic == "" {
			return nil, errors.New("routerOption.Topic needs to be set to a valid topic name")
		}

		r.messageCh = make(chan *RerouteMessage)
		r.closeCh = make(chan interface{}, 1)
		r.log = logger.SubLogger(log.Fields{"rlq-topic": option.Topic})
		go r.run()
	}
	return r, nil
}

func (r *router) Chan() chan *RerouteMessage {
	return r.messageCh
}

func (r *router) run() {
	for {
		select {
		case rm := <-r.messageCh:
			r.log.WithField("msgID", rm.consumerMsg.ID()).Debugf("Got message for topic: %s", r.option.Topic)
			producer := r.getProducer()

			msgID := rm.consumerMsg.ID()
			producer.SendAsync(context.Background(), &rm.producerMsg, func(messageID pulsar.MessageID,
				producerMessage *pulsar.ProducerMessage, err error) {
				if err != nil {
					r.log.WithError(err).WithField("msgID", msgID).Errorf("Failed to send message to topic: %s", r.option.Topic)
					rm.consumerMsg.Consumer.Nack(rm.consumerMsg)
				} else {
					r.log.WithField("msgID", msgID).Debugf("Succeed to send message to topic: %s", r.option.Topic)
					rm.consumerMsg.Consumer.AckID(msgID)
				}
			})

		case <-r.closeCh:
			if r.producer != nil {
				r.producer.Close()
			}
			r.log.Debugf("Closed router for topic: %s", r.option.Topic)
			return
		}
	}
}

func (r *router) close() {
	// Attempt to write on the close channel, without blocking
	select {
	case r.closeCh <- nil:
	default:
	}
}

func (r *router) getProducer() pulsar.Producer {
	if r.producer != nil {
		// Producer was already initialized
		return r.producer
	}

	// Retry to create producer indefinitely
	backoff := &backoff.Backoff{}
	for {
		producer, err := r.client.CreateProducer(pulsar.ProducerOptions{
			Topic:                   r.option.Topic,
			CompressionType:         pulsar.LZ4,
			BatchingMaxPublishDelay: 100 * time.Millisecond,
		})

		if err != nil {
			r.log.WithError(err).Errorf("Failed to create producer for topic: %s", r.option.Topic)
			time.Sleep(backoff.Next())
			continue
		} else {
			r.producer = producer
			return producer
		}
	}
}
