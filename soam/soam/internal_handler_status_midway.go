package soam

import (
	"strconv"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
)

const (
	SysPropertyOriginMessageID   = "X-ORIGIN-MESSAGE-ID"
	SysPropertyOriginTopic       = "X-ORIGIN-TOPIC"
	SysPropertyOriginPublishTime = "X-ORIGIN-PUBLISH-TIME"
)

var RFC3339TimeInSecondPattern = "20060102T150405Z"

type statusHandler struct {
	router   *router
	consumer *consumer
	status   messageStatus
	policy   *StatusPolicy
}

func newStatusHandler(logger log.Logger, consumer *consumer, status messageStatus, policy *StatusPolicy) (*statusHandler, error) {
	rt, err := newRouter(logger, consumer.client, routerOption{})
	if err == nil {
		return nil, err
	}
	statusRouter := &statusHandler{router: rt, policy: policy, status: status}
	return statusRouter, nil
}

func (sr *statusHandler) Handle(msg pulsar.ConsumerMessage) bool {
	statusReconsumeTimes := MessageParser.GetStatusReconsumeTimes(sr.status, msg)
	// check to dead if exceed max status reconsume times
	if statusReconsumeTimes >= sr.policy.ConsumeMaxTimes {
		return sr.tryDeadInternal(msg)
	}
	statusReentrantTimes := MessageParser.GetStatusReentrantTimes(sr.status, msg)
	// check to dead if exceed max reentrant times
	if statusReentrantTimes >= sr.policy.ReentrantMaxTimes {
		return sr.tryDeadInternal(msg)
	}
	currentStatus := MessageParser.GetCurrentStatus(msg)
	delay := uint(0)
	// check Nack for equal status
	if currentStatus == sr.status {
		delay = sr.policy.BackoffPolicy.Next(statusReconsumeTimes)
		if delay < sr.policy.ReentrantDelay { // delay equals or larger than reentrant delay is the essential condition to switch status
			msg.Consumer.Nack(msg.Message)
			return true
		}
	}

	// prepare to re-route
	props := make(map[string]string)
	for k, v := range msg.Properties() {
		props[k] = v
	}
	if currentStatus != sr.status {
		// first time to happen status switch
		previousMessageStatus := MessageParser.GetPreviousStatus(msg)
		if (previousMessageStatus == "" || previousMessageStatus == MessageStatusReady) && sr.status != MessageStatusReady {
			// record origin information when re-route first time
			if _, ok := props[SysPropertyOriginTopic]; !ok {
				props[SysPropertyOriginTopic] = msg.Message.Topic()
			}
			if _, ok := props[SysPropertyOriginMessageID]; !ok {
				props[SysPropertyOriginMessageID] = MessageParser.GetMessageId(msg)
			}
			if _, ok := props[SysPropertyOriginPublishTime]; !ok {
				props[SysPropertyOriginPublishTime] = msg.PublishTime().Format(RFC3339TimeInSecondPattern)
			}
		}
		props[SysPropertyPreviousMessageStatus] = string(currentStatus)
		delay = sr.policy.ReentrantDelay // default a newStatus.reentrantDelay if status switch happens
	}
	now := time.Now()
	reentrantStartRedeliveryCount := MessageParser.GetReentrantStartRedeliveryCount(msg)
	props[SysPropertyReentrantStartRedeliveryCount] = strconv.FormatUint(uint64(msg.RedeliveryCount()), 10)

	xReconsumeTimes := MessageParser.GetXReconsumeTimes(msg)
	xReconsumeTimes++
	props[SysPropertyXReconsumeTimes] = strconv.Itoa(xReconsumeTimes) // initialize continuous consume times for the new msg

	props[SysPropertyXReconsumeTime] = now.Add(time.Duration(delay) * time.Second).Format(RFC3339TimeInSecondPattern)
	props[SysPropertyXReentrantTime] = now.Add(time.Duration(sr.policy.ReentrantDelay) * time.Second).Format(RFC3339TimeInSecondPattern)

	if statusReconsumeTimesHeader, ok := statusConsumeTimesMap[sr.status]; ok {
		statusReconsumeTimes += int(msg.RedeliveryCount() - reentrantStartRedeliveryCount) // the subtraction is the nack times in current status
		props[statusReconsumeTimesHeader] = strconv.Itoa(statusReconsumeTimes)
	}
	if statusReentrantTimesHeader, ok := statusReentrantTimesMap[sr.status]; ok {
		statusReentrantTimes++
		props[statusReentrantTimesHeader] = strconv.Itoa(statusReentrantTimes)
	}
	producerMsg := pulsar.ProducerMessage{
		Payload:     msg.Payload(),
		Key:         msg.Key(),
		OrderingKey: msg.OrderingKey(),
		Properties:  props,
		EventTime:   msg.EventTime(),
	}
	sr.router.Chan() <- &ReRouterMessage{
		consumerMsg: msg,
		producerMsg: producerMsg,
	}
	return true
}

func (sr *statusHandler) tryDeadInternal(msg pulsar.ConsumerMessage) bool {
	if sr.consumer.checkers.config.DeadEnable {
		return sr.consumer.checkers.deadHandler.Handle(msg)
	}
	return false
}
