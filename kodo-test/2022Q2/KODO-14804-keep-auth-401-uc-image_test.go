package _022Q2

import (
	"net/http"
	"testing"

	"github.com/qianjin/kodo-common/auth"
	"github.com/qianjin/kodo-common/authkey"
	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-common/env"
	"github.com/qianjin/kodo-sample/bucket-uc/ucconfig"
	"github.com/qianjin/kodo-sample/bucket-uc/uccrud"
	"github.com/qianjin/kodo-sample/bucket-uc/ucmodel"
	"github.com/qianjin/kodo-sample/bucket/bucketconfig"
	"github.com/qianjin/kodo-sample/bucket/bucketcrud"
	"github.com/stretchr/testify/assert"
)

func TestKODO14804_KeepAuth401_UC_Image_dev(t *testing.T) {
	client.DebugMode = true
	bucketconfig.SetupEnv("10.200.20.25:10221", "10.200.20.25:10221")
	ucconfig.SetupEnv("10.200.20.25:10221", "10.200.20.25:10221")
	testKODO14804_KeepAuth401_UC_Image(t, authkey.Dev_Key_IAM_Parent_fansiqiong, authkey.Dev_Key_IAM_Child_shenqianjin_01)
}

func TestKODO14804_KeepAuth401_UC_Image_prod(t *testing.T) {
	bucketconfig.SetupEnv(env.HostDefaultUc, env.HostDefaultUc)
	ucconfig.SetupEnv(env.HostDefaultUc, env.HostDefaultUc)
	testKODO14804_KeepAuth401_UC_Image(t, authkey.Dev_Key_IAM_Parent_fansiqiong, authkey.Dev_Key_IAM_Child_shenqianjin_01)
}

func testKODO14804_KeepAuth401_UC_Image(t *testing.T, bucketAuthKey authkey.AuthKey, authKey authkey.AuthKey) {
	siteUrl := "https://file-examples.com/storage/fef456d9a1627440e9d1c9f/2017/02/file_example_JSON_1kb.json"
	// prepare bucket data
	bucketCli := client.NewManageClientWithHost(bucketconfig.Env.Domain).
		WithKeys(bucketAuthKey.AK, bucketAuthKey.SK).WithSignType(auth.SignTypeQiniu)
	bucket, createBucketResp1 := bucketcrud.Create(bucketCli)
	assert.Equal(t, http.StatusOK, createBucketResp1.StatusCode)
	assert.NotNil(t, bucket)
	defer func() {
		deleteBucketResp := bucketcrud.Delete(bucketCli, bucket)
		assert.Equal(t, http.StatusOK, deleteBucketResp.StatusCode)
	}()

	ucCli := client.NewManageClientWithHost(bucketconfig.Env.Domain).
		WithKeys(authKey.AK, authKey.SK).WithSignType(auth.SignTypeQiniu)
	setImageReq := ucmodel.SetImageReq{Bucket: bucket, SiteURL: siteUrl}
	_, setImageResp := uccrud.SetImage(ucCli, setImageReq)
	assert.Equal(t, http.StatusUnauthorized, setImageResp.StatusCode)
	//assert.Equal(t, http.StatusForbidden, setImageResp.StatusCode)
}
