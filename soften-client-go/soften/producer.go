package soften

import (
	"context"
	"sync"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
	"github.com/shenqianjin/soften-client-go/soften/config"
	"github.com/shenqianjin/soften-client-go/soften/internal"
)

type producer struct {
	pulsar.Producer
	client      *client
	logger      log.Logger
	routeEnable bool
	routePolicy *config.RoutePolicy
	checkers    []internal.RouteChecker
	routers     map[string]*router
	routersLock sync.RWMutex
	metrics     *internal.ProducerMetrics
}

func newProducer(client *client, conf *config.ProducerConfig, checkers ...internal.RouteChecker) (*producer, error) {
	options := pulsar.ProducerOptions{
		Topic: conf.Topic,
	}
	pulsarProducer, err := client.Client.CreateProducer(options)
	if err != nil {
		return nil, err
	}
	p := &producer{
		Producer:    pulsarProducer,
		client:      client,
		logger:      client.logger.SubLogger(log.Fields{"topic": options.Topic}),
		routeEnable: conf.RouteEnable,
		routePolicy: conf.Route,
		checkers:    checkers,
		routers:     map[string]*router{},
		metrics:     client.metricsProvider.GetProducerMetrics(conf.Topic),
	}
	p.logger.Infof("created soften producer")
	p.metrics.ProducersOpened.Inc()
	return p, nil
}

// Send aim to send message synchronously
func (p *producer) Send(ctx context.Context, msg *pulsar.ProducerMessage) (pulsar.MessageID, error) {
	if !p.routeEnable {
		start := time.Now()
		msgId, err := p.Producer.Send(ctx, msg)
		p.metrics.PublishLatency.Observe(time.Now().Sub(start).Seconds())
		return msgId, err
	}
	for _, chk := range p.checkers {
		routeTopic := chk(msg)
		if routeTopic == "" {
			continue
		}
		rtr, err := p.internalSafeGetRouterInAsync(routeTopic)
		if err != nil {
			p.logger.Warnf("failed to create router for topic: %s", routeTopic)
			continue
		}
		if !rtr.ready {
			// wait router until it's ready
			if p.routePolicy.ConnectInSyncEnable {
				<-rtr.readyCh
			} else {
				// back to other router or main topic before the checked router is ready
				p.logger.Warnf("router is still not ready for topic: %s", routeTopic)
				continue
			}
		}
		if mid, err2 := rtr.Send(ctx, msg); err2 != nil {
			p.logger.Warnf("failed to send message to topic: %s", routeTopic)
			if !p.routePolicy.BackEnable {
				return mid, err2
			}
		} else {
			return mid, err2
		}
	}

	start := time.Now()
	msgId, err := p.Producer.Send(ctx, msg)
	p.metrics.PublishLatency.Observe(time.Now().Sub(start).Seconds())
	return msgId, err
}

// SendAsync send message asynchronously
func (p *producer) SendAsync(ctx context.Context, msg *pulsar.ProducerMessage,
	callback func(pulsar.MessageID, *pulsar.ProducerMessage, error)) {
	start := time.Now()
	callbackNew := func(msgID pulsar.MessageID, msg *pulsar.ProducerMessage, err error) {
		p.metrics.PublishLatency.Observe(time.Now().Sub(start).Seconds())
		callback(msgID, msg, err)
	}
	if !p.routeEnable {
		p.Producer.SendAsync(ctx, msg, callbackNew)
		return
	}
	for _, chk := range p.checkers {
		routeTopic := chk(msg)
		if routeTopic == "" {
			continue
		}
		rtr, err := p.internalSafeGetRouterInAsync(routeTopic)
		if err != nil {
			p.logger.Warnf("failed to create router for topic: %s", routeTopic)
			continue
		}
		if !rtr.ready {
			// wait router until it's ready
			if p.routePolicy.ConnectInSyncEnable {
				<-rtr.readyCh
			} else {
				// back to other router or main topic before the checked router is ready
				p.logger.Warnf("router is still not ready for topic: %s", routeTopic)
				continue
			}
		}
		// route record metrics individually
		rtr.SendAsync(ctx, msg, callback)
		//
		return
	}
	p.Producer.SendAsync(ctx, msg, callbackNew)
	return
}

func (p *producer) internalSafeGetRouterInAsync(topic string) (*router, error) {
	p.routersLock.RLock()
	rtr, ok := p.routers[topic]
	p.routersLock.RUnlock()
	if ok {
		return rtr, nil
	}
	options := routerOptions{Topic: topic, connectInSyncEnable: false}
	p.routersLock.Lock()
	defer p.routersLock.Unlock()
	rtr, ok = p.routers[topic]
	if ok {
		return rtr, nil
	}
	if newRtr, err := newRouter(p.logger, p.client.Client, options); err != nil {
		return nil, err
	} else {
		rtr = newRtr
		p.routers[topic] = newRtr
		return rtr, nil
	}
}

func (p *producer) Close() {
	p.Producer.Close()
	p.logger.Info("closed soften producer")
	p.metrics.ProducersOpened.Dec()
}
