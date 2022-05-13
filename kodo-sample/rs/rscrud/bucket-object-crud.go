package rscrud

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/qianjin/kodo-sample/rs/rsmodel"

	"github.com/qianjin/kodo-sample/rs/rsconfig"

	"github.com/qianjin/kodo-common/client"
)

func DeleteObject(cli *client.Client, reqBody rsmodel.DeleteObjectReq) (respBody rsmodel.DeleteObjectResp, resp *client.Resp) {
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
