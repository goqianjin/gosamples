package checker

import (
	"errors"
	"fmt"

	"github.com/shenqianjin/soften-client-go/soften/internal"
)

var Validator = &validator{}

type validator struct {
}

func (v *validator) ValidateConsumeCheckpoint(checkpoints []Checkpoint) (map[internal.CheckType]*Checkpoint, error) {
	// 校验checker: checker可以在对应配置enable=false的情况下存在
	checkpointMap := make(map[internal.CheckType]*Checkpoint)
	for _, checkOpt := range checkpoints {
		if checkOpt.CheckType == "" {
			return nil, errors.New(" internal.CheckType can not be empty")
		}
		if v.isPreStatusCheckType(checkOpt.CheckType) {
			if checkOpt.Before == nil {
				return nil, errors.New(fmt.Sprintf("BeforeCheckFunc can not be nil for input checkOption: %s", checkOpt.CheckType))
			}
		} else if v.isPostStatusCheckType(checkOpt.CheckType) {
			if checkOpt.After == nil {
				return nil, errors.New(fmt.Sprintf("AfterCheckFunc can not be nil for input checkOption: %s", checkOpt.CheckType))
			}
		} else if v.isPreRerouteCheckType(checkOpt.CheckType) {
			if checkOpt.Before == nil {
				return nil, errors.New(fmt.Sprintf("BeforeForReroute can not be nil for input checkOption: %s", checkOpt.CheckType))
			}
		} else if v.isPostRerouteCheckType(checkOpt.CheckType) {
			if checkOpt.After == nil {
				return nil, errors.New(fmt.Sprintf("AfterForReroute can not be nil for input checkOption: %s", checkOpt.CheckType))
			}
		}
		checkpointMap[checkOpt.CheckType] = &checkOpt
	}
	// 一致性校验
	/*if conf.PendingEnable {
		if conf.Pending.CheckerMandatory && v.findCheckpointByType(checkpointMap, CheckTypePrePending, CheckTypePostPending) == nil {
			return nil, errors.New(fmt.Sprintf("[%s] checkOption is missing. please add one or disable the mandatory if necessary", message.StatusPending))
		}
	}
	if conf.BlockingEnable {
		if conf.Pending.CheckerMandatory && v.findCheckpointByType(checkpointMap, CheckTypePreBlocking, CheckTypePostBlocking) == nil {
			return nil, errors.New(fmt.Sprintf("[%s] checkOption is missing. please add one or disable the mandatory if necessary", message.StatusBlocking))
		}
	}
	if conf.RetryingEnable {
		if conf.Pending.CheckerMandatory && v.findCheckpointByType(checkpointMap, CheckTypePreRetrying, CheckTypePostRetrying) == nil {
			return nil, errors.New(fmt.Sprintf("[%s] checkOption is missing. please add one or disable the mandatory if necessary", message.StatusRetrying))
		}
	}*/
	return checkpointMap, nil
}

func (v *validator) findCheckpointByType(checkpointMap map[internal.CheckType]*Checkpoint, checkTypes ...internal.CheckType) *Checkpoint {
	for _, checkType := range checkTypes {
		if opt, ok := checkpointMap[checkType]; ok {
			return opt
		}
	}
	return nil
}

func (v *validator) isPreStatusCheckType(checkType internal.CheckType) bool {
	for _, ct := range PreCheckTypes() {
		if v.isPreRerouteCheckType(ct) {
			continue
		}
		if ct == checkType {
			return true
		}
	}
	return false
}

func (v *validator) isPostStatusCheckType(checkType internal.CheckType) bool {
	for _, ct := range DefaultPostCheckTypes() {
		if v.isPostRerouteCheckType(ct) {
			continue
		}
		if ct == checkType {
			return true
		}
	}
	return false
}

func (v *validator) isPreRerouteCheckType(checkType internal.CheckType) bool {
	return checkType == CheckTypePreReroute
}

func (v *validator) isPostRerouteCheckType(checkType internal.CheckType) bool {
	return checkType == CheckTypePostReroute
}
