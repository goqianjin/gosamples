package soften

import (
	"time"

	"github.com/shenqianjin/soften-client-go/soften/message"

	"github.com/shenqianjin/soften-client-go/soften/config"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
	"github.com/shenqianjin/soften-client-go/soften/internal"
)

type statusConsumer struct {
	pulsar.Consumer
	logger          log.Logger
	status          internal.MessageStatus
	policy          *config.StatusPolicy
	statusMessageCh chan ConsumerMessage // channel used to deliver message to clients
	handler         internal.Handler
}

func newStatusConsumer(pulsarConsumer pulsar.Consumer, status internal.MessageStatus, policy *config.StatusPolicy, handler internal.Handler) *statusConsumer {
	sc := &statusConsumer{
		Consumer:        pulsarConsumer,
		status:          status,
		policy:          policy,
		handler:         handler,
		statusMessageCh: make(chan ConsumerMessage, 5),
	}
	go sc.start()
	return sc
}

func (sc *statusConsumer) start() {
	for {
		// block to read multiStatusConsumeFacade chan
		msg := <-sc.Consumer.Chan()

		// wait to reentrant time
		reentrantTime := message.Parser.GetReentrantTime(msg)
		if !reentrantTime.IsZero() {
			time.Sleep(time.Now().Sub(reentrantTime))
		}

		reconsumeTime := message.Parser.GetReconsumeTime(msg)
		// delivery message to client if reconsumeTime doesn't exist or before/equal now
		if reconsumeTime.IsZero() || !reconsumeTime.After(time.Now()) {
			sc.statusMessageCh <- ConsumerMessage{
				ConsumerMessage: msg,
				StatusMessage:   &statusMessage{status: sc.status},
			}
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
		if ok := sc.handler.Handle(msg); ok {
			continue
		}

		// let application handle other case
		sc.statusMessageCh <- ConsumerMessage{
			ConsumerMessage: msg,
			StatusMessage:   &statusMessage{status: sc.status},
		}
	}
}

func (sc *statusConsumer) StatusChan() <-chan ConsumerMessage {
	return sc.statusMessageCh
}
