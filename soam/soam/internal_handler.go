package soam

import "github.com/apache/pulsar-client-go/pulsar"

type internalHandler interface {
	Handle(msg pulsar.ConsumerMessage) (success bool)
}

type internalRerouteHandler interface {
	Handle(msg pulsar.ConsumerMessage, topic string) (success bool)
}
