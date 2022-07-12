package test

import (
	"crypto/rand"
	"log"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/shenqianjin/soften-client-go/soften"
	"github.com/shenqianjin/soften-client-go/soften/config"
)

var (
	defaultPulsarUrl     = "pulsar://localhost:6650"
	defaultPulsarHttpUrl = "http://localhost:8080"

	defaultTopic = "my-topic"
	defaultSub   = "my-sub"

	size64 = 64
	size1K = 1024
	size2K = size1K * 2
	size3K = size1K * 3
	size4K = size1K * 4
	size5K = size1K * 5
	size1M = size1K * 1024
	size1G = size1M * 1024

	timeFormat = "20060102150405"

	topicCount = int32(0)
)

func generateTestTopic() string {
	index := atomic.AddInt32(&topicCount, 1)
	now := time.Now().Format(timeFormat)
	return strings.Join([]string{defaultTopic, now, strconv.Itoa(int(index))}, "-")
}

func generateSubscribeName() string {
	return generateSubscribeNameByTopic(generateTestTopic())
}

func generateSubscribeNameByTopic(topic string) string {
	return topic + "-sub"
}

func NewClient(url string) soften.Client {
	if url == "" {
		url = defaultPulsarUrl
	}
	client, err := soften.NewClient(config.ClientConfig{
		URL:               url,
		ConnectionTimeout: 1,
	})
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func CreateProducer(client soften.Client, topic string) soften.Producer {
	if topic == "" {
		topic = generateTestTopic()
	}
	producer, err := client.CreateProducer(config.ProducerConfig{
		Topic: topic,
	})
	if err != nil {
		log.Fatal(err)
	}
	return producer
}

func CreateListener(client soften.Client, conf config.ConsumerConfig) soften.Listener {
	if conf.Topic == "" {
		conf.Topic = generateTestTopic()
	}
	if conf.SubscriptionName == "" {
		conf.SubscriptionName = generateSubscribeNameByTopic(conf.Topic)
	}
	listener, err := client.CreateListener(conf)
	if err != nil {
		log.Fatal(err)
	}
	return listener
}

func GenerateProduceMessage(size int, kvs ...string) *pulsar.ProducerMessage {
	if size < 0 {
		size = size1K
	}
	data := make([]byte, size)
	rand.Read(data)
	properties := make(map[string]string)
	for k := 0; k < len(kvs)-1; k += 2 {
		properties[kvs[k]] = kvs[k+1]
	}
	return &pulsar.ProducerMessage{
		Payload:    data,
		Properties: properties,
	}
}
