package lib_smc

import "time"

type ClientConfig struct {
	URL               string
	ConnectionTimeout time.Duration
}

type ProducerConfig struct {
}

type ComsumerConfig struct {
	Topic string

	PendingNackBackoffPolicy   []string
	PendingNackBackoffMaxTimes int
	//PendingTo                  string // default ${TOPIC}_PENDING

	BlockingNackBackoffPolicy   []string
	BlockingNackBackoffMaxTimes int
	//BlockingTo                  string // default ${TOPIC}_BLOCKING

	RetryingNackBackoffPolicy   []string
	RetryingNackBackoffMaxTimes int
	//RetryingTo                  string // default ${TOPIC}_RETRYING

	Topics []string

	TopicsPattern string

	AutoDiscoveryPeriod time.Duration

	SubscriptionName string

	Properties map[string]string

	SubscriptionProperties map[string]string

	Type SubscriptionType

	SubscriptionInitialPosition

	DLQ *DLQPolicy

	KeySharedPolicy *KeySharedPolicy

	RetryEnable bool

	MessageChannel chan ConsumerMessage

	ReceiverQueueSize int

	NackRedeliveryDelay time.Duration

	Name string

	ReadCompacted bool

	ReplicateSubscriptionState bool

	Interceptors ConsumerInterceptors

	Schema Schema

	MaxReconnectToBroker *uint

	Decryption *MessageDecryptionInfo

	EnableDefaultNackBackoffPolicy bool

	NackBackoffPolicy NackBackoffPolicy
}
