package upcrud_partsv2

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"up/upconfig"
	"up/upmodel"

	"github.com/qianjin/kodo-common/client"
)

func InitParts(cli *client.UpClient, reqBody upmodel.InitPartsReq) (respBody upmodel.InitPartsResp, resp *client.Resp) {
	encodedKey := "~" // 注意默认值非空
	if reqBody.Key != "" {
		encodedKey = base64.URLEncoding.EncodeToString([]byte(reqBody.Key))
	}
	path := "/buckets/" + reqBody.Bucket + "/objects/" + encodedKey + "/uploads"

	bytesBody, err := json.Marshal(&reqBody)
	if err != nil {
		return
	}
	req := client.NewReq(http.MethodPost, path).
		RawQuery("").
		AddHeader("Host", upconfig.Env.Host).
		AddHeader("Content-Type", "application/json").
		Body(bytes.NewReader(bytesBody))
	resp = cli.CallWithRet(req, &respBody)
	return
}

func UploadParts(cli *client.UpClient, body io.Reader, reqBody upmodel.UploadPartsReq) (respBody upmodel.UploadPartsResp, resp *client.Resp) {
	encodedKey := "~" // 注意默认值非空
	if reqBody.Key != "" {
		encodedKey = base64.URLEncoding.EncodeToString([]byte(reqBody.Key))
	}
	path := "/buckets/" + reqBody.Bucket + "/objects/" + encodedKey + "/uploads/" + reqBody.UploadId + "/" + strconv.FormatInt(reqBody.PartNumber, 10)

	req := client.NewReq(http.MethodPut, path).
		RawQuery("").
		AddHeader("Host", upconfig.Env.Host).
		AddHeader("Content-Type", "application/json").
		Body(body)
	if reqBody.Size != 0 {
		req.AddHeader("Content-Length", strconv.Itoa(reqBody.Size))
	}
	resp = cli.CallWithRet(req, &respBody)
	return
}

func CompleteParts(cli *client.UpClient, reqBody upmodel.CompletePartsReq) (respBody upmodel.CompletePartsResp, resp *client.Resp) {
	type CompletePartBody struct {
		Parts      []upmodel.UploadPartInfo `json:"parts"`
		MimeType   string                   `json:"mimeType,omitempty"`
		Metadata   map[string]string        `json:"metadata,omitempty"`
		CustomVars map[string]string        `json:"customVars,omitempty"`
	}
	if reqBody.Extra == nil {
		reqBody.Extra = &upmodel.RputV2Extra{}
	}
	completePartBody := CompletePartBody{
		Parts:      reqBody.Extra.Progresses,
		MimeType:   reqBody.Extra.MimeType,
		Metadata:   reqBody.Extra.Metadata,
		CustomVars: make(map[string]string),
	}
	for k, v := range reqBody.Extra.CustomVars {
		if strings.HasPrefix(k, "x:") && v != "" {
			completePartBody.CustomVars[k] = v
		}
	}

	encodedKey := "~" // 注意默认值非空
	if reqBody.Key != "" {
		encodedKey = base64.URLEncoding.EncodeToString([]byte(reqBody.Key))
	}
	path := "/buckets/" + reqBody.Bucket + "/objects/" + encodedKey + "/uploads/" + reqBody.UploadId

	bytesBody, err := json.Marshal(&completePartBody)
	if err != nil {
		return
	}

	req := client.NewReq(http.MethodPost, path).
		RawQuery("").
		AddHeader("Host", upconfig.Env.Host).
		AddHeader("Content-Type", "application/json").
		AddHeader("Content-Length", string(len(bytesBody))).
		Body(bytes.NewReader(bytesBody))
	resp = cli.CallWithRet(req, &respBody)
	return
}
