package soam

type messageStatus string

const (
	MessageStatusReady    = messageStatus("Ready")
	MessageStatusPending  = messageStatus("Pending")
	MessageStatusBlocking = messageStatus("Blocking")
	MessageStatusRetrying = messageStatus("Retrying")
	MessageStatusDead     = messageStatus("Dead")
	MessageStatusDone     = messageStatus("Done")
	MessageStatusDiscard  = messageStatus("Discard")
	messageStatusNewReady = messageStatus("NewReady")
)

const (
	MessageStatusNewReadyByReRoute = messageStatus("NewReadyByReRoute")
	MessageStatusNewReadyByUpgrade = messageStatus("NewReadyByUpgrade")
	MessageStatusNewReadyByDegrade = messageStatus("NewReadyByDegrade")
)

var messageStatusMap = map[string]messageStatus{
	string(MessageStatusReady):    MessageStatusReady,
	string(MessageStatusPending):  MessageStatusPending,
	string(MessageStatusBlocking): MessageStatusBlocking,
	string(MessageStatusRetrying): MessageStatusRetrying,
	string(MessageStatusDead):     MessageStatusDead,
}

var statusTopicSuffixMap = map[messageStatus]string{
	MessageStatusReady:    "",
	MessageStatusPending:  "-" + string(MessageStatusPending),
	MessageStatusBlocking: "-" + string(MessageStatusBlocking),
	MessageStatusRetrying: "-" + string(MessageStatusRetrying),
	MessageStatusDead:     "-DLQ",
}
