package rscrud

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/qianjin/kodo-sample/rs/rsmodel"

	"github.com/qianjin/kodo-sample/rs/rsconfig"

	"github.com/qianjin/kodo-common/client"
)

func DeleteObject(cli *client.ManageClient, reqBody rsmodel.DeleteObjectReq) (respBody rsmodel.DeleteObjectResp, resp *client.Resp) {
	entry := fmt.Sprintf("%s:%s", reqBody.Bucket, reqBody.Key)
	path := fmt.Sprintf("/delete/%s", base64.URLEncoding.EncodeToString([]byte(entry)))
	reqD := client.NewReq(http.MethodPost, path).
		RawQuery("").
		AddHeader("Host", rsconfig.Env.Host).
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		BodyStr("")

	resp = cli.CallWithRet(reqD, &respBody)
	return
}

func GetStat(cli *client.ManageClient, reqBody rsmodel.GetObjectStatReq, opts ...client.ReqOption) (respBody rsmodel.GetObjectStatResp, resp *client.Resp) {
	entry := fmt.Sprintf("%s:%s", reqBody.Bucket, reqBody.Key)
	encodedEntry := base64.URLEncoding.EncodeToString([]byte(entry))
	path := fmt.Sprintf("/stat/%s", encodedEntry)
	rawQuery := ""
	if reqBody.NeedParts {
		rawQuery = "needparts=true"
	}
	req := client.NewReq(http.MethodPost, path).
		RawQuery(rawQuery).
		AddHeader("Host", rsconfig.Env.Host).
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		BodyStr("")

	for _, opt := range opts {
		opt(req)
	}

	resp = cli.CallWithRet(req, &respBody)
	return
}
