package _022Q1

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/qianjin/kodo-sample/bucket/bucketcrud"

	"github.com/qianjin/kodo-common/env"

	"github.com/qianjin/kodo-common/authkey"

	"github.com/qianjin/kodo-common/auth"
	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-sample/bucket/bucketconfig"
)

func TestKODO9272_MkBucketV3_AuthExpire_MultiCases_dev(t *testing.T) {
	bucketconfig.SetupEnv("10.200.20.25:10221", "10.200.20.25:10221")
	testKODO9272_MkBucketV3_AuthExpire_MultiCases(t, authkey.Dev_Key_general_storage_011)
}

func TestKODO9272_MkBucketV3_AuthExpire_MultiCases_prod(t *testing.T) {
	bucketconfig.SetupEnv(env.HostDefaultUc, env.HostDefaultUc)
	testKODO9272_MkBucketV3_AuthExpire_MultiCases(t, authkey.Prod_Key_shenqianjin)
}

func testKODO9272_MkBucketV3_AuthExpire_MultiCases(t *testing.T, authKey authkey.AuthKey) {
	mkbucketPath := "/mkbucketv3/%s/region/z0/nodomain/true"
	cli := client.NewManageClientWithHost(bucketconfig.Env.Domain).
		WithKeys(authKey.AK, authKey.SK).
		WithSignType(auth.SignTypeQiniu)

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
		{subject: "???X-Qiniu-Date & ???X-Qiniu-Expires",
			headerXDate: "", headerXExpires: "", expectedCode: 200, expectedErr: ""},
		{subject: "???X-Qiniu-Date (??????:now???????????????15??????) & ?????????X-Qiniu-Expires",
			headerXDate: time10MinsFromNow, headerXExpires: "", expectedCode: 200, expectedErr: ""},
		{subject: "???X-Qiniu-Date (??????:now???????????????15??????) & ?????????X-Qiniu-Expires",
			headerXDate: time10MinsBeforeNow, headerXExpires: "", expectedCode: 200, expectedErr: ""},
		{subject: "???X-Qiniu-Date (??????:now????????????15??????) & ?????????X-Qiniu-Expires",
			headerXDate: time20MinsFromNow, headerXExpires: "", expectedCode: 403, expectedErr: "request time is too skewed"},
		{subject: "???X-Qiniu-Date (??????:now????????????15??????) & ?????????X-Qiniu-Expires",
			headerXDate: time20MinsBeforeNow, headerXExpires: "", expectedCode: 403, expectedErr: "request time is too skewed"},
		{subject: "???X-Qiniu-Date (??????:now????????????15??????) & ??????X-Qiniu-Expires (??????)",
			headerXDate: time20MinsFromNow, headerXExpires: time10MinsFromNow, expectedCode: 200, expectedErr: ""},
		{subject: "???X-Qiniu-Date (??????:now????????????15??????) & ??????X-Qiniu-Expires (??????)",
			headerXDate: time20MinsFromNow, headerXExpires: time10MinsBeforeNow, expectedCode: 403, expectedErr: "token is expired"},
		{subject: "???X-Qiniu-Date (??????:now???????????????15??????) & ??????X-Qiniu-Expires (??????)",
			headerXDate: time10MinsFromNow, headerXExpires: time10MinsFromNow, expectedCode: 200, expectedErr: ""},
		{subject: "???X-Qiniu-Date (??????:now???????????????15??????) & ??????X-Qiniu-Expires (??????)",
			headerXDate: time10MinsFromNow, headerXExpires: time10MinsBeforeNow, expectedCode: 403, expectedErr: "token is expired"},

		{subject: "???X-Qiniu-Date (???????????????) & ?????????X-Qiniu-Expires",
			headerXDate: "not_number", headerXExpires: "", expectedCode: 400, expectedErr: "invalid X-Qiniu-Date"},
		{subject: "???X-Qiniu-Date (???????????????) & ??????X-Qiniu-Expires (??????)",
			headerXDate: "not_number", headerXExpires: time10MinsFromNow, expectedCode: 200, expectedErr: ""},
		{subject: "???X-Qiniu-Date (???????????????) & ??????X-Qiniu-Expires (??????)",
			headerXDate: "not_number", headerXExpires: time10MinsBeforeNow, expectedCode: 403, expectedErr: "token is expired"},
		{subject: "?????????X-Qiniu-Date & ??????X-Qiniu-Expires (???????????????)",
			headerXDate: time10MinsBeforeNow, headerXExpires: "not_number", expectedCode: 400, expectedErr: "invalid X-Qiniu-Expires"},
		{subject: "???X-Qiniu-Date (??????:now???????????????15??????) & ??????X-Qiniu-Expires (???????????????)",
			headerXDate: time10MinsBeforeNow, headerXExpires: "not_number", expectedCode: 400, expectedErr: "invalid X-Qiniu-Expires"},
		{subject: "???X-Qiniu-Date (??????:now????????????15??????) & ??????X-Qiniu-Expires (???????????????)",
			headerXDate: time20MinsBeforeNow, headerXExpires: "not_number", expectedCode: 400, expectedErr: "invalid X-Qiniu-Expires"},
		{subject: "???X-Qiniu-Date (???????????????) & ??????X-Qiniu-Expires (???????????????)",
			headerXDate: "not_number", headerXExpires: "not_number", expectedCode: 400, expectedErr: "invalid X-Qiniu-Expires"},

		{subject: "Header???X-Qiniu-Date (??????:now???????????????15??????) & Query???X-Qiniu-Date (??????:now???????????????15??????)",
			headerXDate: time10MinsFromNow, headerXExpires: "", expectedCode: 200, expectedErr: "", queryXDate: time10MinsFromNow},
		{subject: "Header???X-Qiniu-Date (??????:now???????????????15??????) & Query???X-Qiniu-Date (??????:now????????????15??????)",
			headerXDate: time10MinsFromNow, headerXExpires: "", expectedCode: 403, expectedErr: "request time is too skewed", queryXDate: time20MinsFromNow},
		{subject: "Header???X-Qiniu-Expires (??????) & Query???X-Qiniu-Expires (??????)",
			headerXDate: "", headerXExpires: time10MinsFromNow, expectedCode: 200, expectedErr: "", queryXExpires: time10MinsFromNow},
		{subject: "Header???X-Qiniu-Expires (??????) & Query???X-Qiniu-Expires (??????)",
			headerXDate: "", headerXExpires: time10MinsFromNow, expectedCode: 403, expectedErr: "token is expired", queryXExpires: time10MinsBeforeNow},
		{subject: "Header???X-Qiniu-Date (??????) & Header???X-Qiniu-Expires (??????) & Query???X-Qiniu-Date (??????) & Query???X-Qiniu-Expires (??????)",
			headerXDate: time10MinsFromNow, headerXExpires: time10MinsFromNow, expectedCode: 200, expectedErr: "", queryXDate: time10MinsFromNow, queryXExpires: time10MinsFromNow},
		{subject: "Header???X-Qiniu-Date (??????) & Header???X-Qiniu-Expires (??????) & Query???X-Qiniu-Date (??????) & Query???X-Qiniu-Expires (??????)",
			headerXDate: time10MinsFromNow, headerXExpires: time10MinsFromNow, expectedCode: 403, expectedErr: "token is expired", queryXDate: time10MinsFromNow, queryXExpires: time10MinsBeforeNow},
		{subject: "Header???X-Qiniu-Date (??????) & Header???X-Qiniu-Expires (??????) & Query???X-Qiniu-Date (??????)",
			headerXDate: time10MinsFromNow, headerXExpires: time10MinsBeforeNow, expectedCode: 403, expectedErr: "token is expired", queryXDate: time10MinsFromNow},
		{subject: "Header???X-Qiniu-Date (??????) & Header???X-Qiniu-Expires (??????) & Query???X-Qiniu-Date (??????)",
			headerXDate: time10MinsFromNow, headerXExpires: time10MinsBeforeNow, expectedCode: 403, expectedErr: "token is expired", queryXDate: time20MinsFromNow},
	}

	for _, c := range cases {
		bucket := bucketconfig.GenerateBucketName()
		req := client.NewReq(http.MethodPost, fmt.Sprintf(mkbucketPath, bucket)).
			RawQuery("").
			AddHeader("Host", bucketconfig.Env.Host).
			AddHeader("Content-Type", "application/x-www-form-urlencoded")
		if c.headerXDate != "" {
			req.SetHeader("X-Qiniu-Date", c.headerXDate)
		}
		if c.headerXExpires != "" {
			req.SetHeader("X-Qiniu-Expires", c.headerXExpires)
		}
		// query
		query := req.GetRawQuery()
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
			req.RawQuery(query)
		}
		resp := cli.Call(req)
		if resp.StatusCode == 200 {
			deleteResp := bucketcrud.Delete(cli, bucket)
			assert.Equal(t, http.StatusOK, deleteResp.StatusCode)
		}
		passed := resp.StatusCode == c.expectedCode
		if passed && c.expectedErr != "" {
			passed = strings.Contains(string(resp.Body), c.expectedErr)
		}
		prefix := "???Failed???"
		if passed {
			prefix = "???Passed???"
		}
		errmsg := ""
		if resp.Err != nil {
			errmsg = resp.Err.Error()
		}
		fmt.Println(prefix + "[" + c.subject + "] responseCode: " + strconv.Itoa(resp.StatusCode) + ", body: " + string(resp.Body) + ", err: " + errmsg)
	}
}
