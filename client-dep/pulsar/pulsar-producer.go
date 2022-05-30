package main

import (
	"fmt"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

func main() {
	a := pulsar.Shared
	fmt.Println(a)
	var c pulsar.Consumer
	client, err := pulsar.NewClient(pulsar.ClientOptions{})
	producer, errP := client.CreateProducer(pulsar.ProducerOptions{})
	consumer, errC := client.Subscribe(pulsar.ConsumerOptions{})
	consumer.Receive(nil)

	fmt.Println(c)
	fmt.Println(err)
	fmt.Println(errP)
	fmt.Println(errC)
	//producer.Send()
	producer.LastSequenceID()
	consumer.Ack(nil)
	consumer.ReconsumeLater(nil, time.Minute)
	consumer.Nack(nil)
	//consumer.Nack()
	//consumer.Ack()
	//pulsar.ClientOptions

}
