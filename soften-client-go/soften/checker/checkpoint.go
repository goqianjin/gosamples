package checker

import (
	"github.com/shenqianjin/soften-client-go/soften/internal"
)

// ------ consume checkpoint ------

type Checkpoint struct {
	CheckType internal.CheckType
	Before    BeforeCheckFunc
	After     AfterCheckFunc
}

var NilCheckpoint = &Checkpoint{}

// ------ status checkers ------

func PreDiscardChecker(checker BeforeCheckFunc) Checkpoint {
	return Checkpoint{CheckType: CheckTypePreDiscard, Before: checker}
}

func PostDiscardChecker(checker AfterCheckFunc) Checkpoint {
	return Checkpoint{CheckType: CheckTypePostDiscard, After: checker}
}

func PrePendingChecker(checker BeforeCheckFunc) Checkpoint {
	return Checkpoint{CheckType: CheckTypePrePending, Before: checker}
}

func PostPendingChecker(checker AfterCheckFunc) Checkpoint {
	return Checkpoint{CheckType: CheckTypePostPending, After: checker}
}

func PreBlockingChecker(checker BeforeCheckFunc) Checkpoint {
	return Checkpoint{CheckType: CheckTypePreBlocking, Before: checker}
}

func PostBlockingChecker(checker AfterCheckFunc) Checkpoint {
	return Checkpoint{CheckType: CheckTypePostBlocking, After: checker}
}

func PreRetryingChecker(checker BeforeCheckFunc) Checkpoint {
	return Checkpoint{CheckType: CheckTypePreRetrying, Before: checker}
}

func PostRetryingChecker(checker AfterCheckFunc) Checkpoint {
	return Checkpoint{CheckType: CheckTypePostRetrying, After: checker}
}

func PreDeadChecker(checker BeforeCheckFunc) Checkpoint {
	return Checkpoint{CheckType: CheckTypePreDead, Before: checker}
}

func PostDeadChecker(checker AfterCheckFunc) Checkpoint {
	return Checkpoint{CheckType: CheckTypePostDead, After: checker}
}

func PreUpgradeChecker(checker BeforeCheckFunc) Checkpoint {
	return Checkpoint{CheckType: CheckTypePreUpgrade, Before: checker}
}

func PostUpgradeChecker(checker AfterCheckFunc) Checkpoint {
	return Checkpoint{CheckType: CheckTypePostUpgrade, After: checker}
}

func PreDegradeChecker(checker BeforeCheckFunc) Checkpoint {
	return Checkpoint{CheckType: CheckTypePreDegrade, Before: checker}
}

func PostDegradeChecker(checker AfterCheckFunc) Checkpoint {
	return Checkpoint{CheckType: CheckTypePostDegrade, After: checker}
}

// ------ reroute checkers ------

func PreRerouteChecker(checker BeforeCheckFunc) Checkpoint {
	return Checkpoint{CheckType: CheckTypePreReroute, Before: checker}
}

func PostRerouteChecker(checker AfterCheckFunc) Checkpoint {
	return Checkpoint{CheckType: CheckTypePostReroute, After: checker}
}

// ------ produce checkpoint ------

type ProduceCheckpoint struct {
	CheckType internal.CheckType
	CheckFunc ProduceCheckFunc
}

var NilProduceCheckpoint = &ProduceCheckpoint{}

// ------ route checker ------

func ProduceDiscardChecker(checker ProduceCheckFunc) ProduceCheckpoint {
	return ProduceCheckpoint{CheckType: ProduceCheckTypeDiscard, CheckFunc: checker}
}

func ProducePendingChecker(checker ProduceCheckFunc) ProduceCheckpoint {
	return ProduceCheckpoint{CheckType: ProduceCheckTypePending, CheckFunc: checker}
}

func ProduceBlockingChecker(checker ProduceCheckFunc) ProduceCheckpoint {
	return ProduceCheckpoint{CheckType: ProduceCheckTypeBlocking, CheckFunc: checker}
}

func ProduceRetryingChecker(checker ProduceCheckFunc) ProduceCheckpoint {
	return ProduceCheckpoint{CheckType: ProduceCheckTypeRetrying, CheckFunc: checker}
}

func ProduceDeadChecker(checker ProduceCheckFunc) ProduceCheckpoint {
	return ProduceCheckpoint{CheckType: ProduceCheckTypeDead, CheckFunc: checker}
}

func ProduceUpgradeChecker(checker ProduceCheckFunc) ProduceCheckpoint {
	return ProduceCheckpoint{CheckType: ProduceCheckTypeUpgrade, CheckFunc: checker}
}

func ProduceDegradeChecker(checker ProduceCheckFunc) ProduceCheckpoint {
	return ProduceCheckpoint{CheckType: ProduceCheckTypeDegrade, CheckFunc: checker}
}

func ProduceRouteChecker(checker ProduceCheckFunc) ProduceCheckpoint {
	return ProduceCheckpoint{CheckType: ProduceCheckTypeRoute, CheckFunc: checker}
}
