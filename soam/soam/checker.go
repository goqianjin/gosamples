package soam

import "github.com/apache/pulsar-client-go/pulsar"

// ------ status checker ------

type preStatusChecker func(pulsar.ConsumerMessage) (passed bool)

type postStatusChecker func(pulsar.ConsumerMessage, error) (passed bool)

var nilPreStatusChecker = func(pulsar.ConsumerMessage) (passed bool) {
	return false
}

var nilPostStatusChecker = func(pulsar.ConsumerMessage, error) (passed bool) {
	return false
}

func PrePendingChecker(checker preStatusChecker) checkpoint {
	return checkpoint{checkType: CheckTypePrePending, preStatusChecker: checker}
}

func PreBlockingChecker(checker preStatusChecker) checkpoint {
	return checkpoint{checkType: CheckTypePreBlocking, preStatusChecker: checker}
}

func PreRetryingChecker(checker preStatusChecker) checkpoint {
	return checkpoint{checkType: CheckTypePreRetrying, preStatusChecker: checker}
}

func PreDeadChecker(checker preStatusChecker) checkpoint {
	return checkpoint{checkType: CheckTypePreDead, preStatusChecker: checker}
}

func PostPendingChecker(checker postStatusChecker) checkpoint {
	return checkpoint{checkType: CheckTypePostPending, postStatusChecker: checker}
}

func PostBlockingChecker(checker postStatusChecker) checkpoint {
	return checkpoint{checkType: CheckTypePostBlocking, postStatusChecker: checker}
}

func PostRetryingChecker(checker postStatusChecker) checkpoint {
	return checkpoint{checkType: CheckTypePostRetrying, postStatusChecker: checker}
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

// ------ re-router checker ------

type preReRouterChecker func(pulsar.ConsumerMessage) string

type postReRouterChecker func(pulsar.ConsumerMessage, error) string

var nilPreReRouterChecker = func(pulsar.ConsumerMessage) string {
	return ""
}

var nilPostReRouterChecker = func(pulsar.ConsumerMessage, error) string {
	return ""
}

func PreReRouterChecker(checker preReRouterChecker) checkpoint {
	return checkpoint{checkType: CheckTypePreReRoute, preReRouterChecker: checker}
}

func PostReRouterChecker(checker postReRouterChecker) checkpoint {
	return checkpoint{checkType: CheckTypePostReRoute, postReRouterChecker: checker}
}
