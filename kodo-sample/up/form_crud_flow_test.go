package up

import (
	"net/http"
	"testing"
	"time"

	"github.com/qianjin/kodo-common/auth"
	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-sample/bucket/bucketconfig"
	"github.com/qianjin/kodo-sample/bucket/bucketcrud"
	"github.com/qianjin/kodo-sample/up/upconfig"
	"github.com/qianjin/kodo-sample/up/upcrud_form"
	"github.com/qianjin/kodo-sample/up/upmodel"
	"github.com/qianjin/kodo-sample/up/uputil"
	"github.com/qianjin/kodo-security/kodokey"
	"github.com/stretchr/testify/assert"
)

func TestUpForm_BaseFlow_Dev(t *testing.T) {
	upconfig.SetupEnv("10.200.20.23:5010", "10.200.20.23:5010")
	bucketconfig.SetupEnv("10.200.20.25:10221", "10.200.20.25:10221")
	// prepare bucket data
	bucketCli := client.NewManageClientWithHost(bucketconfig.Env.Domain).
		WithKeys(kodokey.Dev_AK_admin, kodokey.Dev_SK_admin).WithSignType(auth.SignTypeQiniuAdmin).WithSuInfo(kodokey.Dev_UID_general_torage_011, 0)
	bucket, createBucketResp1 := bucketcrud.Create(bucketCli)
	assert.Equal(t, http.StatusOK, createBucketResp1.StatusCode)
	assert.NotNil(t, bucket)
	defer bucketcrud.Delete(bucketCli, bucket)

	// 打开sdk 日志
	client.DebugMode = true
	//client.DeepDebugInfo = true
	key := "test01.txt-" + time.Now().Format("20060102150405")
	putPolicyV2 := &auth.PutPolicyV2{Scope: bucket, InsertOnly: 1}
	upCli := client.NewUpClientWithHost(upconfig.Env.Domain).
		WithAuthKey(kodokey.Dev_AK_general_storage_011, kodokey.Dev_SK_general_torage_011).
		WithPutPolicyV2(putPolicyV2)
	// 表单上传
	// 1M随机字符
	size_1M := 1024 * 1024
	token := auth.NewUpTokenGenerator(kodokey.Dev_AK_general_storage_011, kodokey.Dev_SK_general_torage_011).
		WithPutPolicyV2(putPolicyV2).GenerateRawToken()
	uploadRespBody, uploadResp := upcrud_form.FormUpload(upCli, uputil.NewRandomBody(size_1M), upmodel.FormUploadReq{
		Key:         key,
		FileSize:    int64(size_1M),
		UploadToken: token,
	})
	assert.True(t, uploadResp.Err == nil)
	assert.True(t, &uploadRespBody != nil)
	assert.True(t, uploadRespBody.Key == key)
}
