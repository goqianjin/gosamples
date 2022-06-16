package message

import (
	"errors"
	"fmt"
	"strings"

	"github.com/shenqianjin/soften-client-go/soften/internal"
)

const (
	StatusReady    = internal.MessageStatus("Ready")
	StatusPending  = internal.MessageStatus("Pending")
	StatusBlocking = internal.MessageStatus("Blocking")
	StatusRetrying = internal.MessageStatus("Retrying")
	StatusDead     = internal.MessageStatus("Dead")
	StatusDone     = internal.MessageStatus("Done")
	StatusDiscard  = internal.MessageStatus("Discard")
	//statusNewReady = internal.MessageStatus("NewReady")
)

func StatusOf(status string) (internal.MessageStatus, error) {
	for _, v := range StatusValues() {
		if v.String() == status {
			return v, nil
		}
	}
	return "", errors.New(fmt.Sprintf("invalid message status: %s", status))
}

func StatusValues() []internal.MessageStatus {
	values := []internal.MessageStatus{
		StatusReady,
		StatusBlocking, StatusPending, StatusRetrying,
		StatusDead, StatusDone, StatusDiscard,
	}
	return values
}

func TopicSuffixOf(status internal.MessageStatus) (string, error) {
	if suffix, ok := statusTopicSuffixMap[status]; ok {
		return suffix, nil
	} else {
		return "", errors.New(fmt.Sprintf("invalid persistence status: %v", status))
	}
}

var statusTopicSuffixMap = map[internal.MessageStatus]string{
	StatusReady:    "",
	StatusPending:  "-" + strings.ToUpper(string(StatusPending)),
	StatusBlocking: "-" + strings.ToUpper(string(StatusBlocking)),
	StatusRetrying: "-" + strings.ToUpper(string(StatusRetrying)),
	StatusDead:     "-DLQ",
}
