package internal

import "github.com/apache/pulsar-client-go/pulsar"

// ------ status checker ------

type PreStatusChecker func(pulsar.Message) (passed bool)

type PostStatusChecker func(pulsar.Message, error) (passed bool)

var NilPreStatusChecker = func(pulsar.Message) (passed bool) {
	return false
}

var NilPostStatusChecker = func(pulsar.Message, error) (passed bool) {
	return false
}
