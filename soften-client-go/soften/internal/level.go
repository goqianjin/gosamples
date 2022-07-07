package internal

import (
	"strings"
)

// ------ topic level ------

type TopicLevel string

func (e TopicLevel) String() string {
	return string(e)
}

// ------ topic level parser ------

var TopicLevelParser = topicLevelParser{}

type topicLevelParser struct {
}

func (p topicLevelParser) FormatList(levels []TopicLevel) string {
	if len(levels) <= 0 {
		return ""
	}
	ls := make([]string, len(levels))
	for i := 0; i < len(levels)-1; i++ {
		ls[i] = levels[i].String()
	}
	return strings.Join(ls, ", ")
}

// ------ Balance Strategy Symbol ------

type BalanceStrategy string

func (e BalanceStrategy) String() string {
	return string(e)
}
