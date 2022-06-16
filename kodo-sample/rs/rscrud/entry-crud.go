package rscrud

import (
	"net/http"

	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-sample/rs/rsconfig"
	"github.com/qianjin/kodo-sample/rs/rsmodel"
)

func EntryInfo(cli *client.ManageClient, reqBody rsmodel.GetEntryInfoReq, opts ...client.ReqOption) (respBody rsmodel.GetEntryInfoResp, resp *client.Resp) {
	body := "itbl=504699156&key=fragments/z1.wypd.wypd/1631977359914-1631977367102.ts"

	req := client.NewReq(http.MethodPost, "/entryinfo").
		RawQuery("").
		AddHeader("Host", rsconfig.Env.Host).
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		BodyStr(body)

	resp = cli.Call(req)
	return
}
