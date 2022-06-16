package rscrud

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/qianjin/kodo-common/auth"
	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-sample/rs/rsconfig"
	"github.com/qianjin/kodo-sample/rs/rsmodel"
	"github.com/qianjin/kodo-security/kodokey"
	"github.com/stretchr/testify/assert"
)

func TestEntryInfo(t *testing.T) {
	rsconfig.SetupEnv("10.200.20.23:9433", "10.200.20.23:9433")

	cli := client.NewManageClientWithHost(rsconfig.Env.Domain).
		WithKeys(kodokey.Dev_AK_admin, kodokey.Dev_SK_admin).
		WithSignType(auth.SignTypeQiniuAdmin).WithSuInfo(1380538466, 0)

	itbl := uint32(1234)
	bucket := "uid=1380469264&bucket=qj-test-rs_sample-rollback&key=test01.txt"
	key := "test01.txt"
	respBody, resp := EntryInfo(cli, rsmodel.GetEntryInfoReq{Itbl: itbl, Bucket: bucket, Key: key})
	fmt.Printf("get entry info: %v (for itbl:%d, key: %s)\n", respBody, itbl, key)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGet2(t *testing.T) {
	rsconfig.SetupEnv("127.0.0.1:9433", "127.0.0.1:9433")

	body := "itbl=504699156&key=fragments/z1.wypd.wypd/1631977359914-1631977367102.ts"

	cli := client.NewManageClientWithHost(rsconfig.Env.Domain).
		WithKeys(kodokey.Dev_AK_admin, kodokey.Dev_SK_admin).
		WithSignType(auth.SignTypeQiniuAdmin).WithSuInfo(1380469264, 1)
	req := client.NewReq(http.MethodPost, "/entryinfo").
		RawQuery("").
		AddHeader("Host", rsconfig.Env.Host).
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		BodyStr(body)
	var respBody rsmodel.GetEntryInfoResp
	resp := cli.CallWithRet(req, &respBody)
	fmt.Println(resp)
}
