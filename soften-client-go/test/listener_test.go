package test

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/shenqianjin/soften-client-go/soften/admin"
	"github.com/shenqianjin/soften-client-go/soften/checker"
	"github.com/shenqianjin/soften-client-go/soften/config"
	"github.com/shenqianjin/soften-client-go/soften/message"
	"github.com/shenqianjin/soften-client-go/soften/topic"
	"github.com/stretchr/testify/assert"
)

func TestListen_1Msg_Ready(t *testing.T) {
	testListenBySingleStatus(t, string(message.StatusReady))
}

func TestListen_1Msg_Retrying(t *testing.T) {
	testListenBySingleStatus(t, string(message.StatusRetrying))
}

func TestListen_1Msg_Pending(t *testing.T) {
	testListenBySingleStatus(t, string(message.StatusPending))
}

func TestListen_1Msg_Blocking(t *testing.T) {
	testListenBySingleStatus(t, string(message.StatusBlocking))
}

func testListenBySingleStatus(t *testing.T, status string) {
	topic := generateTestTopic()
	storedTopic := topic
	if status != string(message.StatusReady) {
		storedTopic = topic + "-" + strings.ToUpper(status)
	}
	manager := admin.NewAdminManager(defaultPulsarHttpUrl)
	// clean up topic
	err := manager.Delete(storedTopic)
	assert.True(t, err == nil || strings.Contains(err.Error(), "404 Not Found"))
	defer func() {
		err = manager.Delete(storedTopic)
		assert.True(t, err == nil || strings.Contains(err.Error(), "404 Not Found"))
	}()
	// create client
	client := NewClient(defaultPulsarUrl)
	defer client.Close()
	// create producer
	producer, err := client.CreateProducer(config.ProducerConfig{
		Topic:       topic,
		RouteEnable: true,
		Route:       &config.RoutePolicy{ConnectInSyncEnable: true},
	}, checker.ProduceRouteChecker(func(msg *pulsar.ProducerMessage) checker.CheckStatus {
		if storedTopic, ok := msg.Properties["routeTopic"]; ok {
			return checker.CheckStatusPassed.WithRerouteTopic(storedTopic)
		}
		return checker.CheckStatusRejected
	}))
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()
	// send messages
	var msgID pulsar.MessageID
	if status == string(message.StatusReady) {
		msgID, err = producer.Send(context.Background(), GenerateProduceMessage(size1K))
	} else {
		msgID, err = producer.Send(context.Background(), GenerateProduceMessage(size1K, "routeTopic", storedTopic))
	}
	assert.True(t, err == nil)
	fmt.Println("produced message: ", msgID)
	// check send stats
	stats, err := manager.Stats(storedTopic)
	assert.True(t, err == nil)
	assert.True(t, stats.MsgInCounter == 1)

	// ---------------

	// create listener
	listener := CreateListener(client, config.ConsumerConfig{
		Topic:                       topic,
		SubscriptionName:            generateSubscribeNameByTopic(topic),
		SubscriptionInitialPosition: pulsar.SubscriptionPositionEarliest,
		RetryingEnable:              string(message.StatusRetrying) == status, // enable retrying if matches
		PendingEnable:               string(message.StatusPending) == status,  // enable pending if matches
		BlockingEnable:              string(message.StatusBlocking) == status, // enable blocking if matches
	})
	defer listener.Close()
	// listener starts
	ctx, cancel := context.WithCancel(context.Background())
	err = listener.Start(ctx, func(message pulsar.Message) (bool, error) {
		fmt.Printf("consumed message size: %v, headers: %v\n", len(message.Payload()), message.Properties())
		return true, nil
	})
	if err != nil {
		log.Fatal(err)
	}
	// wait for consuming the message
	time.Sleep(30 * time.Millisecond)
	// check stats
	stats, err = manager.Stats(storedTopic)
	assert.True(t, err == nil)
	assert.True(t, stats.MsgOutCounter == 1)
	assert.True(t, stats.MsgOutCounter == stats.MsgInCounter)
	// stop listener
	cancel()
}

