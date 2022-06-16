package soften

import (
	"errors"
	"github.com/shenqianjin/soften-client-go/soften/topic"

	"github.com/shenqianjin/soften-client-go/soften/message"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
	"github.com/shenqianjin/soften-client-go/soften/internal"
)

type gradeHandler struct {
	router *router
	logger log.Logger
	level  internal.TopicLevel
}

func newGradeHandler(client *client, tpc string, level internal.TopicLevel) (*gradeHandler, error) {
	if tpc == "" {
		return nil, errors.New("topic cannot be blank")
	}
	if level == "" {
		return nil, errors.New("topic level is empty")
	}
	suffix, err := topic.NameSuffixOf(level)
	if err != nil {
		return nil, err
	}
	routerOption := routerOption{Enable: true, Topic: tpc + suffix}
	rt, err := newRouter(client.logger, client, routerOption)
	if err != nil {
		return nil, err
	}
	hd := &gradeHandler{router: rt, logger: client.logger, level: level}
	return hd, nil
}

func (hd *gradeHandler) Handle(msg pulsar.ConsumerMessage) bool {
	// prepare to upgrade / degrade
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
	hd.router.Chan() <- &RerouteMessage{
		consumerMsg: msg,
		producerMsg: producerMsg,
	}
	return true
}
