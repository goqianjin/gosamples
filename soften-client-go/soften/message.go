package soften

import (
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/shenqianjin/soften-client-go/soften/internal"
)

// ------ re-reRouter message ------

type RerouteMessage struct {
	producerMsg pulsar.ProducerMessage
	consumerMsg pulsar.ConsumerMessage
}

// ------ custom consumer message ------

type ConsumerMessage struct {
	pulsar.ConsumerMessage
	StatusMessage
	LeveledMessage
}

// ------ status message interface ------

type StatusMessage interface {
	Status() internal.MessageStatus
}

// ------ leveled message interface ------

type LeveledMessage interface {
	Level() internal.TopicLevel
}

// ---------------------------------------

// ------ status message implementation ------

type statusMessage struct {
	status internal.MessageStatus
}

func (m *statusMessage) Status() internal.MessageStatus {
	return m.status
}

// ------ leveled message implementation ------

type leveledMessage struct {
	level internal.TopicLevel
}

func (m *leveledMessage) Level() internal.TopicLevel {
	return m.level
}