func TestListen_4Msg_Ready_Retrying_Pending_Blocking(t *testing.T) {
	topic := generateTestTopic()
	// format topics
	statuses := []string{string(message.StatusReady), string(message.StatusRetrying), string(message.StatusPending), string(message.StatusBlocking)}
	storedTopics := make([]string, len(statuses))
	for index, status := range statuses {
		statusTopic := topic
		if status != string(message.StatusReady) {
			statusTopic = topic + "-" + strings.ToUpper(status)
		}
		storedTopics[index] = statusTopic
	}

	manager := admin.NewAdminManager(defaultPulsarHttpUrl)
	// clean up topic
	for _, storedTopic := range storedTopics {
		err := manager.Delete(storedTopic)
		assert.True(t, err == nil || strings.Contains(err.Error(), "404 Not Found"))
	}
	defer func() {
		for _, storedTopic := range storedTopics {
			err := manager.Delete(storedTopic)
			assert.True(t, err == nil || strings.Contains(err.Error(), "404 Not Found"))
		}
	}()
	// create client
	client := NewClient(defaultPulsarUrl)
	defer client.Close()
	// create producer
	producer, err := client.CreateProducer(config.ProducerConfig{
		Topic:       topic,
		RouteEnable: true,
		Route:       &config.RoutePolicy{ConnectInSyncEnable: true},
	}, checker.ProduceRouteChecker(func(msg *pulsar.ProducerMessage) checker.CheckStatus {
		if storedTopic, ok := msg.Properties["routeTopic"]; ok {
			return checker.CheckStatusPassed.WithRerouteTopic(storedTopic)
		}
		return checker.CheckStatusRejected
	}))
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()
	// send messages
	for _, storedTopic := range storedTopics {
		var msgID pulsar.MessageID
		if storedTopic == topic { // ready status
			msgID, err = producer.Send(context.Background(), GenerateProduceMessage(size1K))
		} else {
			msgID, err = producer.Send(context.Background(), GenerateProduceMessage(size1K, "routeTopic", storedTopic))
		}
		assert.True(t, err == nil)
		fmt.Println("produced message: ", msgID)
		// check send stats
		stats, err := manager.Stats(storedTopic)
		assert.True(t, err == nil)
		assert.True(t, stats.MsgInCounter == 1)
	}

	// ---------------

	// create listener
	listener := CreateListener(client, config.ConsumerConfig{
		Topic:                       topic,
		SubscriptionName:            generateSubscribeNameByTopic(topic),
		SubscriptionInitialPosition: pulsar.SubscriptionPositionEarliest,
		RetryingEnable:              true, // enable retrying if matches
		PendingEnable:               true, // enable pending if matches
		BlockingEnable:              true, // enable blocking if matches
	})
	defer listener.Close()
	// listener starts
	ctx, cancel := context.WithCancel(context.Background())
	err = listener.Start(ctx, func(message pulsar.Message) (bool, error) {
		fmt.Printf("consumed message size: %v, headers: %v\n", len(message.Payload()), message.Properties())
		return true, nil
	})
	if err != nil {
		log.Fatal(err)
	}
	// wait for consuming the message
	time.Sleep(100 * time.Millisecond)
	// check stats
	for _, storedTopic := range storedTopics {
		stats, err := manager.Stats(storedTopic)
		assert.True(t, err == nil)
		assert.True(t, stats.MsgOutCounter == 1)
		assert.True(t, stats.MsgOutCounter == stats.MsgInCounter)
	}
	// stop listener
	cancel()
}

func TestListen_2Msg_L2(t *testing.T) {
	testListenByMultiLevels(t, topic.Levels{topic.L2})
}

func TestListen_2Msg_L1_L2(t *testing.T) {
	testListenByMultiLevels(t, topic.Levels{topic.L1, topic.L2})
}

func TestListen_2Msg_B1(t *testing.T) {
	testListenByMultiLevels(t, topic.Levels{topic.B1})
}

func TestListen_2Msg_L1_B1(t *testing.T) {
	testListenByMultiLevels(t, topic.Levels{topic.L1, topic.B1})
}

func TestListen_2Msg_S1(t *testing.T) {
	testListenByMultiLevels(t, topic.Levels{topic.S1})
}

