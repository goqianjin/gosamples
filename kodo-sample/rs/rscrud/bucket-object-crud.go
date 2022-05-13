package bucketcrud

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-sample/bucket/bucketconfig"
	"github.com/qianjin/kodo-sample/bucket/bucketmodel"
)

func DeleteObject(cli *client.Client, reqBody bucketmodel.DeleteObjectReq) (respBody bucketmodel.DeleteObjectResp, resp *client.Resp) {
	entry := fmt.Sprintf("%s:%s", reqBody.Bucket, reqBody.Key)
	path := fmt.Sprintf("/delete/%s", base64.URLEncoding.EncodeToString([]byte(entry)))
	reqD := client.NewReq(http.MethodPost, path).
		RawQuery("").
		AddHeader("Host", bucketconfig.Env.Host).
		//AddHeader("Content-Type", "application/x-www-form-urlencoded").
		BodyStr("")

	resp = cli.CallWithRet(reqD, &respBody)
	return
}
