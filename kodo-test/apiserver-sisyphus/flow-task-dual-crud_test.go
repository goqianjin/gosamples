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
	"github.com/qianjin/kodo-test/bucket/bucketconfig"
	"github.com/qianjin/kodo-test/bucket/bucketcrud"
	"github.com/qianjin/kodo-test/bucket/bucketmodel"
	"github.com/stretchr/testify/assert"
)

func TestTaskDual_CRUDFlow(t *testing.T) {
	sisyphusconfig.SetupEnv("10.200.20.23:12500", "10.200.20.23:12500")
	bucketconfig.SetupEnv("10.200.20.25:10221", "10.200.20.25:10221")
	taskName := sisyphusconfig.GenerateTaskName()
	// prepare bucket data
	bucketCli := client.NewClientWithHost(bucketconfig.Env.Domain).
		WithAuthKey(kodokey.Dev_AK_general_storage_011, kodokey.Dev_SK_general_torage_011).WithSignType(auth.SignTypeQiniu)
	//WithAuthKey(kodokey.Dev_AK_admin, kodokey.Dev_SK_admin).WithSignType(auth.SignTypeQiniuAdmin).WithSuInfo(kodokey.Dev_UID_general_torage_011, 0)
	//WithAuthKey(kodokey.Dev_AK_general_storage_011, kodokey.Dev_SK_general_torage_011).WithSignType(auth.SignTypeQBox)
	//WithAuthKey(kodokey.Dev_AK_admin, kodokey.Dev_SK_admin).WithSignType(auth.SignTypeQBoxAdmin).WithSuInfo(kodokey.Dev_UID_general_torage_011, 0)
	srcBucket, createBucketResp1 := bucketcrud.Create(bucketCli)
	assert.Equal(t, http.StatusOK, createBucketResp1.StatusCode)
	assert.NotNil(t, srcBucket)
	defer bucketcrud.Delete(bucketCli, srcBucket)
	dstBucket, createBucketResp2 := bucketcrud.CreateWithOption(bucketCli, &bucketmodel.CreateOption{Region: "z1"})
	assert.Equal(t, http.StatusOK, createBucketResp2.StatusCode)
	assert.NotNil(t, dstBucket)
	defer bucketcrud.Delete(bucketCli, dstBucket)
	time.Sleep(time.Second)

	cli := client.NewClientWithHost(sisyphusconfig.Env.Domain).
		WithAuthKey(kodokey.Dev_AK_general_storage_011, kodokey.Dev_SK_general_torage_011).WithSignType(auth.SignTypeQiniu)
	//WithAuthKey(kodokey.Dev_AK_admin, kodokey.Dev_SK_admin).WithSignType(auth.SignTypeQiniuAdmin).WithSuInfo(kodokey.Dev_UID_general_torage_011, 0)
	//WithAuthKey(kodokey.Dev_AK_general_storage_011, kodokey.Dev_SK_general_torage_011).WithSignType(auth.SignTypeQBox)
	//WithAuthKey(kodokey.Dev_AK_admin, kodokey.Dev_SK_admin).WithSignType(auth.SignTypeQBoxAdmin).WithSuInfo(kodokey.Dev_UID_general_torage_011, 0)

	// create
	time.Sleep(time.Second)
	createReqBody := sisyphusmodel.CreateDualTaskReq{CreateReq: sisyphusmodel.CreateReq{Name: taskName, Bkts: []string{srcBucket, dstBucket}, IsSync: true}}
	createRespBody, createResp := sisyphuscrud.CreateDualTask(cli, createReqBody)
	fmt.Printf("result: %+v\n", createResp)
	fmt.Printf("result body: %+v\n", createRespBody)
	assert.Equal(t, http.StatusOK, createResp.StatusCode)

	// query
	time.Sleep(time.Second)
	time.Sleep(time.Second)
	queryReqBody := sisyphusmodel.QueryDualTaskReq{TaskId: createRespBody.TaskId}
	queryRespBody, queryResp := sisyphuscrud.QueryDualTask(cli, queryReqBody)
	fmt.Printf("result: %+v\n", queryResp)
	fmt.Printf("result body: %+v\n", queryRespBody)
	assert.Equal(t, http.StatusOK, queryResp.StatusCode)

	// stop
	time.Sleep(time.Second)
	time.Sleep(time.Second)
	stopReqBody := sisyphusmodel.StopDualTaskReq{TaskId: createRespBody.TaskId}
	stopResp := sisyphuscrud.StopDualTask(cli, stopReqBody)
	fmt.Printf("result: %+v\n", stopResp)
	assert.Equal(t, http.StatusOK, stopResp.StatusCode)

	// start
	time.Sleep(time.Second)
	time.Sleep(time.Second)
	startReqBody := sisyphusmodel.StartDualTaskReq{TaskId: createRespBody.TaskId}
	startResp := sisyphuscrud.StartDualTask(cli, startReqBody)
	fmt.Printf("result: %+v\n", startResp)
	assert.Equal(t, http.StatusOK, startResp.StatusCode)

	// delete
	time.Sleep(time.Second)
	time.Sleep(time.Second)
	deleteReqBody := sisyphusmodel.DeleteDualTaskReq{TaskId: createRespBody.TaskId}
	deleteResp := sisyphuscrud.DeleteDualTask(cli, deleteReqBody)
	fmt.Printf("result: %+v\n", deleteResp)
	assert.Equal(t, http.StatusOK, deleteResp.StatusCode)
	time.Sleep(time.Second)
}
