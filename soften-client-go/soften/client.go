package soften

import (
	"errors"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
	"github.com/shenqianjin/soften-client-go/soften/checker"
	"github.com/shenqianjin/soften-client-go/soften/config"
	"github.com/shenqianjin/soften-client-go/soften/internal"
)

type client struct {
	pulsar.Client
	logger log.Logger
}

func NewClient(conf config.ClientConfig) (*client, error) {
	// validate and default conf
	if err := config.Validator.ValidateAndDefaultClientConfig(&conf); err != nil {
		return nil, err
	}
	// create client
	clientOption := pulsar.ClientOptions{
		URL:                     conf.URL,
		ConnectionTimeout:       time.Duration(conf.ConnectionTimeout) * time.Second,
		OperationTimeout:        time.Duration(conf.ConnectionTimeout) * time.Second,
		MaxConnectionsPerBroker: int(conf.MaxConnectionsPerBroker),
		Logger:                  conf.Logger,
	}
	pulsarClient, err := pulsar.NewClient(clientOption)
	if err != nil {
		return nil, err
	}
	cli := &client{Client: pulsarClient, logger: conf.Logger}
	return cli, nil
}

func (c *client) CreateSoftenProducer(conf config.ProducerConfig, checkers ...internal.RouteChecker) (*producer, error) {
	if conf.Topic == "" {
		return nil, errors.New("topic is empty")
	}
	return newProducer(c, &conf, checkers...)

}

func (c *client) SubscribeRegular(conf config.ConsumerConfig, handler Handler, checkpoints ...internal.Checkpoint) (*consumeFacade, error) {
	// convert handler
	handlerWithState := func(message pulsar.Message) HandleStatus {
		success, err := handler(message)
		if success {
			return HandleStatusOk
		} else {
			return HandleStatusFail.Err(err)
		}
	}
	// forward the call to c.SubscribePremium
	return c.SubscribePremium(conf, handlerWithState, checkpoints...)

}

func (c *client) SubscribePremium(conf config.ConsumerConfig, handler PremiumHandler, checkpoints ...internal.Checkpoint) (*consumeFacade, error) {
	// validate and default config
	if err := config.Validator.ValidateAndDefaultConsumerConfig(&conf); err != nil {
		return nil, err
	}
	// validate handler
	if handler == nil {
		panic("handler parameter is nil")
	}
	// validate checkpoints
	checkpointMap, err := checker.Validator.ValidateConsumeCheckpoint(&conf, checkpoints...)
	if err != nil {
		return nil, err
	}
	// create consumer
	if consumer, err := newMultiStatusConsumeFacade(c, conf, handler, checkpointMap); err != nil {
		return nil, err
	} else {
		return consumer, err
	}
}

func (c *client) SubscribeMultiLevel(conf config.MultiLevelConsumerConfig, handler PremiumHandler, checkpoints ...internal.Checkpoint) (*consumeFacade, error) {
	// validate and default config
	if err := config.Validator.ValidateAndDefaultMultiLevelConsumerConfig(&conf); err != nil {
		return nil, err
	}
	// validate handler
	if handler == nil {
		panic("handler parameter is nil")
	}
	// validate checkpoints
	checkpointMap, err := checker.Validator.ValidateConsumeCheckpoint(conf.ConsumerConfig, checkpoints...)
	if err != nil {
		return nil, err
	}
	// create consumer
	if consumer, err := newMultiLevelConsumeFacade(c, conf, handler, checkpointMap); err != nil {
		return nil, err
	} else {
		return consumer, err
	}
}
