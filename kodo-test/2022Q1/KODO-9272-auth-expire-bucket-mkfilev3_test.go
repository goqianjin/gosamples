package _022Q1

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/qianjin/kodo-common/env"

	"github.com/qianjin/kodo-common/authkey"

	"github.com/qianjin/kodo-common/auth"
	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-sample/bucket/bucketconfig"
)

func TestKODO9272_MkBucketV3_AuthExpire_MultiCases_dev(t *testing.T) {
	bucketconfig.SetupEnv("10.200.20.25:10221", "10.200.20.25:10221")
	testKODO9272_MkBucketV3_AuthExpire_MultiCases(t, bucketconfig.Env, authkey.Dev_Key_general_storage_011)

}

func TestKODO9272_MkBucketV3_AuthExpire_MultiCases_prod(t *testing.T) {
	bucketconfig.SetupEnv(env.DefaultUcHost, env.DefaultUcHost)
	testKODO9272_MkBucketV3_AuthExpire_MultiCases(t, bucketconfig.Env, authkey.Prod_Key_shenqianjin)

}
func testKODO9272_MkBucketV3_AuthExpire_MultiCases(t *testing.T, env env.EnvInfo, key authkey.AuthKey) {
	RFC3339Nano := "20060102150405.999999999"
	mkbucketPath := "/mkbucketv3/qianjin-%s/region/z0/nodomain/true"
	body := ""
	cli := client.NewClientWithHost(env.Domain).
		WithKeys(key.AK, key.SK).
		WithSignType(auth.SignTypeQiniu)
	req := client.NewReq(http.MethodPost, "").
		RawQuery("").
		AddHeader("Host", env.Host).
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		BodyStr(body)
	resp := cli.Call(req)
	fmt.Printf("Resp: %+v\n", resp)

	type UserCase struct {
		subject        string
		headerXDate    string
		headerXExpires string
		expectedCode   int
		expectedErr    string
		queryXDate     string
		queryXExpires  string
	}
	/*time10MinsFromNow := strconv.FormatInt(time.Now().Add(10*time.Minute).Unix(), 10)
	time20MinsFromNow := strconv.FormatInt(time.Now().Add(20*time.Minute).Unix(), 10)
	time10MinsBeforeNow := strconv.FormatInt(time.Now().Add(-10*time.Minute).Unix(), 10)
	time20MinsBeforeNow := strconv.FormatInt(time.Now().Add(-20*time.Minute).Unix(), 10)*/
	var RFC3339TimeInSecondPattern = "20060102T150405Z"
	time10MinsFromNow := time.Now().Add(10 * time.Minute).UTC().Format(RFC3339TimeInSecondPattern)
	time20MinsFromNow := time.Now().Add(20 * time.Minute).UTC().Format(RFC3339TimeInSecondPattern)
	time10MinsBeforeNow := time.Now().Add(-10 * time.Minute).UTC().Format(RFC3339TimeInSecondPattern)
	time20MinsBeforeNow := time.Now().Add(-20 * time.Minute).UTC().Format(RFC3339TimeInSecondPattern)
	cases := []UserCase{
		{subject: "无X-Qiniu-Date & 无X-Qiniu-Expires",
			headerXDate: "", headerXExpires: "", expectedCode: 200, expectedErr: ""},
		{subject: "含X-Qiniu-Date (有效:now往后不超过15分钟) & 不包含X-Qiniu-Expires",
			headerXDate: time10MinsFromNow, headerXExpires: "", expectedCode: 200, expectedErr: ""},
		{subject: "含X-Qiniu-Date (有效:now往前不超过15分钟) & 不包含X-Qiniu-Expires",
			headerXDate: time10MinsBeforeNow, headerXExpires: "", expectedCode: 200, expectedErr: ""},
		{subject: "含X-Qiniu-Date (无效:now往后超过15分钟) & 不包含X-Qiniu-Expires",
			headerXDate: time20MinsFromNow, headerXExpires: "", expectedCode: 403, expectedErr: "request time is too skewed"},
		{subject: "含X-Qiniu-Date (无效:now往前超过15分钟) & 不包含X-Qiniu-Expires",
			headerXDate: time20MinsBeforeNow, headerXExpires: "", expectedCode: 403, expectedErr: "request time is too skewed"},
		{subject: "含X-Qiniu-Date (无效:now往后超过15分钟) & 包含X-Qiniu-Expires (有效)",
			headerXDate: time20MinsFromNow, headerXExpires: time10MinsFromNow, expectedCode: 200, expectedErr: ""},
		{subject: "含X-Qiniu-Date (无效:now往后超过15分钟) & 包含X-Qiniu-Expires (无效)",
			headerXDate: time20MinsFromNow, headerXExpires: time10MinsBeforeNow, expectedCode: 403, expectedErr: "token is expired"},
		{subject: "含X-Qiniu-Date (有效:now往后不超过15分钟) & 包含X-Qiniu-Expires (有效)",
			headerXDate: time10MinsFromNow, headerXExpires: time10MinsFromNow, expectedCode: 200, expectedErr: ""},
		{subject: "含X-Qiniu-Date (有效:now往后不超过15分钟) & 包含X-Qiniu-Expires (无效)",
			headerXDate: time10MinsFromNow, headerXExpires: time10MinsBeforeNow, expectedCode: 403, expectedErr: "token is expired"},

		{subject: "含X-Qiniu-Date (格式不合法) & 不包含X-Qiniu-Expires",
			headerXDate: "not_number", headerXExpires: "", expectedCode: 400, expectedErr: "invalid X-Qiniu-Date"},
		{subject: "含X-Qiniu-Date (格式不合法) & 包含X-Qiniu-Expires (有效)",
			headerXDate: "not_number", headerXExpires: time10MinsFromNow, expectedCode: 200, expectedErr: ""},
		{subject: "含X-Qiniu-Date (格式不合法) & 包含X-Qiniu-Expires (无效)",
			headerXDate: "not_number", headerXExpires: time10MinsBeforeNow, expectedCode: 403, expectedErr: "token is expired"},
		{subject: "不包含X-Qiniu-Date & 包含X-Qiniu-Expires (格式不合法)",
			headerXDate: time10MinsBeforeNow, headerXExpires: "not_number", expectedCode: 400, expectedErr: "invalid X-Qiniu-Expires"},
		{subject: "含X-Qiniu-Date (有效:now往后不超过15分钟) & 包含X-Qiniu-Expires (格式不合法)",
			headerXDate: time10MinsBeforeNow, headerXExpires: "not_number", expectedCode: 400, expectedErr: "invalid X-Qiniu-Expires"},
		{subject: "含X-Qiniu-Date (无效:now往后超过15分钟) & 包含X-Qiniu-Expires (格式不合法)",
			headerXDate: time20MinsBeforeNow, headerXExpires: "not_number", expectedCode: 400, expectedErr: "invalid X-Qiniu-Expires"},
		{subject: "含X-Qiniu-Date (格式不合法) & 包含X-Qiniu-Expires (格式不合法)",
			headerXDate: "not_number", headerXExpires: "not_number", expectedCode: 400, expectedErr: "invalid X-Qiniu-Expires"},

		{subject: "Header含X-Qiniu-Date (有效:now往后不超过15分钟) & Query含X-Qiniu-Date (有效:now往后不超过15分钟)",
			headerXDate: time10MinsFromNow, headerXExpires: "", expectedCode: 200, expectedErr: "", queryXDate: time10MinsFromNow},
		{subject: "Header含X-Qiniu-Date (有效:now往后不超过15分钟) & Query含X-Qiniu-Date (无效:now往后超过15分钟)",
			headerXDate: time10MinsFromNow, headerXExpires: "", expectedCode: 403, expectedErr: "request time is too skewed", queryXDate: time20MinsFromNow},
		{subject: "Header含X-Qiniu-Expires (有效) & Query含X-Qiniu-Expires (有效)",
			headerXDate: "", headerXExpires: time10MinsFromNow, expectedCode: 200, expectedErr: "", queryXExpires: time10MinsFromNow},
		{subject: "Header含X-Qiniu-Expires (有效) & Query含X-Qiniu-Expires (无效)",
			headerXDate: "", headerXExpires: time10MinsFromNow, expectedCode: 403, expectedErr: "token is expired", queryXExpires: time10MinsBeforeNow},
		{subject: "Header含X-Qiniu-Date (有效) & Header含X-Qiniu-Expires (有效) & Query含X-Qiniu-Date (有效) & Query含X-Qiniu-Expires (有效)",
			headerXDate: time10MinsFromNow, headerXExpires: time10MinsFromNow, expectedCode: 200, expectedErr: "", queryXDate: time10MinsFromNow, queryXExpires: time10MinsFromNow},
		{subject: "Header含X-Qiniu-Date (有效) & Header含X-Qiniu-Expires (有效) & Query含X-Qiniu-Date (有效) & Query含X-Qiniu-Expires (无效)",
			headerXDate: time10MinsFromNow, headerXExpires: time10MinsFromNow, expectedCode: 403, expectedErr: "token is expired", queryXDate: time10MinsFromNow, queryXExpires: time10MinsBeforeNow},
		{subject: "Header含X-Qiniu-Date (有效) & Header含X-Qiniu-Expires (无效) & Query含X-Qiniu-Date (有效)",
			headerXDate: time10MinsFromNow, headerXExpires: time10MinsBeforeNow, expectedCode: 403, expectedErr: "token is expired", queryXDate: time10MinsFromNow},
		{subject: "Header含X-Qiniu-Date (有效) & Header含X-Qiniu-Expires (无效) & Query含X-Qiniu-Date (无效)",
			headerXDate: time10MinsFromNow, headerXExpires: time10MinsBeforeNow, expectedCode: 403, expectedErr: "token is expired", queryXDate: time20MinsFromNow},
	}

	for _, c := range cases {
		copiedReq := req.DeepClone()
		newPath := fmt.Sprintf(mkbucketPath, strings.ReplaceAll(time.Now().Format(RFC3339Nano), ".", ""))
		copiedReq.Path(newPath)
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
		passed := resp.StatusCode == c.expectedCode
		if c.expectedErr != "" {
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
		fmt.Println(prefix + "[" + c.subject + "] responseCode: " + strconv.Itoa(resp.StatusCode) + ", body: " + fmt.Sprintf("%v", resp.Body) + ", err: " + errmsg)
	}
}
