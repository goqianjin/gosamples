package soften

import (
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
	"github.com/shenqianjin/soften-client-go/soften/config"
)

type producer struct {
	pulsar.Producer
	client *client
	logger log.Logger
}

func newProducer(client *client, conf *config.ProducerConfig) (*producer, error) {
	options := pulsar.ProducerOptions{
		Topic: conf.Topic,
	}
	pulsarProducer, err := client.Client.CreateProducer(options)
	if err != nil {
		return nil, err
	}
	producer := &producer{
		Producer: pulsarProducer,
		client:   client,
		logger:   client.logger.SubLogger(log.Fields{"topic": options.Topic}),
	}
	return producer, nil
}
