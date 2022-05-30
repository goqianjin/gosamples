package main

import (
	"fmt"
	"net/http"
	"sisyphus/sisyphusconfig"
	"sisyphus/sisyphuscrud"
	"sisyphus/sisyphusmodel"
	"testing"
	"time"

	"github.com/qianjin/kodo-security/kodokey"

	"github.com/qianjin/kodo-common/auth"
	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-sample/bucket/bucketconfig"
	"github.com/qianjin/kodo-sample/bucket/bucketcrud"
	"github.com/qianjin/kodo-sample/bucket/bucketmodel"
	"github.com/stretchr/testify/assert"
)

func TestCRUDFlow(t *testing.T) {
	sisyphusconfig.SetupEnv("10.200.20.23:12500", "10.200.20.23:12500")
	bucketconfig.SetupEnv("10.200.20.25:10221", "10.200.20.25:10221")
	taskName := sisyphusconfig.GenerateTaskName()
	// prepare bucket data
	bucketCli := client.NewManageClientWithHost(bucketconfig.Env.Domain).
		//WithKeys(kodokey.Dev_AK_general_storage_011, kodokey.Dev_SK_general_torage_011).WithSignType(auth.SignTypeQiniu)
		WithKeys(kodokey.Dev_AK_admin, kodokey.Dev_SK_admin).WithSignType(auth.SignTypeQiniuAdmin).WithSuInfo(kodokey.Dev_UID_general_torage_011, 0)
	//WithKeys(kodokey.Dev_AK_general_storage_011, kodokey.Dev_SK_general_torage_011).WithSignType(auth.SignTypeQBox)
	//WithKeys(kodokey.Dev_AK_admin, kodokey.Dev_SK_admin).WithSignType(auth.SignTypeQBoxAdmin).WithSuInfo(kodokey.Dev_UID_general_torage_011, 0)
	srcBucket, createBucketResp1 := bucketcrud.Create(bucketCli)
	assert.Equal(t, http.StatusOK, createBucketResp1.StatusCode)
	assert.NotNil(t, srcBucket)
	defer bucketcrud.Delete(bucketCli, srcBucket)
	dstBucket, createBucketResp2 := bucketcrud.CreateWithOption(bucketCli, &bucketmodel.CreateOption{Region: "z1"})
	assert.Equal(t, http.StatusOK, createBucketResp2.StatusCode)
	assert.NotNil(t, dstBucket)
	defer bucketcrud.Delete(bucketCli, dstBucket)
	time.Sleep(time.Second)
	time.Sleep(time.Second)

	cli := client.NewManageClientWithHost(sisyphusconfig.Env.Domain).
		//WithKeys(kodokey.Dev_AK_general_storage_011, kodokey.Dev_SK_general_torage_011).WithSignType(auth.SignTypeQiniu)
		//WithKeys(kodokey.Dev_AK_admin, kodokey.Dev_SK_admin).WithSignType(auth.SignTypeQiniuAdmin).WithSuInfo(kodokey.Dev_UID_general_torage_011, 0)
		//WithKeys(kodokey.Dev_AK_general_storage_011, kodokey.Dev_SK_general_torage_011).WithSignType(auth.SignTypeQBox)
		WithKeys(kodokey.Dev_AK_admin, kodokey.Dev_SK_admin).WithSignType(auth.SignTypeQBoxAdmin).WithSuInfo(kodokey.Dev_UID_general_torage_011, 0)

	// create
	createReqBody := sisyphusmodel.CreateReq{Name: taskName, SrcBkt: srcBucket, DstBkt: dstBucket, IsSync: true}
	createRespBody, createResp := sisyphuscrud.Create(cli, createReqBody)
	fmt.Printf("result: %+v\n", createResp)
	fmt.Printf("result body: %+v\n", createRespBody)
	assert.Equal(t, http.StatusOK, createResp.StatusCode)
	time.Sleep(time.Second)

	// query
	queryReqBody := sisyphusmodel.QueryReq{TaskId: createRespBody.TaskId}
	queryRespBody, queryResp := sisyphuscrud.Query(cli, queryReqBody)
	fmt.Printf("result: %+v\n", queryResp)
	fmt.Printf("result body: %+v\n", queryRespBody)
	assert.Equal(t, http.StatusOK, queryResp.StatusCode)

	// stop
	stopReqBody := sisyphusmodel.StopReq{TaskId: createRespBody.TaskId}
	stopResp := sisyphuscrud.Stop(cli, stopReqBody)
	fmt.Printf("result: %+v\n", stopResp)
	assert.Equal(t, http.StatusOK, stopResp.StatusCode)
	time.Sleep(time.Second)

	// start
	startReqBody := sisyphusmodel.StartReq{TaskId: createRespBody.TaskId}
	startResp := sisyphuscrud.Start(cli, startReqBody)
	fmt.Printf("result: %+v\n", startResp)
	assert.Equal(t, http.StatusOK, startResp.StatusCode)
	time.Sleep(time.Second)

	// delete
	deleteReqBody := sisyphusmodel.DeleteReq{TaskId: createRespBody.TaskId}
	deleteResp := sisyphuscrud.Delete(cli, deleteReqBody)
	fmt.Printf("result: %+v\n", deleteResp)
	assert.Equal(t, http.StatusOK, deleteResp.StatusCode)
	time.Sleep(time.Second)
}
