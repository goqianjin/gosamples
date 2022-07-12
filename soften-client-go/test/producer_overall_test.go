package test

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/shenqianjin/soften-client-go/soften/admin"
	"github.com/shenqianjin/soften-client-go/soften/checker"
	"github.com/shenqianjin/soften-client-go/soften/config"
	"github.com/stretchr/testify/assert"
)

func TestProduceOverall_Send3Msg_Route1MsgToL1(t *testing.T) {
	topic := generateTestTopic()
	routedTopic := topic + "-L2"
	manager := admin.NewAdminManager(defaultPulsarHttpUrl)

	err := manager.Delete(topic)
	assert.True(t, err == nil || strings.Contains(err.Error(), "404 Not Found"))
	err = manager.Delete(routedTopic)
	assert.True(t, err == nil || strings.Contains(err.Error(), "404 Not Found"))
	defer func() {
		err = manager.Delete(topic)
		assert.True(t, err == nil || strings.Contains(err.Error(), "404 Not Found"))
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
			if index, ok := msg.Properties["Index"]; ok && index == "2" {
				return checker.CheckStatusPassed.WithRerouteTopic(routedTopic)
			}
			return checker.CheckStatusRejected
		}))
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	for i := 1; i <= 3; i++ {
		msg := GenerateProduceMessage(size64)
		msg.Properties["Index"] = strconv.Itoa(i)
		msgID, err := producer.Send(context.Background(), msg)
		assert.Nil(t, err)
		fmt.Println(msgID)
	}

	stats, err := manager.Stats(topic)
	assert.Nil(t, err)
	assert.Equal(t, 2, stats.MsgInCounter)

	stats, err = manager.Stats(routedTopic)
	assert.Nil(t, err)
	assert.Equal(t, 1, stats.MsgInCounter)
}

func TestProduceOverall_SendAsync3Msg_Route1MsgToL1(t *testing.T) {
	topic := generateTestTopic()
	routedTopic := topic + "-L2"
	manager := admin.NewAdminManager(defaultPulsarHttpUrl)

	err := manager.Delete(topic)
	assert.True(t, err == nil || strings.Contains(err.Error(), "404 Not Found"))
	err = manager.Delete(routedTopic)
	assert.True(t, err == nil || strings.Contains(err.Error(), "404 Not Found"))
	defer func() {
		err = manager.Delete(topic)
		assert.True(t, err == nil || strings.Contains(err.Error(), "404 Not Found"))
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
			if index, ok := msg.Properties["Index"]; ok && index == "2" {
				return checker.CheckStatusPassed.WithRerouteTopic(routedTopic)
			}
			return checker.CheckStatusRejected
		}))
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	wg := sync.WaitGroup{}
	wg.Add(3)
	for i := 1; i <= 3; i++ {
		msg := GenerateProduceMessage(size64)
		msg.Properties["Index"] = strconv.Itoa(i)
		producer.SendAsync(context.Background(), msg,
			func(id pulsar.MessageID, message *pulsar.ProducerMessage, err error) {
				fmt.Println("sent async message: ", id)
				wg.Done()
			})
	}
	wg.Wait()

	stats, err := manager.Stats(topic)
	assert.Nil(t, err)
	assert.Equal(t, 2, stats.MsgInCounter)

	stats, err = manager.Stats(routedTopic)
	assert.Nil(t, err)
}
