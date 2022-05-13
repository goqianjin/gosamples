package inflight

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/qianjin/kodo-common/env"

	"github.com/qianjin/kodo-common/auth"
	"github.com/qianjin/kodo-common/authkey"
	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-sample/rs/rsconfig"
)

func Test_Rs_Stats_Dev(t *testing.T) {
	rsconfig.SetupEnv("10.200.20.23:9433", "10.200.20.23:9433")
	test_Rs_Stats(t, authkey.Dev_Key_general_storage_011)
}

func Test_Rs_Stats_Prod(t *testing.T) {
	rsconfig.SetupEnv(env.HostDefaultRs, env.HostDefaultRs)
	test_Rs_Stats(t, authkey.Prod_Key_kodolog)
}

func test_Rs_Stats(t *testing.T, authKey authkey.AuthKey) {

	// rs stats
	cli := client.NewClientWithHost(rsconfig.Env.Domain).
		WithKeys(authKey.AK, authKey.SK).
		WithSignType(auth.SignTypeQiniu)

	bucket := "qianjin-bucket-20220513173323932606"
	key := "test01.txt-20220513173323"
	path := "/stat/" + base64.URLEncoding.EncodeToString([]byte(bucket+":"+key))
	req := client.NewReq(http.MethodPost, path).
		RawQuery("needparts=true").
		AddHeader("Host", rsconfig.Env.Host).
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		BodyStr("")

	resp := cli.Call(req)
	fmt.Println("responseCode: " + strconv.Itoa(resp.StatusCode) + ", body: " + string(resp.Body) + ", err: " + resp.Err.Error())
}
