package soften

import (
	"sync"

	"github.com/shenqianjin/soften-client-go/soften/checker"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
	"github.com/shenqianjin/soften-client-go/soften/config"
	"github.com/shenqianjin/soften-client-go/soften/internal"
	"github.com/shenqianjin/soften-client-go/soften/message"
)

type rerouteHandler struct {
	routers     map[string]*reRouter
	routersLock sync.RWMutex
	client      pulsar.Client
	logger      log.Logger
	policy      *config.ReroutePolicy
	metrics     *internal.ListenerDecideGotoMetrics
}

func newRerouteHandler(client *client, listener *consumeListener, policy *config.ReroutePolicy) (*rerouteHandler, error) {
	routers := make(map[string]*reRouter)
	metrics := client.metricsProvider.GetListenerDecideGotoMetrics(listener.logTopics, listener.logLevels, internalGotoReroute)
	rtrHandler := &rerouteHandler{logger: client.logger, routers: routers, policy: policy}
	metrics.DecidersOpened.Inc()
	return rtrHandler, nil
}

func (d *rerouteHandler) Decide(msg pulsar.ConsumerMessage, cheStatus checker.CheckStatus) bool {
	if !cheStatus.IsPassed() || cheStatus.GetRerouteTopic() == "" {
		return false
	}
	rtr, err := d.internalSafeGetReRouterInAsync(cheStatus.GetRerouteTopic())
	if err != nil {
		return false
	}
	if !rtr.ready {
		if d.policy.ConnectInSyncEnable {
			<-rtr.readyCh
		} else {
			return false
		}
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
		props[message.XPropertyOriginPublishTime] = msg.PublishTime().UTC().Format(internal.RFC3339TimeInSecondPattern)
	}

	producerMsg := pulsar.ProducerMessage{
		Payload:     msg.Payload(),
		Key:         msg.Key(),
		OrderingKey: msg.OrderingKey(),
		Properties:  props,
		EventTime:   msg.EventTime(),
	}
	rtr.Chan() <- &RerouteMessage{
		consumerMsg: msg,
		producerMsg: producerMsg,
	}
	return true
}

func (d *rerouteHandler) internalSafeGetReRouterInAsync(topic string) (*reRouter, error) {
	d.routersLock.RLock()
	rtr, ok := d.routers[topic]
	d.routersLock.RUnlock()
	if ok {
		return rtr, nil
	}
	rtOption := reRouterOptions{Topic: topic, connectInSyncEnable: false}
	d.routersLock.Lock()
	defer d.routersLock.Unlock()
	rtr, ok = d.routers[topic]
	if ok {
		return rtr, nil
	}
	if newRtr, err := newReRouter(d.logger, d.client, rtOption); err != nil {
		return nil, err
	} else {
		rtr = newRtr
		d.routers[topic] = newRtr
		return rtr, nil
	}
}

func (d *rerouteHandler) close() {
	d.metrics.DecidersOpened.Dec()
}
