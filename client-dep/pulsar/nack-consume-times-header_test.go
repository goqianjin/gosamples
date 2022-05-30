package main

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

func TestNackConsumeTimesHeader(t *testing.T) {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: "pulsar://localhost:6650",
	})
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	producer, _ := client.CreateProducer(pulsar.ProducerOptions{
		Topic: "my-topic",
	})

	defer producer.Close()

	consumer, _ := client.Subscribe(pulsar.ConsumerOptions{
		Topic:               "my-topic",
		SubscriptionName:    "my-sub",
		Type:                pulsar.Shared,
		NackRedeliveryDelay: 2 * time.Second,
	})

	defer consumer.Close()

	_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload:    []byte("hello"),
		Properties: map[string]string{pulsar.SysPropertyReconsumeTimes: "10"},
	})
	if err != nil {
		fmt.Println("Failed to publish message", err)
	} else {
		fmt.Println("Published message")
	}
	for i := 0; i <= 2; i++ {
		msg, err := consumer.Receive(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		consumer.Ack(msg)
		fmt.Printf("Received message msgId: %#v -- content: '%s'\n", msg.ID(), string(msg.Payload()))
	}

	consumer2, _ := client.Subscribe(pulsar.ConsumerOptions{
		Topic:               "my-topic",
		SubscriptionName:    "my-sub",
		Type:                pulsar.Shared,
		NackRedeliveryDelay: 2 * time.Second,
	})
	time.Sleep(4 * time.Second)

	fmt.Println(consumer2)
	defer consumer2.Close()
	msg, err := consumer2.Receive(context.Background())
	consumer2.Nack(msg)

	/*time.Sleep(4 * time.Second)
	client2, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: "pulsar://localhost:6650",
	})
	consumer2, _ := client2.Subscribe(pulsar.ConsumerOptions{
		Topic:               "my-topic",
		SubscriptionName:    "my-sub",
		Type:                pulsar.Shared,
		NackRedeliveryDelay: 2 * time.Second,
	})
	msg, err := consumer2.Receive(context.Background())
	consumer2.Ack(msg)*/

}
