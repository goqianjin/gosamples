package inflight

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/qianjin/kodo-common/auth"
	"github.com/qianjin/kodo-common/authkey"
	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-common/env"
	"github.com/qianjin/kodo-sample/bucket/bucketconfig"
	"github.com/qianjin/kodo-sample/rs/rsconfig"
	"github.com/qianjin/kodo-sample/up/upconfig"
)

func TestKODO9272_ZoneProxy_RsStats_AuthExpire_MultiCases_dev(t *testing.T) {
	rsconfig.SetupEnv("10.200.20.25:12501", "10.200.20.25:12501")
	bucketconfig.SetupEnv("10.200.20.25:10221", "10.200.20.25:10221")
	upconfig.SetupEnv("10.200.20.23:5010", "10.200.20.23:5010")
	testKODO9272_RsStats_AuthExpire_MultiCases(t, authkey.Dev_Key_general_storage_011)
}

func TestKODO9272_RsStats_AuthExpire_MultiCases_dev(t *testing.T) {
	rsconfig.SetupEnv("10.200.20.23:9433", "10.200.20.23:9433")
	bucketconfig.SetupEnv("10.200.20.25:10221", "10.200.20.25:10221")
	upconfig.SetupEnv("10.200.20.23:5010", "10.200.20.23:5010")
	testKODO9272_RsStats_AuthExpire_MultiCases(t, authkey.Dev_Key_general_storage_011)
}

func TestKODO9272_RsStats_AuthExpire_MultiCases_prod(t *testing.T) {
	rsconfig.SetupEnv(env.HostDefaultRs, env.HostDefaultRs)
	bucketconfig.SetupEnv(env.HostDefaultUc, env.HostDefaultUc)
	upconfig.SetupEnv(env.HostZ0Up, env.HostZ0Up)
	testKODO9272_RsStats_AuthExpire_MultiCases(t, authkey.Prod_Key_admin)
}

func testKODO9272_RsStats_AuthExpire_MultiCases(t *testing.T, authKey authkey.AuthKey) {
	// prepare bucket data

	bucket := "qianjin-bucket-20220513173323932606"
	key := "test01.txt-20220513173323"
	// rs stats
	fmt.Println("bucket: " + bucket + ", key: " + key)
	path := "/stat/" + base64.URLEncoding.EncodeToString([]byte(bucket+":"+key))
	cli := client.NewClientWithHost(rsconfig.Env.Domain).
		WithKeys(authKey.AK, authKey.SK).
		WithSignType(auth.SignTypeQiniu)
	req := client.NewReq(http.MethodPost, path).
		RawQuery("needparts=true").
		AddHeader("Host", rsconfig.Env.Host).
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		BodyStr("")

	type UserCase struct {
		subject        string
		headerXDate    string
		headerXExpires string
		expectedCode   int
		expectedErr    string
		queryXDate     string
		queryXExpires  string
	}
	cases := []UserCase{
		{subject: "无X-Qiniu-Date & 无X-Qiniu-Expires",
			headerXDate: "", headerXExpires: "", expectedCode: 200, expectedErr: ""},
	}
	for _, c := range cases {
		copiedReq := req.DeepClone()
		if c.headerXDate != "" {
			copiedReq.SetHeader("X-Qiniu-Date", c.headerXDate)
		}
		if c.headerXExpires != "" {
			copiedReq.SetHeader("X-Qiniu-Expires", c.headerXExpires)
		}
		// query
		query := copiedReq.GetRawQuery()
		if c.queryXDate != "" {
			if query != "" {
				query = query + "&"
			}
			query = "X-Qiniu-Date=" + c.queryXDate
		}
		if c.queryXExpires != "" {
			if query != "" {
				query = query + "&"
			}
			query = query + "X-Qiniu-Expires=" + c.queryXExpires
		}
		if query != "" {
			copiedReq.RawQuery(query)
		}
		resp := cli.Call(copiedReq)
		fmt.Println("****" + string(resp.Body))
		passed := resp.StatusCode == c.expectedCode
		if passed && c.expectedErr != "" {
			passed = strings.Contains(string(resp.Body), c.expectedErr)
		}
		prefix := "【Failed】"
		if passed {
			prefix = "【Passed】"
		}
		errmsg := ""
		if resp.Err != nil {
			errmsg = resp.Err.Error()
		}
		fmt.Println(prefix + "[" + c.subject + "] responseCode: " + strconv.Itoa(resp.StatusCode) + ", body: " + string(resp.Body) + ", err: " + errmsg)
	}
}
