package soam

const (
	DlqTopicSuffix    = "-DLQ"
	RetryTopicSuffix  = "-RETRY"
	MaxReconsumeTimes = 16

	SysPropertyDelayTime       = "DELAY_TIME"
	SysPropertyRealTopic       = "REAL_TOPIC"
	SysPropertyRetryTopic      = "RETRY_TOPIC"
	SysPropertyReconsumeTimes  = "RECONSUMETIMES"
	SysPropertyOriginMessageID = "ORIGIN_MESSAGE_IDY_TIME"

	SysPropertyPendingNackBackoffTimes  = "X-Pending-Nack-Backoff-Times"
	SysPropertyPendingBackoffTimes      = "X-Pending-Reentrant-Backoff-Times"
	SysPropertyBlockingNackBackoffTimes = "X-Blocking-Nack-Backoff-Times"
	SysPropertyBlockingBackoffTimes     = "X-Blocking-Backoff-Times"
	SysPropertyRetryingNackBackoffTimes = "X-Retrying-Nack-Backoff-Times"
	SysPropertyRetryingBackoffTimes     = "X-Retrying-Backoff-Times"
)
