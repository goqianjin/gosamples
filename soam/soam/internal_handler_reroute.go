package soam

import (
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
)

type multiTopicsRouter struct {
	*router
	policy *ReRoutePolicy
}

func newMultiTopicsRouter(logger log.Logger, client pulsar.Client, policy *ReRoutePolicy) (*multiTopicsRouter, error) {
	router, err := newRouter(logger, client, routerOption{})
	if err == nil {
		return nil, err
	}
	reRouter := &multiTopicsRouter{router: router, policy: policy}
	return reRouter, nil
}

func (r *multiTopicsRouter) RouteTo(message pulsar.ConsumerMessage, topic string) bool {
	return false
}
