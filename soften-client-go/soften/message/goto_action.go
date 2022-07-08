package message

import (
	"errors"
	"fmt"

	"github.com/shenqianjin/soften-client-go/soften/internal"
)

const (
	GotoBlocking = internal.MessageGoto("Blocking")
	GotoPending  = internal.MessageGoto("Pending")
	GotoRetrying = internal.MessageGoto("Retrying")
	GotoDead     = internal.MessageGoto("Dead")
	GotoDone     = internal.MessageGoto("Done")
	GotoDiscard  = internal.MessageGoto("Discard")
	GotoUpgrade  = internal.MessageGoto("Upgrade")
	GotoDegrade  = internal.MessageGoto("Degrade")
	// Reroute 只能通过checkpoint实现，不能在 HandleResult 中显示指定为 GotoAction
	//GotoReroute = internal.MessageGoto("Reroute")
)

func GotoActionOf(action string) (internal.MessageGoto, error) {
	for _, v := range GotoActionValues() {
		if v.String() == action {
			return v, nil
		}
	}
	return "", errors.New(fmt.Sprintf("invalid message goto action: %s", action))
}

func GotoActionValues() []internal.MessageGoto {
	values := []internal.MessageGoto{
		GotoBlocking, GotoPending, GotoRetrying,
		GotoDead, GotoDone, GotoDiscard,
		GotoUpgrade, GotoDegrade,
	}
	return values
}
