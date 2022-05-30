package soam

import "github.com/apache/pulsar-client-go/pulsar"

var MessageParser = &messageParser{}

type messageParser struct {
}

func (p *messageParser) GetCurrentStatus(message pulsar.ConsumerMessage) messageStatus {
	return messageStatus("") // TODO
}
