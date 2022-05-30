package bucketcrud

import (
	"fmt"
	"net/http"

	"github.com/qianjin/kodo-sample/bucket/bucketmodel"

	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-sample/bucket/bucketconfig"
)

func ListDomainsV3(cli *client.ManageClient, bucket string, opts ...client.ReqOption) (respBody bucketmodel.ListBucketDomainsResp, resp *client.Resp) {
	path := fmt.Sprintf("/v3/domains?tbl=%s", bucket)
	req := client.NewReq(http.MethodPost, path).
		RawQuery("").
		AddHeader("Host", bucketconfig.Env.Host).
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		BodyStr("")

	resp = cli.CallWithRet(req, &respBody, opts...)
	return
}
func ListDomainsV7(cli *client.ManageClient, bucket string, opts ...client.ReqOption) (respBody bucketmodel.ListBucketDomainsResp, resp *client.Resp) {
	path := fmt.Sprintf("/v7/domain/list?tbl=%s", bucket)
	bucket = bucketconfig.GenerateBucketName()

	req := client.NewReq(http.MethodPost, path).
		RawQuery("").
		AddHeader("Host", bucketconfig.Env.Host).
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		BodyStr("")

	resp = cli.CallWithRet(req, &respBody, opts...)
	return
}

func SetDomain(cli *client.ManageClient) (bucket string, resp *client.Resp) {
	pathPattern := "/publish"
	bodyStr := "domain=storage01.ylwl-hlj-1.qiniu-solutions.com&tbl=qiniu-storage&domaintype=1&apiscope=0"
	bucket = bucketconfig.GenerateBucketName()

	req := client.NewReq(http.MethodPost, pathPattern).
		RawQuery("").
		AddHeader("Host", "192.168.200.111:10221").
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		BodyStr(bodyStr)

	resp = cli.Call(req)
	return
}
