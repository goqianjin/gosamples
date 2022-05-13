package upcrud_form

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"

	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-sample/up/upconfig"
	"github.com/qianjin/kodo-sample/up/upmodel"
)

func FormUpload(cli *client.UpClient, body io.Reader, reqBody upmodel.FormUploadReq) (respBody upmodel.FormUploadResp, resp *client.Resp) {
	path := "/"

	multipartWriter, newBody, err := composeFormBody(body, reqBody)
	if err != nil {
		return
	}
	req := client.NewReq(http.MethodPost, path).
		RawQuery("").
		AddHeader("Host", upconfig.Env.Host).
		AddHeader("Content-Type", multipartWriter.FormDataContentType()).
		Body(newBody)
	resp = cli.CallWithRet(req, &respBody)
	return
}

func composeFormBody(body io.Reader, reqBody upmodel.FormUploadReq) (multipartWriter *multipart.Writer, bodyReader io.Reader, err error) {
	// multipart writer
	bodyBuffer := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyBuffer)
	//token
	if err = writer.WriteField("token", reqBody.UploadToken); err != nil {
		return
	}
	//key
	if reqBody.Key != "" {
		if err = writer.WriteField("key", reqBody.Key); err != nil {
			return
		}
	}
	//extra.Params
	if reqBody.Extra != nil && reqBody.Extra.Params != nil {
		for k, v := range reqBody.Extra.Params {
			if (strings.HasPrefix(k, "x:") || strings.HasPrefix(k, "x-qn-meta-")) && v != "" {
				err = writer.WriteField(k, v)
				if err != nil {
					return
				}
			}
		}
	}

	/*var dataReader io.Reader
	h := crc32.NewIEEE()
	dataReader = io.TeeReader(form.Data, h)
	crcReader := newCrc32Reader(writer.Boundary(), h)*/

	//write file
	replacer := strings.NewReplacer("\\", "\\\\", `"`, "\\\"")
	head := make(textproto.MIMEHeader)
	head.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`,
		replacer.Replace(reqBody.FileName)))
	if reqBody.Extra != nil && reqBody.Extra.MimeType != "" {
		head.Set("Content-Type", reqBody.Extra.MimeType)
	}

	_, err = writer.CreatePart(head)
	if err != nil {
		return
	}
	head = nil
	lastLine := fmt.Sprintf("\r\n--%s--\r\n", writer.Boundary())
	r := strings.NewReader(lastLine)

	bodyLen := int64(-1)
	if reqBody.FileSize >= 0 {
		bodyLen = int64(bodyBuffer.Len()) + reqBody.FileSize + int64(len(lastLine))
		//bodyLen += crcReader.length()
	}
	_ = bodyLen

	mr := io.MultiReader(bodyBuffer, body, r)
	bodyBuffer = nil
	//dataReader = nil
	//crcReader = nil
	r = nil

	formBytes, err := ioutil.ReadAll(mr)
	if err != nil {
		return
	}
	mr = nil

	getBodyReader := func() (io.Reader, error) {
		var formReader io.Reader = bytes.NewReader(formBytes)
		return formReader, nil
	}
	/*getBodyReadCloser := func() (io.ReadCloser, error) {
		reader, err := getBodyReader()
		if err != nil {
			return nil, err
		}
		return ioutil.NopCloser(reader), nil
	}*/
	bodyReader, err = getBodyReader()
	if err != nil {
		return
	}
	return writer, bodyReader, err
}
