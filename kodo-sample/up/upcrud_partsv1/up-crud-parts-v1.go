package upcrud_partsv1

import (
	"encoding/base64"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-sample/up/upconfig"
	"github.com/qianjin/kodo-sample/up/upmodel"
)

func Mkblk(cli *client.UpClient, body io.Reader, reqBody upmodel.MkblkReq) (respBody upmodel.MkblkResp, resp *client.Resp) {
	path := "/mkblk/" + strconv.FormatInt(reqBody.BlockSize, 10)
	req := client.NewReq(http.MethodPost, path).
		RawQuery("").
		AddHeader("Host", upconfig.Env.Host).
		AddHeader("Content-Type", "application/octet-stream").
		AddHeader("Content-Length", strconv.FormatInt(reqBody.BodyLength, 10)).
		Body(body)
	resp = cli.CallWithRet(req, &respBody)
	return
}

func Bput(cli *client.UpClient, body io.Reader, reqBody upmodel.BputReq) (respBody upmodel.BputResp, resp *client.Resp) {
	path := "/bput/" + reqBody.Ctx + "/" + strconv.FormatUint(uint64(reqBody.Offset), 10)
	req := client.NewReq(http.MethodPost, path).
		RawQuery("").
		AddHeader("Host", upconfig.Env.Host).
		AddHeader("Content-Type", "application/json").
		Body(body)
	resp = cli.CallWithRet(req, &respBody)
	return
}

func Mkfile(cli *client.UpClient, reqBody upmodel.MkfileReq) (respBody upmodel.MkfileResp, resp *client.Resp) {
	url := "/mkfile/" + strconv.FormatInt(reqBody.Fsize, 10)
	if reqBody.Extra == nil {
		reqBody.Extra = &upmodel.RputExtra{}
	}
	if reqBody.Extra.MimeType != "" {
		url += "/mimeType/" + base64.URLEncoding.EncodeToString([]byte(reqBody.Extra.MimeType))
	}
	if reqBody.Key != "" {
		url += "/key/" + base64.URLEncoding.EncodeToString([]byte(reqBody.Key))
	}
	for k, v := range reqBody.Extra.Params {
		if (strings.HasPrefix(k, "x:") || strings.HasPrefix(k, "x-qn-meta-")) && v != "" {
			url += "/" + k + "/" + base64.URLEncoding.EncodeToString([]byte(v))
		}
	}
	ctxs := make([]string, len(reqBody.Extra.Progresses))
	for i, progress := range reqBody.Extra.Progresses {
		ctxs[i] = progress.Ctx
	}
	buf := strings.Join(ctxs, ",")
	body := strings.NewReader(buf)

	req := client.NewReq(http.MethodPost, url).
		RawQuery("").
		AddHeader("Host", upconfig.Env.Host).
		AddHeader("Content-Type", "application/json").
		//AddHeader("Content-Length", string(len(buf))).
		Body(body)
	resp = cli.CallWithRet(req, &respBody)
	return
}
