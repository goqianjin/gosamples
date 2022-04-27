package bucketcrud

import (
	"fmt"
	"net/http"

	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-sample/bucket/bucketconfig"
	"github.com/qianjin/kodo-sample/bucket/bucketmodel"
	"github.com/qiniu/go-sdk/v7/storage"
)

func Create(cli *client.Client) (bucket string, resp *client.Resp) {
	return CreateWithOption(cli, bucketmodel.NewCreateOption().WithRegion("z0"))
}
func CreateWithOption(cli *client.Client, option *bucketmodel.CreateOption) (bucket string, resp *client.Resp) {
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

func Query(cli *client.Client, bucket string) (bucketInfoQ storage.BucketInfo, resp *client.Resp) {
	// refer storage.BucketManager{}
	reqQ := client.NewReq(http.MethodPost, "/v2/bucketInfo").
		RawQuery(fmt.Sprintf("bucket=%s", bucket)).
		AddHeader("Host", bucketconfig.Env.Host).
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		BodyStr("")

	resp = cli.CallWithRet(reqQ, &bucketInfoQ)
	return
}

func Delete(cli *client.Client, bucket string) (resp *client.Resp) {
	reqD := client.NewReq(http.MethodPost, fmt.Sprintf("/drop/%s", bucket)).
		RawQuery("").
		AddHeader("Host", bucketconfig.Env.Host).
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		BodyStr("")

	resp = cli.Call(reqD)
	return
}

func List(cli *client.Client) (buckets []string, resp *client.Resp) {
	req := client.NewReq(http.MethodPost, "/buckets").
		RawQuery("share=false").
		AddHeader("Host", bucketconfig.Env.Host).
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		BodyStr("")

	buckets = make([]string, 0)
	resp = cli.CallWithRet(req, &buckets)
	return
}
