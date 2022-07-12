package checker

import (
	"github.com/shenqianjin/soften-client-go/soften/internal"
)

// ------ produce check type ------

const (
	ProduceCheckTypeDiscard  = internal.CheckType("Discard")
	ProduceCheckTypePending  = internal.CheckType("Pending")
	ProduceCheckTypeBlocking = internal.CheckType("Blocking")
	ProduceCheckTypeRetrying = internal.CheckType("Retrying")
	ProduceCheckTypeDead     = internal.CheckType("Dead")
	ProduceCheckTypeUpgrade  = internal.CheckType("Upgrade")
	ProduceCheckTypeDegrade  = internal.CheckType("Degrade")
	ProduceCheckTypeRoute    = internal.CheckType("Route")
)

func ProduceCheckTypes() []internal.CheckType {
	values := []internal.CheckType{ProduceCheckTypeDiscard, ProduceCheckTypeDead,
		ProduceCheckTypeRoute, ProduceCheckTypeUpgrade, ProduceCheckTypeDegrade,
		ProduceCheckTypeBlocking, ProduceCheckTypePending, ProduceCheckTypeRetrying}
	return values
}

// ------ consume check type ------

const (
	CheckTypePreDiscard  = internal.CheckType("PreDiscard")
	CheckTypePrePending  = internal.CheckType("PrePending")
	CheckTypePreBlocking = internal.CheckType("PreBlocking")
	CheckTypePreRetrying = internal.CheckType("PreRetrying")
	CheckTypePreDead     = internal.CheckType("PreDead")
	CheckTypePreUpgrade  = internal.CheckType("PreUpgrade")
	CheckTypePreDegrade  = internal.CheckType("PreDegrade")
	CheckTypePreReroute  = internal.CheckType("PreReroute")

	CheckTypePostDiscard  = internal.CheckType("PostDiscard")
	CheckTypePostPending  = internal.CheckType("PostPending")
	CheckTypePostBlocking = internal.CheckType("PostBlocking")
	CheckTypePostRetrying = internal.CheckType("PostRetrying")
	CheckTypePostDead     = internal.CheckType("PostDead")
	CheckTypePostUpgrade  = internal.CheckType("PostUpgrade")
	CheckTypePostDegrade  = internal.CheckType("PostDegrade")
	CheckTypePostReroute  = internal.CheckType("PostReroute")
)

func PreCheckTypes() []internal.CheckType {
	values := []internal.CheckType{CheckTypePreDiscard, CheckTypePreDead,
		CheckTypePreReroute, CheckTypePreUpgrade, CheckTypePreDegrade,
		CheckTypePreBlocking, CheckTypePrePending, CheckTypePreRetrying}
	return values
}

func DefaultPostCheckTypes() []internal.CheckType {
	values := []internal.CheckType{CheckTypePostDiscard, CheckTypePostDead,
		CheckTypePostReroute, CheckTypePostUpgrade, CheckTypePostDegrade,
		CheckTypePreBlocking, CheckTypePrePending, CheckTypePostRetrying}
	return values
}
