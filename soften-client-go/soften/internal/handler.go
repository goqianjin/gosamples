package internal

import "github.com/apache/pulsar-client-go/pulsar"

type Handler interface {
	Handle(msg pulsar.ConsumerMessage) (success bool)
}

type RerouteHandler interface {
	Handle(msg pulsar.ConsumerMessage, topic string) (success bool)
}
