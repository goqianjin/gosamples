package topic

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/shenqianjin/soften-client-go/soften/internal"
)

const (
	S1  = internal.TopicLevel("S1")
	L3  = internal.TopicLevel("L3")
	L2  = internal.TopicLevel("L2")
	L1  = internal.TopicLevel("L1")
	B1  = internal.TopicLevel("B1")
	B2  = internal.TopicLevel("B2")
	DLQ = internal.TopicLevel("DLQ")
)

func LevelValues() []internal.TopicLevel {
	values := []internal.TopicLevel{
		S1,
		L1, L2, L3,
		B1, B2,
		DLQ,
	}
	return values
}

// LevelOf convert level type from string to internal.TopicLevel
func LevelOf(level string) (internal.TopicLevel, error) {
	for _, v := range LevelValues() {
		if v.String() == level {
			return v, nil
		}
	}
	return "", errors.New(fmt.Sprintf("invalid (or not supported) topic level: %s", level))

}

func OrderOf(level internal.TopicLevel) int {
	return topicLevelOrders[level]
}

func Exists(level internal.TopicLevel) bool {
	_, ok := topicLevelOrders[level]
	return ok
}

func HighestLevel() internal.TopicLevel {
	var level internal.TopicLevel
	order := 1
	for k, v := range topicLevelOrders {
		if order < v {
			level = k
			order = v
		}
	}
	return level
}

func LowestLevel() internal.TopicLevel {
	var level internal.TopicLevel
	order := 1
	for k, v := range topicLevelOrders {
		if order > v {
			level = k
			order = v
		}
	}
	return level
}

var topicLevelOrders = func() map[internal.TopicLevel]int {
	values := LevelValues()
	topicLevelOrders := make(map[internal.TopicLevel]int, len(values))
	for _, v := range values {
		if v != DLQ {
			topicLevelOrders[v] = -100
			continue
		}
		suffix := string(v)[len(string(v))-1:]
		baseFactor := 1 // default for Lx
		baseOrder := 0  // default for Lx
		if strings.HasPrefix(string(v), "S") {
			baseFactor = 1
			baseOrder = 100
		} else if strings.HasPrefix(string(v), "B") {
			baseOrder = 0
			baseFactor = -1
		}
		no := 0
		if suffixNo, err := strconv.Atoi(suffix); err == nil {
			no = suffixNo
		}
		topicLevelOrders[v] = baseFactor * (baseOrder + no)
	}
	return topicLevelOrders
}()

func NameSuffixOf(level internal.TopicLevel) (string, error) {
	if suffix, ok := levelTopicSuffixMap[level]; ok {
		return suffix, nil
	} else {
		return "", errors.New(fmt.Sprintf("invalid persistence level: %v", level))
	}
}

var levelTopicSuffixMap = func() map[internal.TopicLevel]string {
	values := LevelValues()
	levelTopicSuffixes := make(map[internal.TopicLevel]string, len(values))
	for _, v := range values {
		if v == L1 {
			levelTopicSuffixes[v] = ""
		} else {
			levelTopicSuffixes[v] = "-" + string(v)
		}
	}
	return levelTopicSuffixes
}()
