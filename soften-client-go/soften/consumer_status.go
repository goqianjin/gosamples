package soften

import (
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
	"github.com/shenqianjin/soften-client-go/soften/checker"
	"github.com/shenqianjin/soften-client-go/soften/config"
	"github.com/shenqianjin/soften-client-go/soften/internal"
	"github.com/shenqianjin/soften-client-go/soften/message"
)

type statusConsumer struct {
	pulsar.Consumer
	logger          log.Logger
	status          internal.MessageStatus
	policy          *config.StatusPolicy
	statusMessageCh chan ConsumerMessage // channel used to deliver message to clients
	handler         internalDecider
}

func newStatusConsumer(parentLogger log.Logger, pulsarConsumer pulsar.Consumer, status internal.MessageStatus, policy *config.StatusPolicy, handler internalDecider) *statusConsumer {
	sc := &statusConsumer{
		logger:          parentLogger.SubLogger(log.Fields{"status": status}),
		Consumer:        pulsarConsumer,
		status:          status,
		policy:          policy,
		handler:         handler,
		statusMessageCh: make(chan ConsumerMessage, 10),
	}
	go sc.start()
	sc.logger.Info("created status consumer")
	return sc
}

func (sc *statusConsumer) start() {
	if sc.status != message.StatusReady {
		sc.logger.Info("------start---- %v", sc.status)
	}

	for {
		// block to read pulsar chan
		msg := <-sc.Consumer.Chan()

		now := time.Now().UTC()

		// wait to reentrant time
		reentrantTime := message.Parser.GetReentrantTime(msg)
		if !reentrantTime.IsZero() {
			if reentrantTime.After(now) {
				time.Sleep(reentrantTime.Sub(now))
			}
		}

		reconsumeTime := message.Parser.GetReconsumeTime(msg)
		// delivery message to client if reconsumeTime doesn't exist or before/equal now
		// 'not exist' means it is consumed first time, 'before/equal now' means it is time to consume
		if reconsumeTime.IsZero() || !reconsumeTime.After(now) {
			sc.statusMessageCh <- ConsumerMessage{
				ConsumerMessage: msg,
				StatusMessage:   &statusMessage{status: sc.status},
			}
			continue
		}

		// Nack or wait by hang it in memory util meet the reconsume time if reconsumeTime is before `now + policy.MaxNackDelay`
		maxNackDelayAt := time.Now().UTC().Add(time.Duration(sc.policy.ReentrantDelay))
		if reconsumeTime.Before(maxNackDelayAt) {
			// do not increase reconsume times if it is in order to meet delays
			msg.Consumer.NackLater(msg, maxNackDelayAt.Sub(reconsumeTime))
			continue
		}

		// reentrant again util meet the reconsume time
		if ok := sc.handler.Decide(msg, checker.CheckStatusPassed); ok {
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

func (sc *statusConsumer) Close() {
	sc.Consumer.Close()
	sc.logger.Info("created status consumer")
}
