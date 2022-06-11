package soam

// ------ check type ------

type CheckType string

const (
	CheckTypePreDiscard  = CheckType("PreDiscard")
	CheckTypePrePending  = CheckType("PrePending")
	CheckTypePreBlocking = CheckType("PreBlocking")
	CheckTypePreRetrying = CheckType("PreRetrying")
	CheckTypePreDead     = CheckType("PreDead")
	CheckTypePreReroute  = CheckType("PreReroute")
	CheckTypePreUpgrade  = CheckType("PreUpgrade")
	CheckTypePreDegrade  = CheckType("PreDegrade")

	CheckTypePostDiscard  = CheckType("PostDiscard")
	CheckTypePostPending  = CheckType("PostPending")
	CheckTypePostBlocking = CheckType("PostBlocking")
	CheckTypePostRetrying = CheckType("PostRetrying")
	CheckTypePostDead     = CheckType("PostDead")
	CheckTypePostReroute  = CheckType("PostReroute")
	CheckTypePostUpgrade  = CheckType("PostUpgrade")
	CheckTypePostDegrade  = CheckType("PostDegrade")
)

var defaultPreCheckTypesInTurn = []CheckType{CheckTypePreDiscard, CheckTypePreDead,
	CheckTypePreReroute, CheckTypePreUpgrade, CheckTypePreDegrade,
	CheckTypePreBlocking, CheckTypePrePending, CheckTypePreRetrying}

var DefaultPostCheckTypesInTurn = []CheckType{CheckTypePostDiscard, CheckTypePostDead,
	CheckTypePostReroute, CheckTypePostUpgrade, CheckTypePostDegrade,
	CheckTypePreBlocking, CheckTypePrePending, CheckTypePostRetrying}

var messageStatusToPostCheckTypeMap = map[messageStatus]CheckType{
	MessageStatusDiscard:           CheckTypePostDiscard,
	MessageStatusDead:              CheckTypePostDead,
	messageStatusNewReadyByReroute: CheckTypePostReroute,
	MessageStatusNewReadyByUpgrade: CheckTypePostUpgrade,
	MessageStatusNewReadyByDegrade: CheckTypePostDegrade,
	MessageStatusBlocking:          CheckTypePostBlocking,
	MessageStatusPending:           CheckTypePostPending,
	MessageStatusRetrying:          CheckTypePostRetrying,
}
