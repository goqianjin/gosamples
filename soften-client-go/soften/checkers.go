package soften

import (
	"github.com/apache/pulsar-client-go/pulsar/log"
	"github.com/shenqianjin/soften-client-go/soften/checker"
	"github.com/shenqianjin/soften-client-go/soften/internal"
)

type consumeCheckers struct {
	log log.Logger

	PreBlockingChecker  internal.PreStatusChecker  //
	PostBlockingChecker internal.PostStatusChecker //
	PrePendingChecker   internal.PreStatusChecker  //
	PostPendingChecker  internal.PostStatusChecker //
	PreRetryingChecker  internal.PreStatusChecker  //
	PostRetryingChecker internal.PostStatusChecker //
	PreDeadChecker      internal.PreStatusChecker  //
	PostDeadChecker     internal.PostStatusChecker //
	PreDiscardChecker   internal.PreStatusChecker  //
	PostDiscardChecker  internal.PostStatusChecker //
	PreUpgradeChecker   internal.PreStatusChecker  //
	PostUpgradeChecker  internal.PostStatusChecker //
	PreDegradeChecker   internal.PreStatusChecker  //
	PostDegradeChecker  internal.PostStatusChecker //

	PreRerouteChecker  internal.PreRerouteChecker  //
	PostReRouteChecker internal.PostRerouteChecker //
}

// newConsumeCheckers create consume checkers.
// xxxEnables of config and checkpoints in subscribe parameters is used in this construction.
func newConsumeCheckers(logger log.Logger, enables *internal.StatusEnables, checkpointMap map[internal.CheckType]*internal.Checkpoint) (*consumeCheckers, error) {
	checkers := &consumeCheckers{
		log: logger,
	}
	if enables.RerouteEnable {
		checkers.PreRerouteChecker = checkers.parsePreRerouteChecker(checker.CheckTypePreReroute, checkpointMap)
		checkers.PostReRouteChecker = checkers.parsePostRerouteChecker(checker.CheckTypePostReroute, checkpointMap)
	}
	if enables.PendingEnable {
		checkers.PrePendingChecker = checkers.parsePreStatusChecker(checker.CheckTypePrePending, checkpointMap)
		checkers.PostPendingChecker = checkers.parsePostStatusChecker(checker.CheckTypePostPending, checkpointMap)
	}
	if enables.BlockingEnable {
		checkers.PreBlockingChecker = checkers.parsePreStatusChecker(checker.CheckTypePreBlocking, checkpointMap)
		checkers.PostBlockingChecker = checkers.parsePostStatusChecker(checker.CheckTypePostBlocking, checkpointMap)
	}
	if enables.RetryingEnable {
		checkers.PreRetryingChecker = checkers.parsePreStatusChecker(checker.CheckTypePreRetrying, checkpointMap)
		checkers.PostRetryingChecker = checkers.parsePostStatusChecker(checker.CheckTypePostRetrying, checkpointMap)
	}
	if enables.DeadEnable {
		checkers.PreDeadChecker = checkers.parsePreStatusChecker(checker.CheckTypePreDead, checkpointMap)
		checkers.PostDeadChecker = checkers.parsePostStatusChecker(checker.CheckTypePostDead, checkpointMap)
	}
	if enables.DiscardEnable {
		checkers.PreDiscardChecker = checkers.parsePreStatusChecker(checker.CheckTypePreDiscard, checkpointMap)
		checkers.PostDiscardChecker = checkers.parsePostStatusChecker(checker.CheckTypePostDiscard, checkpointMap)
	}
	if enables.UpgradeEnable {
		checkers.PreUpgradeChecker = checkers.parsePreStatusChecker(checker.CheckTypePreUpgrade, checkpointMap)
		checkers.PostUpgradeChecker = checkers.parsePostStatusChecker(checker.CheckTypePostUpgrade, checkpointMap)
	}
	if enables.DegradeEnable {
		checkers.PreDegradeChecker = checkers.parsePreStatusChecker(checker.CheckTypePreDegrade, checkpointMap)
		checkers.PostDegradeChecker = checkers.parsePostStatusChecker(checker.CheckTypePostDegrade, checkpointMap)
	}
	return checkers, nil
}

func (ch *consumeCheckers) parsePreStatusChecker(checkType internal.CheckType, checkpointMap map[internal.CheckType]*internal.Checkpoint) internal.PreStatusChecker {
	if ckp, ok := checkpointMap[checkType]; ok {
		if ckp.PreStatusChecker != nil {
			return ckp.PreStatusChecker
		}
	}
	return internal.NilPreStatusChecker
}

func (ch *consumeCheckers) parsePostStatusChecker(checkType internal.CheckType, checkpointMap map[internal.CheckType]*internal.Checkpoint) internal.PostStatusChecker {
	if ckp, ok := checkpointMap[checkType]; ok {
		if ckp.PostStatusChecker != nil {
			return ckp.PostStatusChecker
		}
	}
	return internal.NilPostStatusChecker
}

func (ch *consumeCheckers) parsePreRerouteChecker(checkType internal.CheckType, checkpointMap map[internal.CheckType]*internal.Checkpoint) internal.PreRerouteChecker {
	if ckp, ok := checkpointMap[checkType]; ok {
		if ckp.PreRerouteChecker != nil {
			return ckp.PreRerouteChecker
		}
	}
	return internal.NilPreRerouteChecker
}

func (ch *consumeCheckers) parsePostRerouteChecker(checkType internal.CheckType, checkpointMap map[internal.CheckType]*internal.Checkpoint) internal.PostRerouteChecker {
	if ckp, ok := checkpointMap[checkType]; ok {
		if ckp.PostRerouteChecker != nil {
			return ckp.PostRerouteChecker
		}
	}
	return internal.NilPostRerouteChecker
}
