package inflight

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/qianjin/kodo-common/env"

	"github.com/qianjin/kodo-sample/rs/rsconfig"

	"github.com/qianjin/kodo-common/authkey"

	"github.com/qianjin/kodo-common/auth"
	"github.com/qianjin/kodo-common/client"
)

func TestEntryInfo_dev(t *testing.T) {
	body := "uid=1380469264&bucket=qj-test-rs_sample-rollback&key=test01.txt"
	rsconfig.SetupEnv("10.200.20.23:9433", "10.200.20.23:9433")
	test_Rs_EntryInfo(t, authkey.Dev_Key_admin, body)
}

func TestEntryInfo_prod(t *testing.T) {
	body := "itbl=504699156&key=fragments/z1.wypd.wypd/1631977359914-1631977367102.ts"
	rsconfig.SetupEnv(env.HostZ1Rs, env.HostZ1Rs)
	test_Rs_EntryInfo(t, authkey.Prod_Key_admin, body)
}

func test_Rs_EntryInfo(t *testing.T, authKey authkey.AuthKey, body string) {
	cli := client.NewClientWithHost(rsconfig.Env.Domain).
		WithKeys(authKey.AK, authKey.SK).
		WithSignType(auth.SignTypeQiniu)
	req := client.NewReq(http.MethodPost, "/entryinfo").
		RawQuery("").
		AddHeader("Host", rsconfig.Env.Host).
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		BodyStr(body)
	resp := cli.Call(req)
	fmt.Println(resp)
}
