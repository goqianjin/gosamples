package main

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/qianjin/kodo-security/kodokey"

	"github.com/qianjin/kodo-test/bucket/bucketconfig"
	"github.com/qianjin/kodo-test/bucket/bucketcrud"

	"github.com/qianjin/kodo-common/auth"
	"github.com/qianjin/kodo-common/client"
	"github.com/stretchr/testify/assert"
)

func TestDropBucketsByPrefix_Dev(t *testing.T) {
	bucketconfig.SetupEnv("10.200.20.25:10221", "10.200.20.25:10221")
	dropPrefix := "qianjin-bucket-2022"

	cli := client.NewClientWithHost(bucketconfig.Env.Domain).
		WithAuthKey(kodokey.Dev_AK_general_storage_011, kodokey.Dev_SK_general_torage_011).
		WithSignType(auth.SignTypeQiniu)

	// list
	buckets, resp := bucketcrud.List(cli)
	fmt.Printf("ret: %+v\n", buckets)
	fmt.Printf("result: %+v\n", resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	// delete by loop
	for _, bucket := range buckets {
		if !strings.HasPrefix(bucket, dropPrefix) {
			continue
		}
		respD := bucketcrud.Delete(cli, bucket)
		assert.Equal(t, http.StatusOK, respD.StatusCode)
		fmt.Printf("deleted bucket: %v\n", bucket)
	}

}
