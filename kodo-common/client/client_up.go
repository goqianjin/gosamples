package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/qiniu/go-sdk/v7/storage"

	"github.com/qianjin/kodo-common/auth"
)

type UpClient struct {
	httpclient *http.Client

	host      string
	generator *auth.UpTokenGenerator
}

func NewUpClientWithHost(host string) *UpClient {
	client := &UpClient{
		host:       host,
		httpclient: &http.Client{},
	}
	return client
}

func (c *UpClient) WithAuthKey(ak, sk string) *UpClient {
	c.generator = auth.NewUpTokenGenerator(ak, sk)
	return c
}

func (c *UpClient) WithPutPolicy(putPolicy *storage.PutPolicy) *UpClient {
	c.generator.WithPutPolicy(putPolicy)
	return c
}

func (c *UpClient) WithPutPolicyV2(putPolicy *auth.PutPolicyV2) *UpClient {
	c.generator.WithPutPolicyV2(putPolicy)
	return c
}

func (c *UpClient) CallWithRet(req *Req, ret interface{}) *Resp {
	resp := c.Call(req)
	if resp != nil && len(resp.Body) > 0 {
		fmt.Printf("Data: %s\n", string(resp.Body))
		if err := json.Unmarshal(resp.Body, ret); err != nil {
			fmt.Printf("failed to unmarshal resp body, err: %v", err)
		}
	}
	return resp
}

func (c *UpClient) Call(req *Req) *Resp {
	body := req.body
	if body == nil {
		body = strings.NewReader(req.bodyStr)
	}
	request, err := http.NewRequest(req.method, "http://"+c.host+req.path, body)
	if req.rawQuery != "" {
		request.URL.RawQuery = req.rawQuery
	}
	if hosts, ok := req.headers["Host"]; ok {
		request.Header.Add("Host", strings.Join(hosts, "; "))
	}
	for key, values := range req.headers {
		if key == "Content-Length" {
			continue
		}
		for _, v := range values {
			request.Header.Add(key, v)
		}
	}
	if contentLength, ok := req.headers["Content-Length"]; ok && len(contentLength) > 0 {
		request.ContentLength, _ = strconv.ParseInt(contentLength[0], 10, 64)
	}
	request.Header.Add("Authorization", c.generator.GenerateToken())
	response, err := c.httpclient.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("failed to readAll from response.Body: %v\n", err)
	}
	if response.StatusCode/100 != 2 {
		if err == nil {
			err = errors.New(string(data))
		} else {
			err = errors.New(string(data) + " --> " + err.Error())
		}
	}
	return &Resp{Body: data, StatusCode: response.StatusCode, Err: err, Headers: response.Header}
}
