package onecrud

import (
	"net/http"

	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-common/json"
	"github.com/qianjin/kodo-sample/one/oneconfig"
	"github.com/qianjin/kodo-sample/one/onemodel"
)

func PutUserTuneSwitches(cli *client.ProxyClient, reqBo onemodel.PutUserTuneSwitchesReq) (respBody onemodel.PutUserTuneSwitchesResp, resp *client.Resp) {
	// refer storage.BucketManager{}
	req := client.NewReq(http.MethodPut, "/user/tune/switches").
		RawQuery("").
		AddHeader("Host", oneconfig.Env.Host).
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		BodyStr(json.ToJson(reqBo))

	resp = cli.CallWithRet(req, &respBody)
	return
}

func GetUserTuneSwitches(cli *client.ProxyClient, opts ...client.ReqOption) (respBody onemodel.GetUserTuneSwitchesResp, resp *client.Resp) {
	// refer storage.BucketManager{}
	req := client.NewReq(http.MethodGet, "/user/tune/switches").
		RawQuery("").
		AddHeader("Host", oneconfig.Env.Host).
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		BodyStr("")

	for _, opt := range opts {
		opt(req)
	}

	resp = cli.CallWithRet(req, &respBody)
	return
}
