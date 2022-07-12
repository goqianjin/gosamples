package internal

import (
	"fmt"
	"strconv"
	"strings"
)

// ------ topic level ------

type TopicLevel string

var (
	DefaultGroundTopicLevelL1 = TopicLevel("L1")
	DefaultDeadTopicLevelDLQ  = TopicLevel("DLQ")
)

func (lvl TopicLevel) String() string {
	return string(lvl)
}

func (lvl TopicLevel) OrderOf() int {
	if lvl == DefaultDeadTopicLevelDLQ {
		return -100
	}
	suffix := string(lvl)[len(string(lvl))-1:]
	baseFactor := 1 // default for Lx
	baseOrder := 0  // default for Lx
	if strings.HasPrefix(string(lvl), "S") {
		baseFactor = 1
		baseOrder = 100
	} else if strings.HasPrefix(string(lvl), "B") {
		baseFactor = -1
		baseOrder = 0
	} else if strings.HasPrefix(string(lvl), "DLQ") {
		baseFactor = -1
		baseOrder = 100
	} else {
		panic(fmt.Sprintf("invalid topic level: %v", lvl))
	}
	no := 0
	if suffixNo, err := strconv.Atoi(suffix); err == nil {
		no = suffixNo
	}
	return baseFactor * (baseOrder + no)
}

func (lvl TopicLevel) TopicSuffix() string {
	if lvl == DefaultGroundTopicLevelL1 {
		return ""
	} else {
		return "-" + string(lvl)
	}
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
	for i := 0; i < len(levels); i++ {
		ls[i] = levels[i].String()
	}
	return strings.Join(ls, ", ")
}

// ------ Balance Strategy Symbol ------

type BalanceStrategy string

func (e BalanceStrategy) String() string {
	return string(e)
}
