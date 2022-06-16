package soften

import (
	"errors"
	"fmt"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
	"github.com/shenqianjin/soften-client-go/soften/internal"
	"github.com/shenqianjin/soften-client-go/soften/message"
)

type finalStatusHandler struct {
	logger log.Logger
	status internal.MessageStatus
}

func newFinalStatusHandler(logger log.Logger, status internal.MessageStatus) (*finalStatusHandler, error) {
	if status == "" {
		return nil, errors.New("final message status cannot be empty")
	}
	if status != message.StatusDone && status != message.StatusDiscard {
		return nil, errors.New(fmt.Sprintf("%s is not a final message status", status))
	}
	return &finalStatusHandler{logger: logger, status: status}, nil
}

func (h *finalStatusHandler) Handle(msg pulsar.ConsumerMessage) (success bool) {
	switch h.status {
	case message.StatusDone:
		msg.Consumer.Ack(msg.Message)
		h.logger.Warnf("Handle message: done")
		return true
	case message.StatusDiscard:
		msg.Consumer.Ack(msg.Message)
		h.logger.Warnf("Handle message: discard")
		return true
	}
	return false
}
