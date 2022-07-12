package checker

import "github.com/apache/pulsar-client-go/pulsar"

// ------ check func (for consumer handle message) ------

type BeforeCheckFunc func(pulsar.Message) CheckStatus

type AfterCheckFunc func(pulsar.Message, error) CheckStatus

var NilBeforeCheckFunc = func(pulsar.Message) CheckStatus {
	return CheckStatusRejected
}

var NilAfterCheckFunc = func(pulsar.Message, error) CheckStatus {
	return CheckStatusRejected
}

// ------ intercept func (for producer send message) ------

type ProduceCheckFunc func(msg *pulsar.ProducerMessage) CheckStatus
