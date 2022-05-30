package soam

type BackoffPolicy interface {
	Next(redeliveryCount int) uint
}

type StatusBackoffPolicy interface {
	Next(statusReconsumeTimes int) uint
}
