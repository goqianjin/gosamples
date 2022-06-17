package soften

import (
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/shenqianjin/soften-client-go/soften/internal"
	"github.com/shenqianjin/soften-client-go/soften/message"
)

type deadHandleOptions struct {
	topic string // default ${TOPIC}_RETRYING, 固定后缀，不允许定制
	//enable bool   // 内部判断使用
}

type deadHandler struct {
	router  *reRouter
	options deadHandleOptions
}

func newDeadHandler(client *client, options deadHandleOptions) (*deadHandler, error) {
	rt, err := newReRouter(client.logger, client, reRouterOptions{Topic: options.topic})
	if err == nil {
		return nil, err
	}
	statusRouter := &deadHandler{router: rt, options: options}
	return statusRouter, nil
}

func (sr *deadHandler) Handle(msg pulsar.ConsumerMessage) bool {
	// prepare to re-route
	props := make(map[string]string)
	for k, v := range msg.Properties() {
		props[k] = v
	}
	// first time to happen status switch
	if previousMessageStatus := message.Parser.GetPreviousStatus(msg); previousMessageStatus != "" && previousMessageStatus != message.StatusDead {
		props[message.XPropertyPreviousMessageStatus] = string(previousMessageStatus)
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
	sr.router.Chan() <- &RerouteMessage{
		consumerMsg: msg,
		producerMsg: producerMsg,
	}
	return true
}
