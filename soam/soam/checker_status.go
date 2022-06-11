package soam

import "github.com/apache/pulsar-client-go/pulsar"

// ------ status checker ------

type preStatusChecker func(pulsar.Message) (passed bool)

type postStatusChecker func(pulsar.Message, error) (passed bool)

var nilPreStatusChecker = func(pulsar.Message) (passed bool) {
	return false
}

var nilPostStatusChecker = func(pulsar.Message, error) (passed bool) {
	return false
}

// ------ re-router checker ------

type preRerouteChecker func(pulsar.Message) string

type postRerouteChecker func(pulsar.Message, error) string

var nilPreRerouteChecker = func(pulsar.Message) string {
	return ""
}

var nilPostRerouteChecker = func(pulsar.Message, error) string {
	return ""
}

// ------ action checkers ------

func PreDiscardChecker(checker preStatusChecker) checkpoint {
	return checkpoint{checkType: CheckTypePreDiscard, preStatusChecker: checker}
}

func PostDiscardChecker(checker postStatusChecker) checkpoint {
	return checkpoint{checkType: CheckTypePostDiscard, postStatusChecker: checker}
}

func PrePendingChecker(checker preStatusChecker) checkpoint {
	return checkpoint{checkType: CheckTypePrePending, preStatusChecker: checker}
}

func PostPendingChecker(checker postStatusChecker) checkpoint {
	return checkpoint{checkType: CheckTypePostPending, postStatusChecker: checker}
}

func PreBlockingChecker(checker preStatusChecker) checkpoint {
	return checkpoint{checkType: CheckTypePreBlocking, preStatusChecker: checker}
}

func PostBlockingChecker(checker postStatusChecker) checkpoint {
	return checkpoint{checkType: CheckTypePostBlocking, postStatusChecker: checker}
}

func PreRetryingChecker(checker preStatusChecker) checkpoint {
	return checkpoint{checkType: CheckTypePreRetrying, preStatusChecker: checker}
}

func PostRetryingChecker(checker postStatusChecker) checkpoint {
	return checkpoint{checkType: CheckTypePostRetrying, postStatusChecker: checker}
}

func PreDeadChecker(checker preStatusChecker) checkpoint {
	return checkpoint{checkType: CheckTypePreDead, preStatusChecker: checker}
}

func PostDeadChecker(checker postStatusChecker) checkpoint {
	return checkpoint{checkType: CheckTypePostDead, postStatusChecker: checker}
}

func PreUpgradeChecker(checker preStatusChecker) checkpoint {
	return checkpoint{checkType: CheckTypePreUpgrade, preStatusChecker: checker}
}

func PostUpgradeChecker(checker postStatusChecker) checkpoint {
	return checkpoint{checkType: CheckTypePostUpgrade, postStatusChecker: checker}
}

func PreDegradeChecker(checker preStatusChecker) checkpoint {
	return checkpoint{checkType: CheckTypePreDegrade, preStatusChecker: checker}
}

func PostDegradeChecker(checker postStatusChecker) checkpoint {
	return checkpoint{checkType: CheckTypePostDegrade, postStatusChecker: checker}
}

func PreRerouteChecker(checker preRerouteChecker) checkpoint {
	return checkpoint{checkType: CheckTypePreReroute, preReRouterChecker: checker}
}

func PostRerouteChecker(checker postRerouteChecker) checkpoint {
	return checkpoint{checkType: CheckTypePostReroute, postReRouterChecker: checker}
}
