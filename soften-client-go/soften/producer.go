package soften

import (
	"context"
	"sync"

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
}

func newProducer(client *client, conf *config.ProducerConfig, checkers ...internal.RouteChecker) (*producer, error) {
	options := pulsar.ProducerOptions{
		Topic: conf.Topic,
	}
	pulsarProducer, err := client.Client.CreateProducer(options)
	if err != nil {
		return nil, err
	}
	producer := &producer{
		Producer:    pulsarProducer,
		client:      client,
		logger:      client.logger.SubLogger(log.Fields{"topic": options.Topic}),
		routeEnable: conf.RouteEnable,
		routePolicy: conf.Route,
		checkers:    checkers,
	}
	return producer, nil
}

func (p *producer) Send(ctx context.Context, msg *pulsar.ProducerMessage) (pulsar.MessageID, error) {
	if !p.routeEnable {
		return p.Producer.Send(ctx, msg)
	}
	for _, chk := range p.checkers {
		routeTopic := chk(msg)
		if routeTopic == "" {
			continue
		}
		rtr, err := p.internalGetRouter(routeTopic)
		if err != nil {
			p.logger.Warnf("failed to create router for topic: %s", routeTopic)
			continue
		}
		if !rtr.ready {
			p.logger.Warnf("router is still not ready for topic: %s", routeTopic)
			continue
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
	return p.Producer.Send(ctx, msg)
}

func (p *producer) SendAsync(ctx context.Context, msg *pulsar.ProducerMessage,
	callback func(pulsar.MessageID, *pulsar.ProducerMessage, error)) {
	if !p.routeEnable {
		p.Producer.SendAsync(ctx, msg, callback)
		return
	}
	for _, chk := range p.checkers {
		routeTopic := chk(msg)
		if routeTopic == "" {
			continue
		}
		rtr, err := p.internalGetRouter(routeTopic)
		if err != nil {
			p.logger.Warnf("failed to create router for topic: %s", routeTopic)
			continue
		}
		if !rtr.ready {
			p.logger.Warnf("router is still not ready for topic: %s", routeTopic)
			continue
		}
		rtr.SendAsync(ctx, msg, callback)
	}
	p.Producer.SendAsync(ctx, msg, callback)
}

func (p *producer) internalGetRouter(topic string) (*router, error) {
	p.routersLock.RLock()
	rtr, ok := p.routers[topic]
	p.routersLock.RUnlock()
	if !ok {
		options := routerOptions{Topic: topic, connectInSyncEnable: p.routePolicy.ConnectInSyncEnable}
		p.routersLock.Lock()
		defer p.routersLock.Unlock()
		rtr, ok = p.routers[topic]
		if !ok {
			if newRtr, err := newRouter(p.logger, p.client, options); err != nil {
				return nil, err
			} else {
				rtr = newRtr
				p.routers[topic] = newRtr
			}
		}
	}
	return rtr, nil
}