package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
)

type client struct {
	httpclient *http.Client

	host           string
	generatorToken func(req *http.Request) string
}

func newClientWithHost(host string, generatorToken func(req *http.Request) string) *client {
	client := &client{
		host:           host,
		httpclient:     &http.Client{},
		generatorToken: generatorToken,
	}
	return client
}

func (c *client) CallWithRet(req *Req, ret interface{}, opts ...ReqOption) *Resp {
	resp := c.Call(req, opts...)
	if resp != nil && len(resp.Body) > 0 {
		fmt.Printf("Data: %s\n", string(resp.Body))
		// 2xx才解析body
		if resp.StatusCode/100 == 2 {
			if err := json.Unmarshal(resp.Body, ret); err != nil {
				fmt.Printf("failed to unmarshal resp body, err: %v\n", err)
			}
		}
	}
	return resp
}

func (c *client) Call(req *Req, opts ...ReqOption) *Resp {
	for _, opt := range opts {
		opt(req)
	}
	var request *http.Request
	var err error
	if req.body != nil {
		request, err = http.NewRequest(req.method, "http://"+c.host+req.path, req.body)
	} else {
		request, err = http.NewRequest(req.method, "http://"+c.host+req.path, strings.NewReader(req.bodyStr))
	}
	if err != nil {
		return &Resp{StatusCode: http.StatusInternalServerError, Err: errors.New("client error: " + err.Error())}
	}

	if req.rawQuery != "" {
		request.URL.RawQuery = req.rawQuery
	}
	if hosts, ok := req.headers["Host"]; ok {
		request.Header.Add("Host", strings.Join(hosts, "; "))
	}
	for key, values := range req.headers {
		for _, v := range values {
			request.Header.Add(key, v)
		}
	}
	if c.generatorToken == nil {
		panic("client.tokenGenerator cannot be nil")
	}
	request.Header.Add("Authorization", c.generatorToken(request))
	if DebugMode {
		reqbytes, dumpErr := httputil.DumpRequestOut(request, true)
		fmt.Println(string(reqbytes), dumpErr)
	}
	response, err := c.httpclient.Do(request)
	if DebugMode {
		respbytes, dumpErr := httputil.DumpResponse(response, true)
		fmt.Println(string(respbytes), dumpErr)
	}
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
