package up

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/qianjin/kodo-common/auth"
	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-sample/bucket/bucketconfig"
	"github.com/qianjin/kodo-sample/bucket/bucketcrud"
	"github.com/qianjin/kodo-sample/up/upconfig"
	"github.com/qianjin/kodo-sample/up/upcrud_partsv2"
	"github.com/qianjin/kodo-sample/up/upmodel"
	"github.com/qianjin/kodo-security/kodokey"
	"github.com/stretchr/testify/assert"
)

func TestUpPartsV2_BaseFlow_Dev(t *testing.T) {
	upconfig.SetupEnv("10.200.20.23:5010", "10.200.20.23:5010")
	bucketconfig.SetupEnv("10.200.20.25:10221", "10.200.20.25:10221")
	// prepare bucket data
	bucketCli := client.NewManageClientWithHost(bucketconfig.Env.Domain).
		//WithKeys(kodokey.Dev_AK_general_storage_011, kodokey.Dev_SK_general_torage_011).WithSignType(auth.SignTypeQiniu)
		WithKeys(kodokey.Dev_AK_admin, kodokey.Dev_SK_admin).WithSignType(auth.SignTypeQiniuAdmin).WithSuInfo(kodokey.Dev_UID_general_torage_011, 0)
	//WithKeys(kodokey.Dev_AK_general_storage_011, kodokey.Dev_SK_general_torage_011).WithSignType(auth.SignTypeQBox)
	//WithKeys(kodokey.Dev_AK_admin, kodokey.Dev_SK_admin).WithSignType(auth.SignTypeQBoxAdmin).WithSuInfo(kodokey.Dev_UID_general_torage_011, 0)
	bucket, createBucketResp1 := bucketcrud.Create(bucketCli)
	assert.Equal(t, http.StatusOK, createBucketResp1.StatusCode)
	assert.NotNil(t, bucket)
	defer bucketcrud.Delete(bucketCli, bucket)

	// 打开sdk 日志
	client.DebugMode = true
	//client.DeepDebugInfo = true
	key := "test01.txt-" + time.Now().Format("20060102150405")
	putPolicy2 := &auth.PutPolicyV2{
		Scope: bucket,
		//InsertOnly: 1,
		//Exclusive: 1,
		ForceInsertOnly: true,
	}
	upCli := client.NewUpClientWithHost(upconfig.Env.Domain).
		WithAuthKey(kodokey.Dev_AK_general_storage_011, kodokey.Dev_SK_general_torage_011).
		WithPutPolicyV2(putPolicy2)

	// 初始化: 肯定会成功
	initRespBody, initResp := upcrud_partsv2.InitParts(upCli, upmodel.InitPartsReq{Bucket: bucket, Key: key})
	assert.True(t, initResp.Err == nil)
	assert.True(t, &initRespBody != nil)
	assert.True(t, initRespBody.UploadID != "")
	// 分片上传: 传3片,每次上传1M+300b的随机数据
	partSize := 3
	data0 := make([]byte, 1024*1024+300) // 每次上传1M+300b的随机数据(非最后一片)
	// call upload part
	var progresses []upmodel.UploadPartInfo
	for i := 1; i <= partSize; i++ {
		var data []byte
		if i != partSize {
			n, err := rand.Read(data0)
			log.Printf("n = %v, err: %v\n", n, err)
			data = data0
		} else {
			data = []byte(fmt.Sprintf("I am data part %v. The last part!", i))
		}
		dataMd5 := fmt.Sprintf("%x", md5.Sum(data))
		uploadPartsRespBody, uploadPartsResp := upcrud_partsv2.UploadParts(upCli, bytes.NewReader(data), upmodel.UploadPartsReq{
			Bucket: bucket, Key: key,
			UploadId: initRespBody.UploadID, PartNumber: int64(i), Size: len(data),
		})
		assert.True(t, uploadPartsResp.Err == nil)
		assert.True(t, &uploadPartsRespBody != nil)
		assert.True(t, uploadPartsRespBody.MD5 == dataMd5)
		progresses = append(progresses, upmodel.UploadPartInfo{PartNumber: int64(i), Etag: uploadPartsRespBody.Etag})
	}
	// complete part
	completeExtra := &upmodel.RputV2Extra{Progresses: progresses}
	completePartsRespBody, completePartsResp := upcrud_partsv2.CompleteParts(upCli, upmodel.CompletePartsReq{
		Bucket: bucket, Key: key, UploadId: initRespBody.UploadID,
		Extra: completeExtra,
	})
	assert.True(t, completePartsResp.Err == nil)
	assert.True(t, &completePartsRespBody != nil)
	assert.True(t, completePartsRespBody.Key == key)
	// 二次complete报错
	completePartsRespBody, completePartsResp = upcrud_partsv2.CompleteParts(upCli, upmodel.CompletePartsReq{
		Bucket: bucket, Key: key, UploadId: initRespBody.UploadID,
		Extra: completeExtra,
	})
	assert.True(t, completePartsResp.Err != nil)
	assert.True(t, completePartsResp.StatusCode == 612) // ptfdm 返回的code为404, up将404转成了612

	// 初始化: 肯定会成功
	initRespBody, initResp = upcrud_partsv2.InitParts(upCli, upmodel.InitPartsReq{Bucket: bucket, Key: key})
	assert.True(t, initResp.Err == nil)
	assert.True(t, &initRespBody != nil)
	assert.True(t, initRespBody.UploadID != "")
	// 分片上传 一个片(一串字符)
	progresses = make([]upmodel.UploadPartInfo, 0)
	dataStr := "I am the only part data. The last part!"
	data := []byte(dataStr)
	dataMd5 := fmt.Sprintf("%x", md5.Sum(data))
	uploadPartsRespBody, uploadPartsResp := upcrud_partsv2.UploadParts(upCli, bytes.NewReader(data), upmodel.UploadPartsReq{
		Bucket: bucket, Key: key,
		UploadId: initRespBody.UploadID, PartNumber: int64(1), Size: len(data),
	})
	assert.True(t, uploadPartsResp.Err == nil)
	assert.True(t, &uploadPartsRespBody != nil)
	assert.True(t, uploadPartsRespBody.MD5 == dataMd5)
	progresses = append(progresses, upmodel.UploadPartInfo{PartNumber: int64(1), Etag: uploadPartsRespBody.Etag})

	// 同一个key 在force insert only的情况下，不能继续上传了
	completeExtra = &upmodel.RputV2Extra{Progresses: progresses}
	completePartsRespBody, completePartsResp = upcrud_partsv2.CompleteParts(upCli, upmodel.CompletePartsReq{
		Bucket: bucket, Key: key, UploadId: initRespBody.UploadID,
		Extra: completeExtra,
	})
	assert.True(t, completePartsResp.Err != nil)
	assert.True(t, completePartsResp.StatusCode == 614) // file exists
}
