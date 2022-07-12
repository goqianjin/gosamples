package soften

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
	"github.com/shenqianjin/soften-client-go/soften/checker"
	"github.com/shenqianjin/soften-client-go/soften/config"
	"github.com/shenqianjin/soften-client-go/soften/internal"
	"github.com/shenqianjin/soften-client-go/soften/message"
)

type routeDeciderOptions struct {
	connectInSyncEnable bool // Optional: 是否同步建立连接, 首次发送消息需阻塞等待客户端与服务端连接完成

	// extra for upgrade/degrade

	upgradeLevel internal.TopicLevel
	degradeLevel internal.TopicLevel
}

type routeDecider struct {
	logger      log.Logger
	client      pulsar.Client
	topic       string
	gotoAction  internal.MessageGoto
	options     *routeDeciderOptions
	routers     map[string]*router
	routersLock sync.RWMutex
	metrics     *internal.DecideGotoMetrics
	policy      *config.ReroutePolicy
}

func newRouteDecider(producer *producer, gotoAction internal.MessageGoto, options *routeDeciderOptions) (*routeDecider, error) {
	if gotoAction == message.GotoUpgrade {
		if options.upgradeLevel == "" {
			return nil, errors.New(fmt.Sprintf("upgrade level is missing for %v router", gotoAction))
		}
	}
	if gotoAction == message.GotoDegrade {
		if options.degradeLevel == "" {
			return nil, errors.New(fmt.Sprintf("upgrade level is missing for %v router", gotoAction))
		}
	}
	metrics := producer.client.metricsProvider.GetProducerDecideGotoMetrics(producer.topic, gotoAction.String())
	rtrDecider := &routeDecider{
		logger:     producer.logger,
		client:     producer.client.Client,
		topic:      producer.topic,
		gotoAction: gotoAction,
		options:    options,
		routers:    make(map[string]*router),
		metrics:    metrics,
	}
	metrics.DecidersOpened.Inc()
	return rtrDecider, nil
}

func (d *routeDecider) Decide(ctx context.Context, msg *pulsar.ProducerMessage,
	checkStatus checker.CheckStatus) (mid pulsar.MessageID, err error, decided bool) {
	if !checkStatus.IsPassed() {
		return nil, nil, false
	}
	// process discard
	if d.gotoAction == message.GotoDiscard {
		d.logger.Warnf(fmt.Sprintf("discard message. payload size: %v, properties: %v", len(msg.Payload), msg.Properties))
		decided = true
		return
	}
	// parse topic
	routeTopic := d.parseRouteTopic(checkStatus)
	// get or create router
	rtr, err := d.internalSafeGetRouterInAsync(routeTopic)
	if err != nil {
		d.logger.Warnf("failed to create router for topic: %s", routeTopic)
		return nil, err, false
	}
	if !rtr.ready {
		// wait router until it's ready
		if d.options.connectInSyncEnable {
			<-rtr.readyCh
		} else {
			// back to other router or main topic before the checked router is ready
			d.logger.Warnf("router is still not ready for topic: %s", routeTopic)
			return nil, nil, false
		}
	}
	// send
	mid, err = rtr.Send(ctx, msg)
	if err != nil {
		d.logger.Warnf("failed to route message, payload size: %v, properties: %v", len(msg.Payload), msg.Properties)
		return mid, err, false
	}
	return mid, err, true
}

func (d *routeDecider) DecideAsync(ctx context.Context, msg *pulsar.ProducerMessage, checkStatus checker.CheckStatus,
	callback func(pulsar.MessageID, *pulsar.ProducerMessage, error)) (decided bool) {
	if !checkStatus.IsPassed() {
		return false
	}
	// process discard
	if d.gotoAction == message.GotoDiscard {
		d.logger.Warnf(fmt.Sprintf("discard message. payload size: %v, properties: %v", len(msg.Payload), msg.Properties))
		decided = true
		return
	}
	// parse topic
	routeTopic := d.parseRouteTopic(checkStatus)
	// get or create router
	rtr, err := d.internalSafeGetRouterInAsync(routeTopic)
	if err != nil {
		d.logger.Warnf("failed to create router for topic: %s", routeTopic)
		return false
	}
	if !rtr.ready {
		// wait router until it's ready
		if d.options.connectInSyncEnable {
			<-rtr.readyCh
		} else {
			// back to other router or main topic before the checked router is ready
			d.logger.Warnf("router is still not ready for topic: %s", routeTopic)
			return false
		}
	}
	// send
	rtr.SendAsync(ctx, msg, callback)

	return true

}

func (d *routeDecider) parseRouteTopic(checkStatus checker.CheckStatus) string {
	routeTopic := ""
	switch d.gotoAction {
	case message.GotoDead:
		routeTopic = d.topic + message.StatusDead.TopicSuffix()
	case message.GotoUpgrade:
		routeTopic = d.topic + d.options.upgradeLevel.TopicSuffix()
	case message.GotoDegrade:
		routeTopic = d.topic + d.options.degradeLevel.TopicSuffix()
	case message.GotoBlocking:
		routeTopic = d.topic + message.StatusBlocking.TopicSuffix()
	case message.GotoPending:
		routeTopic = d.topic + message.StatusPending.TopicSuffix()
	case message.GotoRetrying:
		routeTopic = d.topic + message.StatusRetrying.TopicSuffix()
	case internalGotoRoute:
		routeTopic = checkStatus.GetRerouteTopic()
	}
	return routeTopic
}

func (d *routeDecider) internalSafeGetRouterInAsync(topic string) (*router, error) {
	d.routersLock.RLock()
	rtr, ok := d.routers[topic]
	d.routersLock.RUnlock()
	if ok {
		return rtr, nil
	}
	rtOption := routerOptions{Topic: topic, connectInSyncEnable: false}
	d.routersLock.Lock()
	defer d.routersLock.Unlock()
	rtr, ok = d.routers[topic]
	if ok {
		return rtr, nil
	}
	if newRtr, err := newRouter(d.logger, d.client, rtOption); err != nil {
		return nil, err
	} else {
		rtr = newRtr
		d.routers[topic] = newRtr
		return rtr, nil
	}
}

func (d *routeDecider) close() {
	d.metrics.DecidersOpened.Dec()
}
