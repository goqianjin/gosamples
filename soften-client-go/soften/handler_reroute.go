package soften

import (
	"sync"

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
}

func newRerouteHandler(client *client, policy *config.ReroutePolicy) (*rerouteHandler, error) {
	routers := make(map[string]*reRouter)
	rtrHandler := &rerouteHandler{logger: client.logger, routers: routers, policy: policy}
	return rtrHandler, nil
}

func (hd *rerouteHandler) Handle(msg pulsar.ConsumerMessage, topic string) bool {
	if topic == "" {
		return false
	}
	rtr, err := hd.internalSafeGetReRouterInAsync(topic)
	if err != nil {
		return false
	}
	if !rtr.ready {
		if hd.policy.ConnectInSyncEnable {
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
		props[message.XPropertyOriginPublishTime] = msg.PublishTime().Format(internal.RFC3339TimeInSecondPattern)
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

func (hd *rerouteHandler) internalSafeGetReRouterInAsync(topic string) (*reRouter, error) {
	hd.routersLock.RLock()
	rtr, ok := hd.routers[topic]
	hd.routersLock.RUnlock()
	if ok {
		return rtr, nil
	}
	rtOption := reRouterOptions{Topic: topic, connectInSyncEnable: false}
	hd.routersLock.Lock()
	defer hd.routersLock.Unlock()
	rtr, ok = hd.routers[topic]
	if ok {
		return rtr, nil
	}
	if newRtr, err := newReRouter(hd.logger, hd.client, rtOption); err != nil {
		return nil, err
	} else {
		rtr = newRtr
		hd.routers[topic] = newRtr
		return rtr, nil
	}
}
