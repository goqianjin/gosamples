package soam

import (
	"errors"
	"fmt"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
)

type finalStatusHandler struct {
	logger log.Logger
	status messageStatus
}

func newFinalStatusHandler(logger log.Logger, status messageStatus) (*finalStatusHandler, error) {
	if status == "" {
		return nil, errors.New("final message status cannot be empty")
	}
	if status != MessageStatusDone && status != MessageStatusDiscard {
		return nil, errors.New(fmt.Sprintf("%s is not a final message status", status))
	}
	return &finalStatusHandler{logger: logger, status: status}, nil
}

func (h *finalStatusHandler) Handle(message pulsar.ConsumerMessage) (success bool) {
	switch h.status {
	case MessageStatusDone:
		message.Consumer.Ack(message.Message)
		h.logger.Warnf("Handle message: done")
		return true
	case MessageStatusDiscard:
		message.Consumer.Ack(message.Message)
		h.logger.Warnf("Handle message: discard")
		return true
	}
	return false
}
