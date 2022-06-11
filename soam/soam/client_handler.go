package soam

import "github.com/apache/pulsar-client-go/pulsar"

// Handler is the regular processing flow, and it is recommended.
// the message will be acknowledged when the return is true; If it returns false,
// the message will be unacknowledged to the main partition, then route to retrying
// parting if the retrying module is enabled. finally, it goto dead letter partition
// when all retrying times exceed the maximum.
type Handler func(pulsar.Message) (success bool, err error)

// HandlerInPremium allows the result contains Done, Retrying, Dead, Pending, Blocking, Degrade and Upgrade statuses.
// different status will deliver current message to the corresponding destination.
// Please note the process will be regressed to regular module when the returned
// status is not enough to do its flow, e.g. HandleStatusPending is returned
// however the pending flow is not enabled in the consumer configuration.
type HandlerInPremium func(pulsar.Message) HandleResult

var (
	HandledOk   = handleResult{handledStatus: MessageStatusDone}
	HandledFail = handleResult{postCheckTypesInTurn: DefaultPostCheckTypesInTurn}
)

type HandleResult interface {
	getHandledStatus() messageStatus
	getPostCheckTypesInTurn() []CheckType
	getErr() error
}

type handleResult struct {
	handledStatus        messageStatus // 状态转移至新状态: 默认为空; 与当前状态一致时，参数无效。优先级比 postCheckTypes 高。
	postCheckTypesInTurn []CheckType   // 事后检查类型列表: 默认 DefaultPostCheckTypesInTurn; 指定时，使用指定类型列表。优先级比 handledStatus 低。
	err                  error         // 后置校验器如果需要依赖处理错误，通过该参数传递。框架在处理的过程不会更新err内容，client自己在传递的过程有更新除外。
}

func (h handleResult) getHandledStatus() messageStatus {
	return h.handledStatus
}

func (h handleResult) getPostCheckTypesInTurn() []CheckType {
	return h.postCheckTypesInTurn
}

func (h handleResult) getErr() error {
	return h.err
}

func (h handleResult) TransferTo(status messageStatus) handleResult {
	h.handledStatus = status
	return h
}

func (h handleResult) PostCheckTypesInTurn(types []CheckType) handleResult {
	h.postCheckTypesInTurn = types
	return h
}

func (h handleResult) Err(err error) handleResult {
	h.err = err
	return h
}
