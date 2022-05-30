package client

import (
	"net/http"

	"github.com/qianjin/kodo-common/proxyuser"

	"github.com/qianjin/kodo-common/auth"
)

type ProxyClient struct {
	*client
	generator *auth.QiniuProxyTokenGenerator
}

func NewProxyClientWithHost(host string) *ProxyClient {
	cli := &ProxyClient{}
	client := newClientWithHost(host, cli.GenerateToken)
	cli.client = client
	return cli
}

func (c *ProxyClient) GenerateToken(_ *http.Request) string {
	return c.generator.GenerateToken()
}

func (c *ProxyClient) WithProxyUser(user proxyuser.ProxyUser) *ProxyClient {
	c.generator = auth.NewQiniuProxyTokenGenerator(user)
	return c
}
