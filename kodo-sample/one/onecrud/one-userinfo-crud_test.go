package onecrud

import (
	"net/http"
	"testing"

	"github.com/qianjin/kodo-common/auth"
	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-sample/one/oneconfig"
	"github.com/qianjin/kodo-security/kodokey"
	"github.com/stretchr/testify/assert"
)

func TestGetUserTuneSwitches(t *testing.T) {
	client.DebugMode = true
	oneconfig.SetupEnv("10.200.20.25:23200", "10.200.20.25:23200")
	oneCli := client.NewProxyClientWithHost(oneconfig.Env.Domain).
		WithProxyUser(auth.ProxyUserInfo{Uid: kodokey.Dev_UID_general_torage_011, Utype: 524316})

	getUserTuneSwitchesRespBody, getUserTuneSwitchesResp := GetUserTuneSwitches(oneCli)
	assert.Equal(t, http.StatusOK, getUserTuneSwitchesResp.StatusCode)
	assert.Equal(t, "", getUserTuneSwitchesRespBody.TuneSwitches)
}
