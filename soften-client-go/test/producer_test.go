package test

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/shenqianjin/soften-client-go/soften/admin"
	"github.com/stretchr/testify/assert"
)

func TestSend_1Msg(t *testing.T) {
	topic := generateTestTopic()
	manager := admin.NewAdminManager(defaultPulsarHttpUrl)

	err := manager.Delete(topic)
	assert.True(t, err == nil || strings.Contains(err.Error(), "404 Not Found"))
	defer func() {
		err = manager.Delete(topic)
		assert.True(t, err == nil || strings.Contains(err.Error(), "404 Not Found"))
	}()

	client := NewClient(defaultPulsarUrl)
	defer client.Close()

	producer := CreateProducer(client, topic)
	defer producer.Close()

	msgID, err := producer.Send(context.Background(), GenerateProduceMessage(size1K))
	assert.True(t, err == nil)
	fmt.Println(msgID)

	stats, err := manager.Stats(topic)
	assert.True(t, err == nil)
	assert.True(t, stats.MsgInCounter == 1)
}

func TestSendAsync_1Msg(t *testing.T) {
	topic := generateTestTopic()
	manager := admin.NewAdminManager(defaultPulsarHttpUrl)

	err := manager.Delete(topic)
	assert.True(t, err == nil || strings.Contains(err.Error(), "404 Not Found"))
	defer func() {
		err = manager.Delete(topic)
		assert.True(t, err == nil || strings.Contains(err.Error(), "404 Not Found"))
	}()

	client := NewClient(defaultPulsarUrl)
	defer client.Close()

	producer := CreateProducer(client, topic)
	defer producer.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)
	producer.SendAsync(context.Background(), GenerateProduceMessage(size1K),
		func(id pulsar.MessageID, message *pulsar.ProducerMessage, err error) {
			fmt.Println("sent async message: ", id)
			wg.Done()
		})
	wg.Wait()
	assert.True(t, err == nil)

	stats, err := manager.Stats(topic)
	assert.True(t, err == nil)
	assert.True(t, stats.MsgInCounter == 1)
}
