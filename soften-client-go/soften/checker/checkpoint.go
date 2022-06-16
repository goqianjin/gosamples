package checker

import (
	"github.com/shenqianjin/soften-client-go/soften/internal"
)

// ------ status checkers ------

func PreDiscardChecker(checker internal.PreStatusChecker) internal.Checkpoint {
	return internal.Checkpoint{CheckType: CheckTypePreDiscard, PreStatusChecker: checker}
}

func PostDiscardChecker(checker internal.PostStatusChecker) internal.Checkpoint {
	return internal.Checkpoint{CheckType: CheckTypePostDiscard, PostStatusChecker: checker}
}

func PrePendingChecker(checker internal.PreStatusChecker) internal.Checkpoint {
	return internal.Checkpoint{CheckType: CheckTypePrePending, PreStatusChecker: checker}
}

func PostPendingChecker(checker internal.PostStatusChecker) internal.Checkpoint {
	return internal.Checkpoint{CheckType: CheckTypePostPending, PostStatusChecker: checker}
}

func PreBlockingChecker(checker internal.PreStatusChecker) internal.Checkpoint {
	return internal.Checkpoint{CheckType: CheckTypePreBlocking, PreStatusChecker: checker}
}

func PostBlockingChecker(checker internal.PostStatusChecker) internal.Checkpoint {
	return internal.Checkpoint{CheckType: CheckTypePostBlocking, PostStatusChecker: checker}
}

func PreRetryingChecker(checker internal.PreStatusChecker) internal.Checkpoint {
	return internal.Checkpoint{CheckType: CheckTypePreRetrying, PreStatusChecker: checker}
}

func PostRetryingChecker(checker internal.PostStatusChecker) internal.Checkpoint {
	return internal.Checkpoint{CheckType: CheckTypePostRetrying, PostStatusChecker: checker}
}

func PreDeadChecker(checker internal.PreStatusChecker) internal.Checkpoint {
	return internal.Checkpoint{CheckType: CheckTypePreDead, PreStatusChecker: checker}
}

func PostDeadChecker(checker internal.PostStatusChecker) internal.Checkpoint {
	return internal.Checkpoint{CheckType: CheckTypePostDead, PostStatusChecker: checker}
}

func PreUpgradeChecker(checker internal.PreStatusChecker) internal.Checkpoint {
	return internal.Checkpoint{CheckType: CheckTypePreUpgrade, PreStatusChecker: checker}
}

func PostUpgradeChecker(checker internal.PostStatusChecker) internal.Checkpoint {
	return internal.Checkpoint{CheckType: CheckTypePostUpgrade, PostStatusChecker: checker}
}

func PreDegradeChecker(checker internal.PreStatusChecker) internal.Checkpoint {
	return internal.Checkpoint{CheckType: CheckTypePreDegrade, PreStatusChecker: checker}
}

func PostDegradeChecker(checker internal.PostStatusChecker) internal.Checkpoint {
	return internal.Checkpoint{CheckType: CheckTypePostDegrade, PostStatusChecker: checker}
}

// ------ reroute checkers ------

func PreRerouteChecker(checker internal.PreRerouteChecker) internal.Checkpoint {
	return internal.Checkpoint{CheckType: CheckTypePreReroute, PreRerouteChecker: checker}
}

func PostRerouteChecker(checker internal.PostRerouteChecker) internal.Checkpoint {
	return internal.Checkpoint{CheckType: CheckTypePostReroute, PostRerouteChecker: checker}
}
