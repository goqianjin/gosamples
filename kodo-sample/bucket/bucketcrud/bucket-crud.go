package bucketcrud

import (
	"fmt"
	"net/http"

	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-sample/bucket/bucketconfig"
	"github.com/qianjin/kodo-sample/bucket/bucketmodel"
)

func Create(cli *client.ManageClient) (bucket string, resp *client.Resp) {
	return CreateWithOption(cli, bucketmodel.NewCreateOption().WithRegion("z0"))
}
func CreateWithOption(cli *client.ManageClient, option *bucketmodel.CreateOption) (bucket string, resp *client.Resp) {
	pathPattern := "/mkbucketv3/%s/region/%s/nodomain/true"
	bucket = bucketconfig.GenerateBucketName()

	req := client.NewReq(http.MethodPost, fmt.Sprintf(pathPattern, bucket, option.Region)).
		RawQuery("").
		AddHeader("Host", bucketconfig.Env.Host).
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		BodyStr("")

	resp = cli.Call(req)
	return
}

func Query(cli *client.ManageClient, bucket string, opts ...client.ReqOption) (respBody bucketmodel.QueryBucketResp, resp *client.Resp) {
	// refer storage.BucketManager{}
	path := fmt.Sprintf("/bucket/%s", bucket)
	req := client.NewReq(http.MethodGet, path).
		RawQuery("").
		AddHeader("Host", bucketconfig.Env.Host).
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		BodyStr("")

	for _, opt := range opts {
		opt(req)
	}

	resp = cli.CallWithRet(req, &respBody)
	return
}

func Delete(cli *client.ManageClient, bucket string) (resp *client.Resp) {
	reqD := client.NewReq(http.MethodPost, fmt.Sprintf("/drop/%s", bucket)).
		RawQuery("").
		AddHeader("Host", bucketconfig.Env.Host).
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		BodyStr("")

	resp = cli.Call(reqD)
	return
}

func List(cli *client.ManageClient) (buckets []string, resp *client.Resp) {
	req := client.NewReq(http.MethodPost, "/buckets").
		RawQuery("share=false").
		AddHeader("Host", bucketconfig.Env.Host).
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		BodyStr("")

	buckets = make([]string, 0)
	resp = cli.CallWithRet(req, &buckets)
	return
}
