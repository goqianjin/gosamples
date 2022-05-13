package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/qianjin/kodo-common/auth"
	"github.com/qianjin/kodo-common/authkey"
)

type Client struct {
	httpclient *http.Client

	host      string
	generator *auth.ManagedTokenGenerator
}

func NewClientWithHost(host string) *Client {
	client := &Client{
		host:       host,
		httpclient: &http.Client{},
	}
	return client
}

func (c *Client) WithKey(key *authkey.AuthKey) *Client {
	c.generator = auth.NewManagedTokenGenerator(key.AK, key.SK)
	return c
}

func (c *Client) WithKeys(ak, sk string) *Client {
	c.generator = auth.NewManagedTokenGenerator(ak, sk)
	return c
}

func (c *Client) WithSignType(signType auth.SignType) *Client {
	c.generator.WithSignType(signType)
	return c
}

func (c *Client) WithSuInfo(uid, appId uint32) *Client {
	c.generator.WithSuInfo(uid, appId)
	return c
}

func (c *Client) CallWithRet(req *Req, ret interface{}) *Resp {
	resp := c.Call(req)
	if resp != nil && len(resp.Body) > 0 {
		fmt.Printf("Data: %s\n", string(resp.Body))
		if err := json.Unmarshal(resp.Body, ret); err != nil {
			fmt.Printf("failed to unmarshal resp body, err: %v\n", err)
		}
	}
	return resp
}

func (c *Client) Call(req *Req) *Resp {
	request, err := http.NewRequest(req.method, "http://"+c.host+req.path, strings.NewReader(req.bodyStr))
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
	request.Header.Add("Authorization", c.generator.GenerateToken(request))
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