func TestListen_2Msg_L1_S1(t *testing.T) {
	testListenByMultiLevels(t, topic.Levels{topic.L1, topic.S1})
}

func TestListen_4Msg_L1_L2_B1_S1(t *testing.T) {
	testListenByMultiLevels(t, topic.Levels{topic.L1, topic.L2, topic.B1, topic.S1})
}

func TestListen_Msg_AllLevels(t *testing.T) {
	testListenByMultiLevels(t, topic.Levels{
		topic.S2, topic.S1,
		topic.L3, topic.L2, topic.L1,
		topic.B1, topic.B2,
	})
}

func testListenByMultiLevels(t *testing.T, levels topic.Levels) {
	testTopic := generateTestTopic()
	// format topics
	storedTopics := make([]string, len(levels))
	for index, level := range levels {
		storedTopic := testTopic
		if level != topic.L1 {
			storedTopic = testTopic + "-" + strings.ToUpper(level.String())
		}
		storedTopics[index] = storedTopic
	}

	manager := admin.NewAdminManager(defaultPulsarHttpUrl)
	// clean up testTopic
	for _, storedTopic := range storedTopics {
		err := manager.Delete(storedTopic)
		assert.True(t, err == nil || strings.Contains(err.Error(), "404 Not Found"))
	}
	defer func() {
		for _, storedTopic := range storedTopics {
			err := manager.Delete(storedTopic)
			assert.True(t, err == nil || strings.Contains(err.Error(), "404 Not Found"))
		}
	}()
	// create client
	client := NewClient(defaultPulsarUrl)
	defer client.Close()
	// create producer
	producer, err := client.CreateProducer(config.ProducerConfig{
		Topic:       testTopic,
		RouteEnable: true,
		Route:       &config.RoutePolicy{ConnectInSyncEnable: true},
	}, checker.ProduceRouteChecker(func(msg *pulsar.ProducerMessage) checker.CheckStatus {
		if storedTopic, ok := msg.Properties["routeTopic"]; ok {
			return checker.CheckStatusPassed.WithRerouteTopic(storedTopic)
		}
		return checker.CheckStatusRejected
	}))
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()
	// send messages
	for index, level := range levels {
		storedTopic := storedTopics[index]
		var msgID pulsar.MessageID
		if level == topic.L1 { // ready level
			msgID, err = producer.Send(context.Background(), GenerateProduceMessage(size1K))
		} else {
			msgID, err = producer.Send(context.Background(), GenerateProduceMessage(size1K, "routeTopic", storedTopic))
		}
		assert.True(t, err == nil)
		fmt.Println("produced message: ", msgID)
		// check send stats
		stats, err := manager.Stats(storedTopic)
		assert.True(t, err == nil)
		assert.True(t, stats.MsgInCounter == 1)
	}

	// ---------------

	// create listener
	listener := CreateListener(client, config.ConsumerConfig{
		Topic:                       testTopic,
		SubscriptionName:            generateSubscribeNameByTopic(testTopic),
		SubscriptionInitialPosition: pulsar.SubscriptionPositionEarliest,
		//RetryingEnable:              true, // enable retrying if matches
		//PendingEnable:               true, // enable pending if matches
		//BlockingEnable:              true, // enable blocking if matches
		Levels: levels,
	})
	defer listener.Close()
	// listener starts
	ctx, cancel := context.WithCancel(context.Background())
	err = listener.Start(ctx, func(message pulsar.Message) (bool, error) {
		fmt.Printf("consumed message size: %v, headers: %v\n", len(message.Payload()), message.Properties())
		return true, nil
	})
	if err != nil {
		log.Fatal(err)
	}
	// wait for consuming the message
	time.Sleep(100 * time.Millisecond)
	// check stats
	for _, storedTopic := range storedTopics {
		stats, err := manager.Stats(storedTopic)
		assert.True(t, err == nil)
		assert.True(t, stats.MsgOutCounter == 1)
		assert.True(t, stats.MsgOutCounter == stats.MsgInCounter)
	}
	// stop listener
	cancel()
}
