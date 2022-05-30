package soam

import (
	"context"
	"errors"
	"soam/soam/internal"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
)

// ------ router interface ------

/*type Router interface {
	Route(message pulsar.ConsumerMessage) (routed bool)
}

type CustomRouter interface {
	RouteTo(message pulsar.ConsumerMessage, topic string) (routed bool)
}*/

// ------ re-router message ------

type ReRouterMessage struct {
	producerMsg pulsar.ProducerMessage
	consumerMsg pulsar.ConsumerMessage
}

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
	messageCh chan *ReRouterMessage
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

		r.messageCh = make(chan *ReRouterMessage)
		r.closeCh = make(chan interface{}, 1)
		r.log = logger.SubLogger(log.Fields{"rlq-topic": option.Topic})
		go r.run()
	}
	return r, nil
}

func (r *router) Chan() chan *ReRouterMessage {
	return r.messageCh
}

func (r *router) run() {
	for {
		select {
		case rm := <-r.messageCh:
			r.log.WithField("msgID", rm.consumerMsg.ID()).Debug("Got message for RLQ")
			producer := r.getProducer()

			msgID := rm.consumerMsg.ID()
			producer.SendAsync(context.Background(), &rm.producerMsg, func(messageID pulsar.MessageID,
				producerMessage *pulsar.ProducerMessage, err error) {
				if err != nil {
					r.log.WithError(err).WithField("msgID", msgID).Error("Failed to send message to RLQ")
					rm.consumerMsg.Consumer.Nack(rm.consumerMsg)
				} else {
					r.log.WithField("msgID", msgID).Debug("Succeed to send message to RLQ")
					rm.consumerMsg.Consumer.AckID(msgID)
				}
			})

		case <-r.closeCh:
			if r.producer != nil {
				r.producer.Close()
			}
			r.log.Debug("Closed RLQ router")
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
	backoff := &internal.Backoff{}
	for {
		producer, err := r.client.CreateProducer(pulsar.ProducerOptions{
			Topic:                   r.option.Topic,
			CompressionType:         pulsar.LZ4,
			BatchingMaxPublishDelay: 100 * time.Millisecond,
		})

		if err != nil {
			r.log.WithError(err).Error("Failed to create RLQ producer")
			time.Sleep(backoff.Next())
			continue
		} else {
			r.producer = producer
			return producer
		}
	}
}
