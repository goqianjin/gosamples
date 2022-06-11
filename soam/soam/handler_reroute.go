package soam

import (
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
)

type rerouteHandler struct {
	routers  map[string]*router
	consumer *consumer
	logger   log.Logger
}

func newRerouteHandler(logger log.Logger, consumer *consumer) (*rerouteHandler, error) {
	routers := make(map[string]*router)
	reRouter := &rerouteHandler{logger: logger, consumer: consumer, routers: routers}
	return reRouter, nil
}

func (hd *rerouteHandler) Handle(msg pulsar.ConsumerMessage, topic string) bool {
	if topic == "" {
		return false
	}
	if _, ok := hd.routers[topic]; !ok {
		rtOption := routerOption{Enable: true, Topic: topic}
		rt, err := newRouter(hd.logger, hd.consumer.client, rtOption)
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
	if _, ok := props[XPropertyOriginTopic]; !ok {
		props[XPropertyOriginTopic] = msg.Message.Topic()
	}
	if _, ok := props[XPropertyOriginMessageID]; !ok {
		props[XPropertyOriginMessageID] = MessageParser.GetMessageId(msg)
	}
	if _, ok := props[XPropertyOriginPublishTime]; !ok {
		props[XPropertyOriginPublishTime] = msg.PublishTime().Format(RFC3339TimeInSecondPattern)
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
