package soften

import (
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/shenqianjin/soften-client-go/soften/checker"
	"github.com/shenqianjin/soften-client-go/soften/internal"
	"github.com/shenqianjin/soften-client-go/soften/message"
)

// ------ consumer biz handle interfaces ------

// Handler is the regular processing flow, and it is recommended.
// the message will be acknowledged when the return is true; If it returns false,
// the message will be unacknowledged to the main partition, then route to retrying
// parting if the retrying module is enabled. finally, it goto dead letter partition
// when all retrying times exceed the maximum.
type Handler func(pulsar.Message) (success bool, err error)

// PremiumHandler allows the result contains Done, Retrying, Dead, Pending, Blocking, Degrade and Upgrade statuses.
// different status will deliver current message to the corresponding destination.
// Please note the process will be regressed to regular module when the returned
// status is not enough to do its flow, e.g. HandleStatusPending is returned
// however the pending flow is not enabled in the multiStatusConsumeFacade configuration.
type PremiumHandler func(pulsar.Message) HandleStatus

// ------ handle status ------

type HandleStatus interface {
	getGotoAction() internal.MessageGoto
	getCheckTypes() []internal.CheckType
	getErr() error
}

var (
	HandleStatusOk      = handleStatus{gotoAction: message.GotoDone}                // handle message successfully
	HandleStatusFail    = handleStatus{checkTypes: checker.DefaultPostCheckTypes()} // handle message failure
	HandleStatusDefault = handleStatus{}                                            // default handle message
)

// ------ handleStatus impl ------

type handleStatus struct {
	gotoAction internal.MessageGoto // 状态转移至新状态: 默认为空; 与当前状态一致时，参数无效。优先级比 checkTypes 高。
	checkTypes []internal.CheckType // 事后检查类型列表: 默认 DefaultPostCheckTypesInTurn; 指定时，使用指定类型列表。优先级比 gotoAction 低。
	err        error                // 后置校验器如果需要依赖处理错误，通过该参数传递。框架在处理的过程不会更新err内容，client自己在传递的过程有更新除外。
}

func (h handleStatus) getGotoAction() internal.MessageGoto {
	return h.gotoAction
}

func (h handleStatus) getCheckTypes() []internal.CheckType {
	return h.checkTypes
}

func (h handleStatus) getErr() error {
	return h.err
}

func (h handleStatus) GotoAction(action internal.MessageGoto) handleStatus {
	h.gotoAction = action
	return h
}

func (h handleStatus) CheckTypes(types []internal.CheckType) handleStatus {
	h.checkTypes = types
	return h
}

func (h handleStatus) Err(err error) handleStatus {
	h.err = err
	return h
}
