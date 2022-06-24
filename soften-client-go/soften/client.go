package soften

import (
	"errors"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
	"github.com/shenqianjin/soften-client-go/soften/config"
	"github.com/shenqianjin/soften-client-go/soften/internal"
)

type Client interface {
	RawClient() pulsar.Client
	CreateProducer(conf config.ProducerConfig, checkers ...internal.RouteChecker) (*producer, error)
	CreateListener(conf config.ConsumerConfig) (*consumeListener, error)
	Close() // close the Client and free associated resources
}

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
	cli.logger.Infof("Soften client (url:%s) is ready", conf.URL)
	return cli, nil
}

func (c *client) RawClient() pulsar.Client {
	return c.Client
}

func (c *client) CreateProducer(conf config.ProducerConfig, checkers ...internal.RouteChecker) (*producer, error) {
	if conf.Topic == "" {
		return nil, errors.New("topic is empty")
	}
	return newProducer(c, &conf, checkers...)

}

func (c *client) CreateListener(conf config.ConsumerConfig) (*consumeListener, error) {
	// validate and default config
	if err := config.Validator.ValidateAndDefaultConsumerConfig(&conf); err != nil {
		return nil, err
	}
	// create consumer
	if consumer, err := newConsumeListener(c, conf); err != nil {
		return nil, err
	} else {
		return consumer, err
	}
}
