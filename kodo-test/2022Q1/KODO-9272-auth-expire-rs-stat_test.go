package _022Q1

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/qianjin/kodo-sample/rs/rsmodel"

	"github.com/qianjin/kodo-sample/rs/rscrud"

	"github.com/qianjin/kodo-sample/up/upcrud_form"
	"github.com/qianjin/kodo-sample/up/upmodel"
	"github.com/qianjin/kodo-sample/up/uputil"

	"github.com/qianjin/kodo-common/auth"
	"github.com/qianjin/kodo-common/authkey"
	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-common/env"
	"github.com/qianjin/kodo-sample/bucket/bucketconfig"
	"github.com/qianjin/kodo-sample/bucket/bucketcrud"
	"github.com/qianjin/kodo-sample/rs/rsconfig"
	"github.com/qianjin/kodo-sample/up/upconfig"
	"github.com/stretchr/testify/assert"
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
	testKODO9272_RsStats_AuthExpire_MultiCases(t, authkey.Prod_Key_kodolog)
}

func testKODO9272_RsStats_AuthExpire_MultiCases(t *testing.T, authKey authkey.AuthKey) {
	// prepare bucket data
	bucketCli := client.NewClientWithHost(bucketconfig.Env.Domain).
		WithKeys(authKey.AK, authKey.SK).WithSignType(auth.SignTypeQiniu)
	bucket, createBucketResp1 := bucketcrud.Create(bucketCli)
	assert.Equal(t, http.StatusOK, createBucketResp1.StatusCode)
	assert.NotNil(t, bucket)
	defer func() {
		deleteBucketResp := bucketcrud.Delete(bucketCli, bucket)
		assert.Equal(t, http.StatusOK, deleteBucketResp.StatusCode)
	}()

	// prepare file
	size1M := 1024 * 1024
	key := "test01.txt-" + time.Now().Format("20060102150405")
	putPolicyV2 := &auth.PutPolicyV2{Scope: bucket, InsertOnly: 1}
	upCli := client.NewUpClientWithHost(upconfig.Env.Domain).
		WithAuthKey(authKey.AK, authKey.SK).
		WithPutPolicyV2(putPolicyV2)
	token := auth.NewUpTokenGenerator(authKey.AK, authKey.SK).
		WithPutPolicyV2(putPolicyV2).GenerateRawToken()
	uploadRespBody, uploadResp := upcrud_form.FormUpload(upCli, uputil.NewRandomBody(size1M), upmodel.FormUploadReq{
		Key:         key,
		FileSize:    int64(size1M),
		UploadToken: token,
	})
	assert.True(t, uploadResp.Err == nil)
	assert.True(t, &uploadRespBody != nil)
	assert.True(t, uploadRespBody.Key == key)
	defer func() {
		rsCli := client.NewClientWithHost(rsconfig.Env.Domain).
			WithKeys(authKey.AK, authKey.SK).WithSignType(auth.SignTypeQiniu)
		_, deleteResp := rscrud.DeleteObject(rsCli, rsmodel.DeleteObjectReq{Bucket: bucket, Key: key})
		assert.Equal(t, http.StatusOK, deleteResp.StatusCode)
	}()

	// rs stats
	fmt.Println("bucket: " + bucket + ", key: " + key)
	path := "/stat/" + base64.URLEncoding.EncodeToString([]byte(bucket+":"+key))
	cli := client.NewClientWithHost(rsconfig.Env.Domain).
		WithKeys(authKey.AK, authKey.SK).
		WithSignType(auth.SignTypeQiniu)
	req := client.NewReq(http.MethodPost, path).
		RawQuery("").
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
		{subject: "含X-Qiniu-Date (格式不合法) & 包含X-Qiniu-Expires (格式不合法)",
			headerXDate: "not_number", headerXExpires: "not_number", expectedCode: 400, expectedErr: "invalid X-Qiniu-Expires"},

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
