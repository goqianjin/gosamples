package soam

import (
	"context"
	"soam/soam/internal"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
)

type RetryMessage struct {
	producerMsg pulsar.ProducerMessage
	consumerMsg pulsar.ConsumerMessage
}

type retryRouter struct {
	client    pulsar.Client
	producer  pulsar.Producer
	policy    *pulsar.DLQPolicy
	messageCh chan RetryMessage
	closeCh   chan interface{}
	log       log.Logger
}

func newRetryRouter(client pulsar.Client, policy *pulsar.DLQPolicy, retryEnabled bool, logger log.Logger) (*retryRouter, error) {
	r := &retryRouter{
		client: client,
		policy: policy,
		log:    logger,
	}

	if policy != nil && retryEnabled {
		if policy.MaxDeliveries <= 0 {
			return nil, nil //newError(InvalidConfiguration, "DLQPolicy.MaxDeliveries needs to be > 0")
		}

		if policy.RetryLetterTopic == "" {
			return nil, nil // newError(InvalidConfiguration, "DLQPolicy.RetryLetterTopic needs to be set to a valid topic name")
		}

		r.messageCh = make(chan RetryMessage)
		r.closeCh = make(chan interface{}, 1)
		r.log = logger.SubLogger(log.Fields{"rlq-topic": policy.RetryLetterTopic})
		go r.run()
	}
	return r, nil
}

func (r *retryRouter) Chan() chan RetryMessage {
	return r.messageCh
}

func (r *retryRouter) run() {
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

func (r *retryRouter) close() {
	// Attempt to write on the close channel, without blocking
	select {
	case r.closeCh <- nil:
	default:
	}
}

func (r *retryRouter) getProducer() pulsar.Producer {
	if r.producer != nil {
		// Producer was already initialized
		return r.producer
	}

	// Retry to create producer indefinitely
	backoff := &internal.Backoff{}
	for {
		producer, err := r.client.CreateProducer(pulsar.ProducerOptions{
			Topic:                   r.policy.RetryLetterTopic,
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
