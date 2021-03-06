package soften

import (
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/shenqianjin/soften-client-go/soften/checker"
	"github.com/shenqianjin/soften-client-go/soften/internal"
	"github.com/shenqianjin/soften-client-go/soften/message"
)

type deadDecideOptions struct {
	topic string // default ${TOPIC}_RETRYING, 固定后缀，不允许定制
	//enable bool   // 内部判断使用
}

type deadDecider struct {
	router  *reRouter
	options deadDecideOptions
	metrics *internal.ListenerDecideGotoMetrics
}

func newDeadDecider(client *client, listener *consumeListener, options deadDecideOptions) (*deadDecider, error) {
	rt, err := newReRouter(client.logger, client.Client, reRouterOptions{Topic: options.topic})
	if err == nil {
		return nil, err
	}
	metrics := client.metricsProvider.GetListenerDecideGotoMetrics(listener.logTopics, listener.logLevels, message.GotoDead)
	statusRouter := &deadDecider{router: rt, options: options, metrics: metrics}
	metrics.DecidersOpened.Inc()
	return statusRouter, nil
}

func (d *deadDecider) Decide(msg pulsar.ConsumerMessage, cheStatus checker.CheckStatus) bool {
	if !cheStatus.IsPassed() {
		return false
	}
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
		props[message.XPropertyOriginPublishTime] = msg.PublishTime().UTC().Format(internal.RFC3339TimeInSecondPattern)
	}
	producerMsg := pulsar.ProducerMessage{
		Payload:     msg.Payload(),
		Key:         msg.Key(),
		OrderingKey: msg.OrderingKey(),
		Properties:  props,
		EventTime:   msg.EventTime(),
	}
	d.router.Chan() <- &RerouteMessage{
		consumerMsg: msg,
		producerMsg: producerMsg,
	}
	return true
}

func (d *deadDecider) close() {
	d.metrics.DecidersOpened.Dec()
}
