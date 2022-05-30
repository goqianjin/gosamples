package main

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

func TestSubscribe(t *testing.T) {

	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: "pulsar://localhost:6650",
	})
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	consumer, _ := client.Subscribe(pulsar.ConsumerOptions{
		Topic:               "my-topic",
		SubscriptionName:    "my-sub",
		Type:                pulsar.Shared,
		NackRedeliveryDelay: 2 * time.Second,
	})

	defer consumer.Close()

	for {
		msg, err := consumer.Receive(context.Background())
		fmt.Println(msg, err)
		consumer.Ack(msg)
		time.Sleep(time.Second)

	}
}
