package soften

import (
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
	"github.com/shenqianjin/soften-client-go/soften/internal"
	"github.com/shenqianjin/soften-client-go/soften/message"
)

type rerouteHandler struct {
	routers map[string]*reRouter
	client  pulsar.Client
	logger  log.Logger
}

func newRerouteHandler(client *client) (*rerouteHandler, error) {
	routers := make(map[string]*reRouter)
	rtrHandler := &rerouteHandler{logger: client.logger, routers: routers}
	return rtrHandler, nil
}

func (hd *rerouteHandler) Handle(msg pulsar.ConsumerMessage, topic string) bool {
	if topic == "" {
		return false
	}
	if _, ok := hd.routers[topic]; !ok {
		rtOption := reRouterOptions{Enable: true, Topic: topic}
		rt, err := newReRouter(hd.logger, hd.client, rtOption)
		if err != nil {
			return false
		}
		hd.routers[topic] = rt
	}
	// prepare to reroute
	props := make(map[string]string)
	for k, v := range msg.Properties() {
		props[k] = v
	}
	// record origin information when re-route first time
	if _, ok := props[message.XPropertyOriginTopic]; !ok {
		props[message.XPropertyOriginTopic] = msg.Message.Topic()
	}
	if _, ok := props[message.XPropertyOriginMessageID]; !ok {
		props[message.XPropertyOriginMessageID] = message.Parser.GetMessageId(msg)
	}
	if _, ok := props[message.XPropertyOriginPublishTime]; !ok {
		props[message.XPropertyOriginPublishTime] = msg.PublishTime().Format(internal.RFC3339TimeInSecondPattern)
	}

	producerMsg := pulsar.ProducerMessage{
		Payload:     msg.Payload(),
		Key:         msg.Key(),
		OrderingKey: msg.OrderingKey(),
		Properties:  props,
		EventTime:   msg.EventTime(),
	}
	hd.routers[topic].Chan() <- &RerouteMessage{
		consumerMsg: msg,
		producerMsg: producerMsg,
	}
	return true
}
