package soam

const (
	SysPropertyPreviousMessageStatus         = "X-Previous-Status"                  // 前一个消息状态
	SysPropertyXCurrentMessageStatus         = "X-Current-Status"                   // 当前消息的状态
	SysPropertyXReconsumeTimes               = "X-Reconsume-Times"                  // 总重试消费次数
	SysPropertyXReconsumeTime                = "X-Reconsume-Time"                   // 消费时间
	SysPropertyXReentrantTime                = "X-Reentrant-Time"                   // 重入时间
	SysPropertyReentrantStartRedeliveryCount = "X-Reentrant-Start-Redelivery-Count" // 当前状态开始的消费次数

	SysPropertyPendingReconsumeTimes  = "X-Pending-Reconsume-Times" // 状态消费次数
	SysPropertyPendingReentrantTimes  = "X-Pending-Reentrant-Times" // 状态重入次数
	SysPropertyBlockingReconsumeTimes = "X-Blocking-Reconsume-Times"
	SysPropertyBlockingReentrantTimes = "X-Blocking-Reentrant-Times"
	SysPropertyRetryingReconsumeTimes = "X-Retrying-Reconsume-Times"
	SysPropertyRetryingReentrantTimes = "X-Retrying-Reentrant-Times"
	SysPropertyReadyReconsumeTimes    = "X-Ready-Reconsume-Times"
	SysPropertyReadyReentrantTimes    = "X-Ready-Reentrant-Times"
	SysPropertyDeadReconsumeTimes     = "X-Dead-Reconsume-Times"
	SysPropertyDeadReentrantTimes     = "X-Dead-Reentrant-Times"
)

var (
	statusConsumeTimesMap = map[messageStatus]string{
		MessageStatusPending:  SysPropertyPendingReconsumeTimes,
		MessageStatusBlocking: SysPropertyBlockingReconsumeTimes,
		MessageStatusRetrying: SysPropertyRetryingReconsumeTimes,
		MessageStatusReady:    SysPropertyReadyReconsumeTimes,
		MessageStatusDead:     SysPropertyDeadReconsumeTimes,
	}
	statusReentrantTimesMap = map[messageStatus]string{
		MessageStatusPending:  SysPropertyPendingReentrantTimes,
		MessageStatusBlocking: SysPropertyBlockingReentrantTimes,
		MessageStatusRetrying: SysPropertyRetryingReentrantTimes,
		MessageStatusReady:    SysPropertyReadyReentrantTimes,
		MessageStatusDead:     SysPropertyDeadReentrantTimes,
	}
)
