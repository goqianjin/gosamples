package rspub

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/qianjin/kodo-security/kodokey"

	"github.com/qianjin/kodo-common/auth"
	"github.com/qianjin/kodo-common/client"
)

func TestEntryInfo_prod(t *testing.T) {
	body := "itbl=504699156&key=fragments/z1.wypd.wypd/1631977359914-1631977367102.ts"
	cli := client.NewClientWithHost("rs-z1.qbox.me").
		WithAuthKey(kodokey.Prod_AK_admin, kodokey.Prod_SK_admin).
		WithSignType(auth.SignTypeQiniu)
	req := client.NewReq(http.MethodPost, "/entryinfo").
		RawQuery("").
		AddHeader("Host", "rs-z1.qbox.me").
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		BodyStr(body)
	resp := cli.Call(req)
	fmt.Println(resp)
}

func TestEntryInfo_dev(t *testing.T) {
	body := "uid=1380469264&bucket=qj-test-rs_sample-rollback&key=test01.txt"
	//var Hosts = []string{"http://10.200.20.23:9433", "http://10.200.20.23:9433"}
	cli := client.NewClientWithHost("10.200.20.23:9433").
		WithAuthKey(kodokey.Dev_AK_admin, kodokey.Dev_SK_admin).
		WithSignType(auth.SignTypeQiniu)
	req := client.NewReq(http.MethodPost, "/entryinfo").
		RawQuery("").
		AddHeader("Host", "10.200.20.23:9433").
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		BodyStr(body)
	resp := cli.Call(req)
	fmt.Println(resp)
}
