package main

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

func TestNackConsumeLaterMix(t *testing.T) {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: "pulsar://localhost:6650",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: "my-topic",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	consumer, _ := client.Subscribe(pulsar.ConsumerOptions{
		Topic:               "my-topic",
		SubscriptionName:    "my-sub",
		Type:                pulsar.Shared,
		NackRedeliveryDelay: 1 * time.Second,
		RetryEnable:         true,
		DLQ: &pulsar.DLQPolicy{
			MaxDeliveries: 11,
		},
	})
	defer consumer.Close()

	_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload: []byte("hello"),
	})
	if err != nil {
		fmt.Println("Failed to publish message", err)
	} else {
		fmt.Println("Published message")
	}
	for count := 1; count <= 1000; count++ {
		msg, err := consumer.Receive(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("count: %v, redeliveryCount: %v, reconsumeTimes: %v, Received message msgId: %#v -- content: '%s'\n",
			count, msg.RedeliveryCount(), msg.Properties()[pulsar.SysPropertyReconsumeTimes], msg.ID(), string(msg.Payload()))
		if count%10 == 0 {
			consumer.ReconsumeLater(msg, time.Second)
		}
		if count/10 == 10 {
			consumer.Ack(msg)
			break
		}
		consumer.Nack(msg)
	}

}
