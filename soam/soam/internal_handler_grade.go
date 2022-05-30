package soam

import "github.com/apache/pulsar-client-go/pulsar"

type gradeRouter struct {
	*router
	topicLevels []TopicLevel
}

func (r *gradeRouter) Route(message pulsar.ConsumerMessage) bool {
	return false
}

func newGradeReRouter(client pulsar.Client, topicLevel TopicLevel) (*gradeRouter, error) {
	return nil, nil
}
