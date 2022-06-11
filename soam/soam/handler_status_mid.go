package soam

import (
	"strconv"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

type statusHandleOptions struct {
	topic       string        // default ${TOPIC}_RETRYING, 固定后缀，不允许定制
	enable      bool          // 内部判断使用
	status      messageStatus // messageStatus
	deadHandler internalHandler
}

type statusHandler struct {
	router  *router
	policy  *StatusPolicy
	options statusHandleOptions
}

func newStatusHandler(client *client, policy *StatusPolicy, options statusHandleOptions) (*statusHandler, error) {
	rt, err := newRouter(client.logger, client, routerOption{Topic: options.topic, Enable: options.enable})
	if err == nil {
		return nil, err
	}
	statusRouter := &statusHandler{router: rt, policy: policy, options: options}
	return statusRouter, nil
}

func (sr *statusHandler) Handle(msg pulsar.ConsumerMessage) bool {
	statusReconsumeTimes := MessageParser.GetStatusReconsumeTimes(sr.options.status, msg)
	// check to dead if exceed max status reconsume times
	if statusReconsumeTimes >= sr.policy.ConsumeMaxTimes {
		return sr.tryDeadInternal(msg)
	}
	statusReentrantTimes := MessageParser.GetStatusReentrantTimes(sr.options.status, msg)
	// check to dead if exceed max reentrant times
	if statusReentrantTimes >= sr.policy.ReentrantMaxTimes {
		return sr.tryDeadInternal(msg)
	}
	currentStatus := MessageParser.GetCurrentStatus(msg)
	delay := uint(0)
	// check Nack for equal status
	if currentStatus == sr.options.status {
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
	if currentStatus != sr.options.status {
		// first time to happen status switch
		previousMessageStatus := MessageParser.GetPreviousStatus(msg)
		if (previousMessageStatus == "" || previousMessageStatus == MessageStatusReady) && sr.options.status != MessageStatusReady {
			// record origin information when re-route first time
			if _, ok := props[XPropertyOriginTopic]; !ok {
				props[XPropertyOriginTopic] = msg.Message.Topic()
			}
			if _, ok := props[XPropertyOriginMessageID]; !ok {
				props[XPropertyOriginMessageID] = MessageParser.GetMessageId(msg)
			}
			if _, ok := props[XPropertyOriginPublishTime]; !ok {
				props[XPropertyOriginPublishTime] = msg.PublishTime().Format(RFC3339TimeInSecondPattern)
			}
		}
		props[XPropertyPreviousMessageStatus] = string(currentStatus)
		delay = sr.policy.ReentrantDelay // default a newStatus.reentrantDelay if status switch happens
	}
	now := time.Now()
	reentrantStartRedeliveryCount := MessageParser.GetReentrantStartRedeliveryCount(msg)
	props[XPropertyReentrantStartRedeliveryCount] = strconv.FormatUint(uint64(msg.RedeliveryCount()), 10)

	xReconsumeTimes := MessageParser.GetXReconsumeTimes(msg)
	xReconsumeTimes++
	props[XPropertyReconsumeTimes] = strconv.Itoa(xReconsumeTimes) // initialize continuous consume times for the new msg

	props[XPropertyReconsumeTime] = now.Add(time.Duration(delay) * time.Second).Format(RFC3339TimeInSecondPattern)
	props[XPropertyReentrantTime] = now.Add(time.Duration(sr.policy.ReentrantDelay) * time.Second).Format(RFC3339TimeInSecondPattern)

	if statusReconsumeTimesHeader, ok := statusConsumeTimesMap[sr.options.status]; ok {
		statusReconsumeTimes += int(msg.RedeliveryCount() - reentrantStartRedeliveryCount) // the subtraction is the nack times in current status
		props[statusReconsumeTimesHeader] = strconv.Itoa(statusReconsumeTimes)
	}
	if statusReentrantTimesHeader, ok := statusReentrantTimesMap[sr.options.status]; ok {
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
	sr.router.Chan() <- &RerouteMessage{
		consumerMsg: msg,
		producerMsg: producerMsg,
	}
	return true
}

func (sr *statusHandler) tryDeadInternal(msg pulsar.ConsumerMessage) bool {
	if sr.options.deadHandler != nil {
		return sr.options.deadHandler.Handle(msg)
	}
	return false
}
