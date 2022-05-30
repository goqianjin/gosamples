package _022Q2

import (
	"net/http"
	"testing"
	"time"

	"github.com/qianjin/kodo-common/proxyuser"

	"github.com/qianjin/kodo-common/auth"
	"github.com/qianjin/kodo-common/authkey"
	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-common/env"
	"github.com/qianjin/kodo-sample/bucket/bucketconfig"
	"github.com/qianjin/kodo-sample/bucket/bucketcrud"
	"github.com/qianjin/kodo-sample/one/oneconfig"
	"github.com/qianjin/kodo-sample/one/onecrud"
	"github.com/qianjin/kodo-sample/one/onemodel"
	"github.com/stretchr/testify/assert"
)

func TestKODO15123_TuneSwtichesDisableAuthCheckExpire_BucketDomain_ListV3_dev(t *testing.T) {
	//client.DebugMode = true
	bucketconfig.SetupEnv("10.200.20.25:10221", "10.200.20.25:10221")
	oneconfig.SetupEnv("10.200.20.25:23200", "10.200.20.25:23200")
	testKODO15123_TuneSwtichesDisableAuthCheckExpire_BucketDomain_ListV3(t, authkey.Dev_Key_general_storage_011, proxyuser.ProxyUser_Dev_general_storage_011)
}

func TestKODO15123_TuneSwtichesDisableAuthCheckExpire_BucketDomain_ListV3_prod(t *testing.T) {
	oneconfig.SetupEnv(env.HostDefaultUc, env.HostDefaultUc)
	testKODO15123_TuneSwtichesDisableAuthCheckExpire_BucketDomain_ListV3(t, authkey.Prod_Key_kodolog, proxyuser.ProxyUser_Dev_general_storage_011)
}

func testKODO15123_TuneSwtichesDisableAuthCheckExpire_BucketDomain_ListV3(t *testing.T, authKey authkey.AuthKey, user proxyuser.ProxyUser) {
	// rs client
	oneCli := client.NewProxyClientWithHost(oneconfig.Env.Domain).
		WithProxyUser(user)

	// prepare bucket data
	bucketCli := client.NewManageClientWithHost(bucketconfig.Env.Domain).
		WithKeys(authKey.AK, authKey.SK).WithSignType(auth.SignTypeQiniu)
	bucket, createBucketResp1 := bucketcrud.Create(bucketCli)
	assert.Equal(t, http.StatusOK, createBucketResp1.StatusCode)
	assert.NotNil(t, bucket)
	defer func() {
		deleteBucketResp := bucketcrud.Delete(bucketCli, bucket)
		assert.Equal(t, http.StatusOK, deleteBucketResp.StatusCode)
	}()

	// 恢复正常环境：关闭 禁用过期时间校验功能
	_, putUserTuneSwitchesResp := onecrud.PutUserTuneSwitches(oneCli, onemodel.PutUserTuneSwitchesReq(""))
	assert.True(t, putUserTuneSwitchesResp.Err == nil)
	assert.True(t, putUserTuneSwitchesResp.StatusCode == http.StatusOK)

	var RFC3339TimeInSecondPattern = "20060102T150405Z"
	time10MinsFromNow := time.Now().Add(10 * time.Minute).UTC().Format(RFC3339TimeInSecondPattern)
	time20MinsFromNow := time.Now().Add(20 * time.Minute).UTC().Format(RFC3339TimeInSecondPattern)
	// 默认 禁用过期时间校验功能
	// 正常请求返回 200
	_, queryBucketResp := bucketcrud.ListDomainsV3(bucketCli, bucket, client.WithReqHeader(map[string]string{"X-Qiniu-Date": time10MinsFromNow}))
	assert.True(t, queryBucketResp.Err == nil)
	assert.True(t, queryBucketResp.StatusCode == http.StatusOK)
	// 正常请求 时间过期 报错403
	_, queryBucketResp = bucketcrud.ListDomainsV3(bucketCli, bucket, client.WithReqHeader(map[string]string{"X-Qiniu-Date": time20MinsFromNow}))
	assert.True(t, queryBucketResp.Err != nil)
	assert.True(t, queryBucketResp.StatusCode == http.StatusForbidden)
	// 开启 禁用过期时间校验功能时：
	_, putUserTuneSwitchesResp = onecrud.PutUserTuneSwitches(oneCli, onemodel.PutUserTuneSwitchesReq("000000001"))
	assert.True(t, putUserTuneSwitchesResp.Err == nil)
	assert.True(t, putUserTuneSwitchesResp.StatusCode == http.StatusOK)
	time.Sleep(3 * time.Second) // bucket服务本地缓存过期时间 3s, redis 过期时间 50ms
	// 正常请求 时间过期 不报错 返回200
	_, queryBucketResp = bucketcrud.ListDomainsV3(bucketCli, bucket, client.WithReqHeader(map[string]string{"X-Qiniu-Date": time20MinsFromNow}))
	assert.True(t, queryBucketResp.Err == nil)
	assert.True(t, queryBucketResp.StatusCode == http.StatusOK)
	// 关闭 禁用过期时间校验功能
	_, putUserTuneSwitchesResp = onecrud.PutUserTuneSwitches(oneCli, onemodel.PutUserTuneSwitchesReq("000000000"))
	assert.True(t, putUserTuneSwitchesResp.Err == nil)
	assert.True(t, putUserTuneSwitchesResp.StatusCode == http.StatusOK)
	time.Sleep(3 * time.Second) // bucket服务本地缓存过期时间 3s, redis 过期时间 50ms
	// 正常请求 时间过期 报错403
	_, queryBucketResp = bucketcrud.ListDomainsV3(bucketCli, bucket, client.WithReqHeader(map[string]string{"X-Qiniu-Date": time20MinsFromNow}))
	assert.True(t, queryBucketResp.Err != nil)
	assert.True(t, queryBucketResp.StatusCode == http.StatusForbidden)

	// 开启 禁用过期时间校验功能
	_, putUserTuneSwitchesResp = onecrud.PutUserTuneSwitches(oneCli, onemodel.PutUserTuneSwitchesReq("000000001"))
	assert.True(t, putUserTuneSwitchesResp.Err == nil)
	assert.True(t, putUserTuneSwitchesResp.StatusCode == http.StatusOK)
	time.Sleep(3 * time.Second) // bucket服务本地缓存过期时间 3s, redis 过期时间 50ms
	// 正常请求 时间过期 不报错 返回200
	_, queryBucketResp = bucketcrud.ListDomainsV3(bucketCli, bucket, client.WithReqHeader(map[string]string{"X-Qiniu-Date": time20MinsFromNow}))
	assert.True(t, queryBucketResp.Err == nil)
	assert.True(t, queryBucketResp.StatusCode == http.StatusOK)
	// 关闭 禁用过期时间校验功能
	_, putUserTuneSwitchesResp = onecrud.PutUserTuneSwitches(oneCli, onemodel.PutUserTuneSwitchesReq(""))
	assert.True(t, putUserTuneSwitchesResp.Err == nil)
	assert.True(t, putUserTuneSwitchesResp.StatusCode == http.StatusOK)
	time.Sleep(3 * time.Second) // bucket服务本地缓存过期时间 3s, redis 过期时间 50ms
	// 正常请求 时间过期 报错403
	_, queryBucketResp = bucketcrud.ListDomainsV3(bucketCli, bucket, client.WithReqHeader(map[string]string{"X-Qiniu-Date": time20MinsFromNow}))
	assert.True(t, queryBucketResp.Err != nil)
	assert.True(t, queryBucketResp.StatusCode == http.StatusForbidden)
}
