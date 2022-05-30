package client

import (
	"net/http"

	"github.com/qianjin/kodo-common/auth"
	"github.com/qianjin/kodo-common/authkey"
)

type ManageClient struct {
	*client
	generator *auth.ManagedTokenGenerator
}

func NewManageClientWithHost(host string) *ManageClient {
	cli := &ManageClient{}
	client := newClientWithHost(host, cli.GenerateToken)
	cli.client = client
	return cli

}

func (c *ManageClient) GenerateToken(req *http.Request) string {
	return c.generator.GenerateToken(req)
}

func (c *ManageClient) WithKey(key *authkey.AuthKey) *ManageClient {
	c.generator = auth.NewManagedTokenGenerator(key.AK, key.SK)
	return c
}

func (c *ManageClient) WithKeys(ak, sk string) *ManageClient {
	c.generator = auth.NewManagedTokenGenerator(ak, sk)
	return c
}

func (c *ManageClient) WithSignType(signType auth.SignType) *ManageClient {
	c.generator.WithSignType(signType)
	return c
}

func (c *ManageClient) WithSuInfo(uid, appId uint32) *ManageClient {
	c.generator.WithSuInfo(uid, appId)
	return c
}
