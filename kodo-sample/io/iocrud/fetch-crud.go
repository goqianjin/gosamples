package iocrud

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-sample/io/ioconfig"
	"github.com/qianjin/kodo-sample/io/iomodel"
)

func Fetch(cli *client.ManageClient, reqBo iomodel.FetchReq, options ...client.ReqOption) (respBo iomodel.FetchResp, resp *client.Resp) {
	entry := reqBo.Bucket
	if reqBo.Key != "" {
		entry += ":" + reqBo.Key
	}
	path := fmt.Sprintf("/fetch/%s/to/%s",
		base64.URLEncoding.EncodeToString([]byte(reqBo.ResURL)), base64.URLEncoding.EncodeToString([]byte(entry)))
	reqQ := client.NewReq(http.MethodPost, path).
		RawQuery("").
		AddHeader("Host", ioconfig.Env.Host).
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		BodyStr("")
	// options
	for _, opt := range options {
		opt(reqQ)
	}

	resp = cli.CallWithRet(reqQ, &respBo)
	return
}
