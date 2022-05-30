package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/qiniu/go-sdk/v7/storage"

	"github.com/qianjin/kodo-common/auth"
)

type UpClient struct {
	*client
	generator *auth.UpTokenGenerator
}

func NewUpClientWithHost(host string) *UpClient {
	cli := &UpClient{}
	client := newClientWithHost(host, cli.GenerateToken)
	cli.client = client
	return cli
}

func (c *UpClient) GenerateToken(_ *http.Request) string {
	return c.generator.GenerateToken()
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
			fmt.Printf("failed to unmarshal resp body, err: %v\n", err)
		}
	}
	return resp
}
