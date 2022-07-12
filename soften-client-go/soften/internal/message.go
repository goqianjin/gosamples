package internal

import (
	"strings"
)

// ------ message status ------

type MessageStatus string

const (
	DefaultMessageStatusReady = "Ready"
)

func (status MessageStatus) String() string {
	return string(status)
}

func (status MessageStatus) TopicSuffix() string {
	if status == DefaultMessageStatusReady {
		return ""
	} else {
		return "-" + strings.ToUpper(string(status))
	}
}

// ------ message goto action ------

type MessageGoto string

func (e MessageGoto) String() string {
	return string(e)
}
