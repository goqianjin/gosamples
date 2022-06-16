package internal

import "github.com/apache/pulsar-client-go/pulsar"

// ------ route checker ------

type RouteChecker func(*pulsar.ProducerMessage) string

var NilRouteChecker = func(*pulsar.ProducerMessage) string {
	return ""
}
