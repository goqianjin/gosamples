package soam

import (
	"errors"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
)

type gradeHandler struct {
	router   *router
	consumer *consumer
	logger   log.Logger
	level    TopicLevel
}

func newGradeHandler(logger log.Logger, consumer *consumer, level TopicLevel) (*gradeHandler, error) {
	if level == "" {
		return nil, errors.New("topic level is empty")
	}
	routerOption := routerOption{Enable: true, Topic: consumer.config.Topic + "_" + string(level)}
	rt, err := newRouter(logger, consumer.client, routerOption)
	if err != nil {
		return nil, err
	}
	hd := &gradeHandler{router: rt, consumer: consumer, logger: logger, level: level}
	return hd, nil
}

func (hd *gradeHandler) Handle(msg pulsar.ConsumerMessage) bool {
	// prepare to upgrade / degrade
	props := make(map[string]string)
	for k, v := range msg.Properties() {
		props[k] = v
	}
	// record origin information when re-route first time
	if _, ok := props[SysPropertyOriginTopic]; !ok {
		props[SysPropertyOriginTopic] = msg.Message.Topic()
	}
	if _, ok := props[SysPropertyOriginMessageID]; !ok {
		props[SysPropertyOriginMessageID] = MessageParser.GetMessageId(msg)
	}
	if _, ok := props[SysPropertyOriginPublishTime]; !ok {
		props[SysPropertyOriginPublishTime] = msg.PublishTime().Format(RFC3339TimeInSecondPattern)
	}

	producerMsg := pulsar.ProducerMessage{
		Payload:     msg.Payload(),
		Key:         msg.Key(),
		OrderingKey: msg.OrderingKey(),
		Properties:  props,
		EventTime:   msg.EventTime(),
	}
	hd.router.Chan() <- &ReRouterMessage{
		consumerMsg: msg,
		producerMsg: producerMsg,
	}
	return true
}
