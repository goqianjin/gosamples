package soam

import (
	"fmt"
	"strconv"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

var MessageParser = &messageParser{}

type messageParser struct {
}

func (p *messageParser) GetCurrentStatus(msg pulsar.ConsumerMessage) messageStatus {
	properties := msg.Message.Properties()
	if status, ok := properties[SysPropertyXCurrentMessageStatus]; ok {
		if messageStatus, ok2 := messageStatusMap[status]; ok2 {
			return messageStatus
		}
	}
	return MessageStatusReady
}

func (p *messageParser) GetPreviousStatus(msg pulsar.ConsumerMessage) messageStatus {
	properties := msg.Message.Properties()
	if status, ok := properties[SysPropertyPreviousMessageStatus]; ok {
		if messageStatus, ok2 := messageStatusMap[status]; ok2 {
			return messageStatus
		}
	}
	return ""
}

func (p *messageParser) GetXReconsumeTimes(msg pulsar.ConsumerMessage) int {
	properties := msg.Message.Properties()
	if timesStr, ok := properties[SysPropertyXReconsumeTimes]; ok {
		if times, err := strconv.Atoi(timesStr); err == nil {
			return times
		}
	}
	return 0
}

func (p *messageParser) GetReentrantStartRedeliveryCount(msg pulsar.ConsumerMessage) uint32 {
	properties := msg.Message.Properties()
	if timesStr, ok := properties[SysPropertyReentrantStartRedeliveryCount]; ok {
		if times, err := strconv.ParseUint(timesStr, 10, 32); err == nil {
			return uint32(times)
		}
	}
	return 0
}

func (p *messageParser) GetStatusReconsumeTimes(status messageStatus, msg pulsar.ConsumerMessage) int {
	statusConsumeTimesHeader, ok := statusConsumeTimesMap[status]
	if !ok {
		panic("invalid status for statusConsumeTimes")
	}
	properties := msg.Message.Properties()
	if timesStr, ok := properties[statusConsumeTimesHeader]; ok {
		if times, err := strconv.Atoi(timesStr); err == nil {
			return times
		}
	}
	return 0
}

func (p *messageParser) GetStatusReentrantTimes(status messageStatus, msg pulsar.ConsumerMessage) int {
	reentrantTimesHeader, ok := statusReentrantTimesMap[status]
	if !ok {
		panic("invalid status for statusReentrantTimes")
	}
	properties := msg.Message.Properties()
	if timesStr, ok := properties[reentrantTimesHeader]; ok {
		if times, err := strconv.Atoi(timesStr); err == nil {
			return times
		}
	}
	return 0
}

func (p *messageParser) GetReconsumeTime(msg pulsar.ConsumerMessage) time.Time {
	properties := msg.Message.Properties()
	if timeStr, ok := properties[SysPropertyXReconsumeTime]; ok {
		if time, err := time.Parse(RFC3339TimeInSecondPattern, timeStr); err == nil {
			return time
		}
	}
	return time.Time{}
}

func (p *messageParser) GetReentrantTime(msg pulsar.ConsumerMessage) time.Time {
	properties := msg.Message.Properties()
	if timeStr, ok := properties[SysPropertyXReentrantTime]; ok {
		if time, err := time.Parse(RFC3339TimeInSecondPattern, timeStr); err == nil {
			return time
		}
	}
	return time.Time{}
}

func (p *messageParser) GetMessageId(msg pulsar.ConsumerMessage) string {
	id := msg.Message.ID()
	return fmt.Sprintf("%d:%d:%d", id.LedgerID(), id.EntryID(), id.PartitionIdx())
}
