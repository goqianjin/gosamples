package soften

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerResult(t *testing.T) {
	assert.True(t, HandleStatusOk.gotoAction != HandleStatusFail.gotoAction)
}
