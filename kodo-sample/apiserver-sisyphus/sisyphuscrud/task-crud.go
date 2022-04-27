package sisyphuscrud

import (
	"net/http"
	"sisyphus/sisyphusconfig"
	"sisyphus/sisyphusmodel"

	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-common/json"
)

func Create(cli *client.Client, reqBody sisyphusmodel.CreateReq) (respBody sisyphusmodel.CreateResp, resp *client.Resp) {
	if reqBody.Name == "" {
		reqBody.Name = sisyphusconfig.GenerateTaskName()
	}

	req := client.NewReq(http.MethodPost, "/transfer/task/create").
		RawQuery("").
		AddHeader("Host", sisyphusconfig.Env.Host).
		AddHeader("Content-Type", "application/json").
		BodyStr(json.ToJson(reqBody))
	resp = cli.CallWithRet(req, &respBody)
	return
}

func Query(cli *client.Client, reqBody sisyphusmodel.QueryReq) (respBody sisyphusmodel.QueryResp, resp *client.Resp) {
	// refer storage.BucketManager{}

	req := client.NewReq(http.MethodPost, "/transfer/task/query").
		RawQuery("").
		AddHeader("Host", sisyphusconfig.Env.Host).
		AddHeader("Content-Type", "application/json").
		BodyStr(json.ToJson(reqBody))

	resp = cli.CallWithRet(req, &respBody)
	return
}

func Stop(cli *client.Client, reqBody sisyphusmodel.StopReq) (resp *client.Resp) {
	req := client.NewReq(http.MethodPost, "/transfer/task/stop").
		RawQuery("").
		AddHeader("Host", sisyphusconfig.Env.Host).
		AddHeader("Content-Type", "application/json").
		BodyStr(json.ToJson(reqBody))

	resp = cli.Call(req)
	return
}

func Start(cli *client.Client, reqBody sisyphusmodel.StartReq) (resp *client.Resp) {
	req := client.NewReq(http.MethodPost, "/transfer/task/start").
		RawQuery("").
		AddHeader("Host", sisyphusconfig.Env.Host).
		AddHeader("Content-Type", "application/json").
		BodyStr(json.ToJson(reqBody))

	resp = cli.Call(req)
	return
}

func Delete(cli *client.Client, reqBody sisyphusmodel.DeleteReq) (resp *client.Resp) {
	req := client.NewReq(http.MethodPost, "/transfer/task/delete").
		RawQuery("").
		AddHeader("Host", sisyphusconfig.Env.Host).
		AddHeader("Content-Type", "application/json").
		BodyStr(json.ToJson(reqBody))

	resp = cli.Call(req)
	return
}
