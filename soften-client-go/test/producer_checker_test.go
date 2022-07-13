package test

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"testing"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/shenqianjin/soften-client-go/soften/admin"
	"github.com/shenqianjin/soften-client-go/soften/checker"
	"github.com/shenqianjin/soften-client-go/soften/config"
	"github.com/stretchr/testify/assert"
)

func TestProduceCheck_RouteToL2_BySend(t *testing.T) {
	topic := generateTestTopic()
	routedTopic := topic + "-L2"

	testProduceCheckBySend(t, topic, routedTopic,
		checker.ProduceRouteChecker(func(msg *pulsar.ProducerMessage) checker.CheckStatus {
			return checker.CheckStatusPassed.WithRerouteTopic(routedTopic)
		}))
}

func TestProduceCheck_RouteToL2_BySend(t *testing.T) {
	topic := generateTestTopic()
	routedTopic := topic + "-L2"

	testProduceCheckBySend(t, topic, routedTopic,
		checker.ProduceRouteChecker(func(msg *pulsar.ProducerMessage) checker.CheckStatus {
			return checker.CheckStatusPassed.WithRerouteTopic(routedTopic)
		}))
}

func testProduceCheckBySend(t *testing.T, topic, storedTopic string, checkers ...checker.ProduceCheckpoint) {
	if storedTopic == "" {
		storedTopic = topic
	}
	manager := admin.NewAdminManager(defaultPulsarHttpUrl)

	err := manager.Delete(storedTopic)
	assert.True(t, err == nil || strings.Contains(err.Error(), "404 Not Found"))
	defer func() {
		err = manager.Delete(storedTopic)
		assert.True(t, err == nil || strings.Contains(err.Error(), "404 Not Found"))
	}()

	client := NewClient(defaultPulsarUrl)
	defer client.Close()

	producer, err := client.CreateProducer(config.ProducerConfig{
		Topic:       topic,
		RouteEnable: true,
		Route:       &config.RoutePolicy{ConnectInSyncEnable: true},
	}, checkers...)
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	msg := GenerateProduceMessage(size64)
	msgID, err := producer.Send(context.Background(), msg)
	assert.Nil(t, err)
	fmt.Println(msgID)

	stats, err := manager.Stats(storedTopic)
	assert.Nil(t, err)
	assert.Equal(t, 1, stats.MsgInCounter)
}

func TestProduceCheck_Route1MsgToL1_BySendAsync(t *testing.T) {
	topic := generateTestTopic()
	routedTopic := topic + "-L2"
	manager := admin.NewAdminManager(defaultPulsarHttpUrl)

	err := manager.Delete(routedTopic)
	assert.True(t, err == nil || strings.Contains(err.Error(), "404 Not Found"))
	defer func() {
		err = manager.Delete(routedTopic)
		assert.True(t, err == nil || strings.Contains(err.Error(), "404 Not Found"))
	}()

	client := NewClient(defaultPulsarUrl)
	defer client.Close()

	producer, err := client.CreateProducer(config.ProducerConfig{
		Topic:       topic,
		RouteEnable: true,
		Route:       &config.RoutePolicy{ConnectInSyncEnable: true},
	},
		checker.ProduceRouteChecker(func(msg *pulsar.ProducerMessage) checker.CheckStatus {
			return checker.CheckStatusPassed.WithRerouteTopic(routedTopic)
		}))
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)
	msg := GenerateProduceMessage(size64)
	producer.SendAsync(context.Background(), msg,
		func(id pulsar.MessageID, message *pulsar.ProducerMessage, err error) {
			fmt.Println("sent async message: ", id)
			wg.Done()
		})
	wg.Wait()

	stats, err := manager.Stats(routedTopic)
	assert.Nil(t, err)
	assert.Equal(t, 1, stats.MsgInCounter)
}
