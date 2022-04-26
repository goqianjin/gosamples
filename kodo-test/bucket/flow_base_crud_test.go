package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/qianjin/kodo-security/kodokey"

	"github.com/qianjin/kodo-common/auth"
	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-test/bucket/bucketconfig"
	"github.com/qianjin/kodo-test/bucket/bucketcrud"
	"github.com/stretchr/testify/assert"
)

func TestCRDFlow(t *testing.T) {
	bucketconfig.SetupEnv("10.200.20.25:10221", "10.200.20.25:10221")

	cli := client.NewClientWithHost(bucketconfig.Env.Domain).
		//WithAuthKey(kodokey.Dev_AK_general_storage_011, kodokey.Dev_SK_general_torage_011).WithSignType(auth.SignTypeQiniu)
		//WithAuthKey(kodokey.Dev_AK_admin, kodokey.Dev_SK_admin).WithSignType(auth.SignTypeQiniuAdmin).WithSuInfo(kodokey.Dev_UID_general_torage_011, 0)
		//WithAuthKey(kodokey.Dev_AK_general_storage_011, kodokey.Dev_SK_general_torage_011).WithSignType(auth.SignTypeQBox)
		WithAuthKey(kodokey.Dev_AK_admin, kodokey.Dev_SK_admin).WithSignType(auth.SignTypeQBoxAdmin).WithSuInfo(kodokey.Dev_UID_general_torage_011, 0)

	// create
	bucket, resp := bucketcrud.Create(cli)
	fmt.Println("created bucket: " + bucket)
	fmt.Printf("result: %+v\n", resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// query
	bucketInfoQ, respQ := bucketcrud.Query(cli, bucket)
	fmt.Printf("result: %+v\n", respQ)
	assert.Equal(t, http.StatusOK, respQ.StatusCode)
	assert.Equal(t, "z0", bucketInfoQ.Region)

	// delete
	respD := bucketcrud.Delete(cli, bucket)
	assert.Equal(t, http.StatusOK, respD.StatusCode)
}
