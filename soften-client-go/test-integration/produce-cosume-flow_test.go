package test_integration

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/shenqianjin/soften-client-go/soften"
	"github.com/shenqianjin/soften-client-go/soften/config"
)

func TestNackConsumeTimesHeader(t *testing.T) {
	client, err := soften.NewClient(config.ClientConfig{
		URL: "pulsar://localhost:6650",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	producer, err := client.CreateSoftenProducer(config.ProducerConfig{
		Topic: "my-topic",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	consumer, err := client.SubscribeRegular(config.ConsumerConfig{
		Topic:               "my-topic",
		SubscriptionName:    "my-sub",
		Type:                pulsar.Shared,
		NackRedeliveryDelay: 2 * time.Second,
	}, func(message pulsar.Message) (bool, error) {
		fmt.Printf("consume message: %v, headers: %v\n", string(message.Payload()), message.Properties())
		return true, nil
	})
	if err != nil {
		fmt.Printf("failed to subscribe. err: %v\n", err)
	}
	defer consumer.Close()

	for count := 0; count < 100; count++ {
		_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
			Payload: []byte(fmt.Sprintf("hello message index: %v at %v", count, time.Now().Format(time.RFC3339))),
		})
		if err != nil {
			fmt.Printf("failed to send message. err: %v\n", err)
		}
	}

	fmt.Println("starting....")

	time.Sleep(10 * time.Second)

}
