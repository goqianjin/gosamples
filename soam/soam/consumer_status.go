package soam

import (
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

type statusConsumer struct {
	pulsar.Consumer
	outerConsumer   *consumer
	status          messageStatus
	policy          *StatusPolicy
	statusMessageCh chan pulsar.ConsumerMessage // channel used to deliver message to clients
}

func newStatusConsumer(outerConsumer *consumer, pulsarConsumer pulsar.Consumer, status messageStatus, policy *StatusPolicy) *statusConsumer {
	sc := &statusConsumer{
		Consumer:      pulsarConsumer,
		outerConsumer: outerConsumer,
		status:        status,
		policy:        policy,
	}
	sc.statusMessageCh = make(chan pulsar.ConsumerMessage, 5)
	go sc.start()
	return sc
}

func (sc *statusConsumer) start() {
	for {
		// block to read consumer chan
		msg := <-sc.Consumer.Chan()

		// wait to reentrant time
		reentrantTime := MessageParser.GetReentrantTime(msg)
		if !reentrantTime.IsZero() {
			time.Sleep(time.Now().Sub(reentrantTime))
		}

		reconsumeTime := MessageParser.GetReconsumeTime(msg)
		// delivery message to client if reconsumeTime doesn't exist or before/equal now
		if reconsumeTime.IsZero() || !reconsumeTime.After(time.Now()) {
			sc.statusMessageCh <- msg
			continue
		}

		// Nack or wait by hang it in memory util meet the reconsume time if reconsumeTime is before `now + policy.MaxNackDelay`
		maxNackDelayAt := time.Now().Add(time.Duration(sc.policy.ReentrantDelay))
		if reconsumeTime.Before(maxNackDelayAt) {
			// do not increase reconsume times if it is in order to meet delays
			msg.Consumer.NackLater(msg, maxNackDelayAt.Sub(reconsumeTime))
			continue
		}

		// reentrant again or Nack util meet the reconsume time
		if ok := sc.outerConsumer.internalRouteByStatus(sc.status, msg); ok {
			continue
		}

		// let application handle other case
		sc.statusMessageCh <- msg
	}
}

func (sc *statusConsumer) StatusChan() <-chan pulsar.ConsumerMessage {
	return sc.statusMessageCh
}
