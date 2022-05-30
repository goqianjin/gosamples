package uccrud

import (
	"fmt"
	"net/http"

	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-sample/bucket-uc/ucconfig"
	"github.com/qiniu/go-sdk/v7/storage"
)

func Query(cli *client.ManageClient, bucket string) (respBody storage.BucketInfo, resp *client.Resp) {
	// refer storage.BucketManager{}
	req := client.NewReq(http.MethodPost, "/v2/bucketInfo").
		RawQuery(fmt.Sprintf("bucket=%s", bucket)).
		AddHeader("Host", ucconfig.Env.Host).
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		BodyStr("")

	resp = cli.CallWithRet(req, &respBody)
	return
}
