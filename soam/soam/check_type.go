package soam

// ------ check type ------

type CheckType string

const (
	CheckTypePreDiscard  = CheckType("PreDiscard")
	CheckTypePrePending  = CheckType("PrePending")
	CheckTypePreBlocking = CheckType("PreBlocking")
	CheckTypePreRetrying = CheckType("PreRetrying")
	CheckTypePreDead     = CheckType("PreDead")
	CheckTypePreReRoute  = CheckType("PreReRoute")
	CheckTypePreUpgrade  = CheckType("PreUpgrade")
	CheckTypePreDegrade  = CheckType("PreDegrade")

	CheckTypePostDiscard  = CheckType("PostDiscard")
	CheckTypePostPending  = CheckType("PostPending")
	CheckTypePostBlocking = CheckType("PostBlocking")
	CheckTypePostRetrying = CheckType("PostRetrying")
	CheckTypePostDead     = CheckType("PostDead")
	CheckTypePostReRoute  = CheckType("PostReRoute")
	CheckTypePostUpgrade  = CheckType("PostUpgrade")
	CheckTypePostDegrade  = CheckType("PostDegrade")
)

var defaultPreCheckTypesInTurn = []CheckType{CheckTypePreDiscard, CheckTypePreDead,
	CheckTypePreReRoute, CheckTypePreUpgrade, CheckTypePreDegrade,
	CheckTypePreBlocking, CheckTypePrePending, CheckTypePreRetrying}

var DefaultPostCheckTypesInTurn = []CheckType{CheckTypePostDiscard, CheckTypePostDead,
	CheckTypePostReRoute, CheckTypePostUpgrade, CheckTypePostDegrade,
	CheckTypePreBlocking, CheckTypePrePending, CheckTypePostRetrying}

var messageStatusToPostCheckTypeMap = map[messageStatus]CheckType{
	MessageStatusDiscard:           CheckTypePostDiscard,
	MessageStatusDead:              CheckTypePostDead,
	MessageStatusNewReadyByReRoute: CheckTypePostReRoute,
	MessageStatusNewReadyByUpgrade: CheckTypePostUpgrade,
	MessageStatusNewReadyByDegrade: CheckTypePostDegrade,
	MessageStatusBlocking:          CheckTypePostBlocking,
	MessageStatusPending:           CheckTypePostPending,
	MessageStatusRetrying:          CheckTypePostRetrying,
}
