package checker

import "github.com/apache/pulsar-client-go/pulsar"

// ------ check func ------

type BeforeCheckFunc func(pulsar.Message) CheckStatus

type AfterCheckFunc func(pulsar.Message, error) CheckStatus

var NilBeforeCheckFunc = func(pulsar.Message) CheckStatus {
	return CheckStatusFailed
}

var NilAfterCheckFunc = func(pulsar.Message, error) CheckStatus {
	return CheckStatusFailed
}
