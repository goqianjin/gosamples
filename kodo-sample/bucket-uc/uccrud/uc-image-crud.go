package uccrud

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-sample/bucket-uc/ucconfig"
	"github.com/qianjin/kodo-sample/bucket-uc/ucmodel"
)

func SetImage(cli *client.Client, reqBo ucmodel.SetImageReq) (respBo ucmodel.SetImageResp, resp *client.Resp) {
	path := fmt.Sprintf("/image/%s/from/%s", reqBo.Bucket, base64.URLEncoding.EncodeToString([]byte(reqBo.SiteURL)))
	reqQ := client.NewReq(http.MethodPost, path).
		RawQuery("").
		AddHeader("Host", ucconfig.Env.Host).
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		BodyStr("")

	resp = cli.CallWithRet(reqQ, &respBo)
	return
}
