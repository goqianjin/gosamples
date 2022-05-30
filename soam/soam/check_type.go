package soam

// ------ check type ------

type CheckType string

const (
	CheckTypePreDiscard  = CheckType("PreDiscard")
	CheckTypePrePending  = CheckType("PrePending")
	CheckTypePreBlocking = CheckType("PreBlocking")
	CheckTypePreRetrying = CheckType("PreRetrying")
	CheckTypePreDead     = CheckType("PreDead")
	CheckTypePreReRouter = CheckType("PreReRouter")
	CheckTypePreUpgrade  = CheckType("PreUpgrade")
	CheckTypePreDegrade  = CheckType("PreDegrade")

	CheckTypePostDiscard  = CheckType("PostDiscard")
	CheckTypePostPending  = CheckType("PostPending")
	CheckTypePostBlocking = CheckType("PostBlocking")
	CheckTypePostRetrying = CheckType("PostRetrying")
	CheckTypePostDead     = CheckType("PostDead")
	CheckTypePostReRouter = CheckType("PostReRouter")
	CheckTypePostUpgrade  = CheckType("PostUpgrade")
	CheckTypePostDegrade  = CheckType("PostDegrade")
)

var defaultPreCheckTypesInTurn = []CheckType{CheckTypePreDead, CheckTypePreDiscard,
	CheckTypePreReRouter, CheckTypePreUpgrade, CheckTypePreDegrade,
	CheckTypePreBlocking, CheckTypePrePending, CheckTypePreRetrying}

var DefaultPostCheckTypesInTurn = []CheckType{CheckTypePostDead, CheckTypePostDiscard,
	CheckTypePostReRouter, CheckTypePostUpgrade, CheckTypePostDegrade,
	CheckTypePreBlocking, CheckTypePrePending, CheckTypePostRetrying}
