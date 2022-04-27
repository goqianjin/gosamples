package sisyphuscrud

import (
	"net/http"
	"sisyphus/sisyphusconfig"
	"sisyphus/sisyphusmodel"

	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-common/json"
)

func CreateDualTask(cli *client.Client, reqBody sisyphusmodel.CreateDualTaskReq) (respBody sisyphusmodel.CreateDualTaskResp, resp *client.Resp) {
	if reqBody.Name == "" {
		reqBody.Name = sisyphusconfig.GenerateTaskName()
	}

	req := client.NewReq(http.MethodPost, "/dualsync/task/create").
		RawQuery("").
		AddHeader("Host", sisyphusconfig.Env.Host).
		AddHeader("Content-Type", "application/json").
		BodyStr(json.ToJson(reqBody))
	resp = cli.CallWithRet(req, &respBody)
	return
}

func QueryDualTask(cli *client.Client, reqBody sisyphusmodel.QueryDualTaskReq) (respBody sisyphusmodel.QueryDualTaskResp, resp *client.Resp) {
	// refer storage.BucketManager{}

	req := client.NewReq(http.MethodPost, "/dualsync/task/query").
		RawQuery("").
		AddHeader("Host", sisyphusconfig.Env.Host).
		AddHeader("Content-Type", "application/json").
		BodyStr(json.ToJson(reqBody))

	resp = cli.CallWithRet(req, &respBody)
	return
}

func StopDualTask(cli *client.Client, reqBody sisyphusmodel.StopDualTaskReq) (resp *client.Resp) {
	req := client.NewReq(http.MethodPost, "/dualsync/task/stop").
		RawQuery("").
		AddHeader("Host", sisyphusconfig.Env.Host).
		AddHeader("Content-Type", "application/json").
		BodyStr(json.ToJson(reqBody))

	resp = cli.Call(req)
	return
}

func StartDualTask(cli *client.Client, reqBody sisyphusmodel.StartDualTaskReq) (resp *client.Resp) {
	req := client.NewReq(http.MethodPost, "/dualsync/task/start").
		RawQuery("").
		AddHeader("Host", sisyphusconfig.Env.Host).
		AddHeader("Content-Type", "application/json").
		BodyStr(json.ToJson(reqBody))

	resp = cli.Call(req)
	return
}

func DeleteDualTask(cli *client.Client, reqBody sisyphusmodel.DeleteDualTaskReq) (resp *client.Resp) {
	req := client.NewReq(http.MethodPost, "/dualsync/task/delete").
		RawQuery("").
		AddHeader("Host", sisyphusconfig.Env.Host).
		AddHeader("Content-Type", "application/json").
		BodyStr(json.ToJson(reqBody))

	resp = cli.Call(req)
	return
}
