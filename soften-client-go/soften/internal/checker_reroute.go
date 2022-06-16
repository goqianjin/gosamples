package internal

import "github.com/apache/pulsar-client-go/pulsar"

// ------ reroute checker ------

type PreRerouteChecker func(pulsar.Message) string

type PostRerouteChecker func(pulsar.Message, error) string

var NilPreRerouteChecker = func(pulsar.Message) string {
	return ""
}

var NilPostRerouteChecker = func(pulsar.Message, error) string {
	return ""
}
