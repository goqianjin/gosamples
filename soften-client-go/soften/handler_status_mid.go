package soften

import (
	"strconv"
	"time"

	"github.com/shenqianjin/soften-client-go/soften/checker"

	"github.com/shenqianjin/soften-client-go/soften/config"
	"github.com/shenqianjin/soften-client-go/soften/message"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/shenqianjin/soften-client-go/soften/internal"
)

type statusHandleOptions struct {
	topic       string                 // default ${TOPIC}_RETRYING, 固定后缀，不允许定制
	status      internal.MessageStatus // MessageStatus
	deadHandler internalHandler        //
	//levels      []TopicLevel    //
	//enable      bool                   // 内部判断使用
}

type statusHandler struct {
	router  *reRouter
	policy  *config.StatusPolicy
	options statusHandleOptions
}

func newStatusHandler(client *client, policy *config.StatusPolicy, options statusHandleOptions) (*statusHandler, error) {
	rt, err := newReRouter(client.logger, client.Client, reRouterOptions{Topic: options.topic})
	if err != nil {
		return nil, err
	}
	statusRouter := &statusHandler{router: rt, policy: policy, options: options}
	return statusRouter, nil
}

func (sr *statusHandler) Handle(msg pulsar.ConsumerMessage, cheStatus checker.CheckStatus) bool {
	if !cheStatus.IsPassed() {
		return false
	}
	statusReconsumeTimes := message.Parser.GetStatusReconsumeTimes(sr.options.status, msg)
	// check to dead if exceed max status reconsume times
	if statusReconsumeTimes >= sr.policy.ConsumeMaxTimes {
		return sr.tryDeadInternal(msg)
	}
	statusReentrantTimes := message.Parser.GetStatusReentrantTimes(sr.options.status, msg)
	// check to dead if exceed max reentrant times
	if statusReentrantTimes >= sr.policy.ReentrantMaxTimes {
		return sr.tryDeadInternal(msg)
	}
	currentStatus := message.Parser.GetCurrentStatus(msg)
	delay := uint(0)
	// check Nack for equal status
	if currentStatus == sr.options.status {
		delay = sr.policy.BackoffPolicy.Next(0, statusReconsumeTimes)
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
		previousMessageStatus := message.Parser.GetPreviousStatus(msg)
		if (previousMessageStatus == "" || previousMessageStatus == message.StatusReady) && sr.options.status != message.StatusReady {
			// record origin information when re-route first time
			if _, ok := props[message.XPropertyOriginTopic]; !ok {
				props[message.XPropertyOriginTopic] = msg.Message.Topic()
			}
			if _, ok := props[message.XPropertyOriginMessageID]; !ok {
				props[message.XPropertyOriginMessageID] = message.Parser.GetMessageId(msg)
			}
			if _, ok := props[message.XPropertyOriginPublishTime]; !ok {
				props[message.XPropertyOriginPublishTime] = msg.PublishTime().Format(internal.RFC3339TimeInSecondPattern)
			}
		}
		props[message.XPropertyPreviousMessageStatus] = string(currentStatus)
		delay = sr.policy.ReentrantDelay // default a newStatus.reentrantDelay if status switch happens
	}
	now := time.Now()
	reentrantStartRedeliveryCount := message.Parser.GetReentrantStartRedeliveryCount(msg)
	props[message.XPropertyReentrantStartRedeliveryCount] = strconv.FormatUint(uint64(msg.RedeliveryCount()), 10)

	xReconsumeTimes := message.Parser.GetXReconsumeTimes(msg)
	xReconsumeTimes++
	props[message.XPropertyReconsumeTimes] = strconv.Itoa(xReconsumeTimes) // initialize continuous consume times for the new msg

	props[message.XPropertyReconsumeTime] = now.Add(time.Duration(delay) * time.Second).Format(internal.RFC3339TimeInSecondPattern)
	props[message.XPropertyReentrantTime] = now.Add(time.Duration(sr.policy.ReentrantDelay) * time.Second).Format(internal.RFC3339TimeInSecondPattern)

	if statusReconsumeTimesHeader, ok := message.XPropertyConsumeTimes(sr.options.status); ok {
		statusReconsumeTimes += int(msg.RedeliveryCount() - reentrantStartRedeliveryCount) // the subtraction is the nack times in current status
		props[statusReconsumeTimesHeader] = strconv.Itoa(statusReconsumeTimes)
	}
	if statusReentrantTimesHeader, ok := message.XPropertyReentrantTimes(sr.options.status); ok {
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
		return sr.options.deadHandler.Handle(msg, checker.CheckStatusPassed)
	}
	return false
}
