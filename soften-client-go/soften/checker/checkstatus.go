package checker

// ------ check status interface ------

type CheckStatus interface {
	IsPassed() bool
	GetHandledDefer() func()

	// extra for reroute checker

	GetRerouteTopic() string
}

// ------ check status implementation ------

type checkStatus struct {
	passed       bool
	handledDefer func()
	//CheckTo      MessageGoto
	rerouteTopic string
}

func (s *checkStatus) WithPassed(passed bool) *checkStatus {
	s.passed = passed
	return s
}

func (s *checkStatus) WithHandledDefer(handledDefer func()) *checkStatus {
	s.handledDefer = handledDefer
	return s
}

func (s *checkStatus) WithRerouteTopic(topic string) *checkStatus {
	s.rerouteTopic = topic
	return s
}

func (s *checkStatus) IsPassed() bool {
	return s.passed
}

func (s *checkStatus) GetHandledDefer() func() {
	return s.handledDefer
}

func (s *checkStatus) GetRerouteTopic() string {
	return s.rerouteTopic
}

// ------ check status enums ------

var (
	CheckStatusPassed  = &checkStatus{passed: true}
	CheckStatusFailed  = &checkStatus{passed: false}
	CheckStatusDefault = &checkStatus{} // default is equal to CheckStatusFailed
)
