package partsv2

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/qianjin/kodo-security/kodokey"

	"github.com/qianjin/kodo-common/auth"
	"github.com/qiniu/go-sdk/v7/client"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/stretchr/testify/assert"
)

func TestUpPartsV2_BaseFlow_Dev(t *testing.T) {
	// 打开sdk 日志
	client.DebugMode = true
	//client.DeepDebugInfo = true
	uphost := "http://10.200.20.23:5010"
	bucket := "file-exist-test-01"
	key := "test01.txt-" + time.Now().Format("20060102150405")
	putPolicy2 := &auth.PutPolicyV2{
		Scope: bucket,
		//InsertOnly: 1,
		//Exclusive: 1,
		ForceInsertOnly: true,
	}
	upToken2 := auth.NewUpTokenGenerator(kodokey.Dev_AK_general_storage_011, kodokey.Dev_SK_general_torage_011).WithPutPolicyV2(putPolicy2).GenerateRawToken()
	resumeUploader := newResumeV2Uploader(uphost)

	// 初始化: 肯定会成功
	initRet := &storage.InitPartsRet{}
	initErr := resumeUploader.InitParts(context.Background(), upToken2, uphost, bucket, key, true, initRet)
	assert.True(t, initErr == nil)
	assert.True(t, initRet != nil)
	assert.True(t, initRet.UploadID != "")
	// 分片上传: 传3片,每次上传1M+300b的随机数据
	partSize := 3
	data0 := make([]byte, 1024*1024+300) // 每次上传1M+300b的随机数据(非最后一片)
	// call upload part
	var progresses []storage.UploadPartInfo
	for i := 1; i <= partSize; i++ {
		uploadPartRet := &storage.UploadPartsRet{}
		var data []byte
		if i != partSize {
			n, err := rand.Read(data0)
			log.Printf("n = %v, err: %v\n", n, err)
			data = data0
		} else {
			data = []byte(fmt.Sprintf("I am data part %v. The last part!", i))
		}
		dataMd5 := fmt.Sprintf("%x", md5.Sum(data))
		uploadPartErr := resumeUploader.UploadParts(context.Background(), upToken2, uphost, bucket, key, true,
			initRet.UploadID, int64(i), dataMd5, uploadPartRet, bytes.NewReader(data), len(data))
		assert.True(t, uploadPartErr == nil)
		assert.True(t, uploadPartRet != nil)
		assert.True(t, uploadPartRet.MD5 == dataMd5)
		progresses = append(progresses, storage.UploadPartInfo{PartNumber: int64(i), Etag: uploadPartRet.Etag})
	}
	// complete part
	completePartRet := &storage.PutRet{}
	completeExtra := &storage.RputV2Extra{Progresses: progresses}
	completePartsErr := resumeUploader.CompleteParts(context.Background(), upToken2, uphost, completePartRet, bucket,
		key, true, initRet.UploadID, completeExtra)
	assert.True(t, completePartsErr == nil)
	assert.True(t, completePartRet != nil)
	assert.True(t, completePartRet.Key == key)
	// 二次complete报错
	completePartsErr = resumeUploader.CompleteParts(context.Background(), upToken2, uphost, completePartRet, bucket,
		key, true, initRet.UploadID, completeExtra)
	assert.True(t, completePartsErr != nil)
	cerr, ok := completePartsErr.(*client.ErrorInfo)
	assert.True(t, ok && cerr.Code == 612) // ptfdm 返回的code为404, up将404转成了612

	// 初始化: 肯定会成功
	initRet = &storage.InitPartsRet{}
	initErr = resumeUploader.InitParts(context.Background(), upToken2, uphost, bucket, key, true, initRet)
	assert.True(t, initErr == nil)
	assert.True(t, initRet != nil)
	assert.True(t, initRet.UploadID != "")
	// 分片上传 一个片(一串字符)
	progresses = make([]storage.UploadPartInfo, 0)
	uploadPartRet := &storage.UploadPartsRet{}
	dataStr := "I am the only part data. The last part!"
	data := []byte(dataStr)
	dataMd5 := fmt.Sprintf("%x", md5.Sum(data))
	uploadPartErr := resumeUploader.UploadParts(context.Background(), upToken2, uphost, bucket, key, true,
		initRet.UploadID, int64(1), dataMd5, uploadPartRet, strings.NewReader(dataStr), len(data))
	assert.True(t, uploadPartErr == nil)
	assert.True(t, uploadPartRet != nil)
	assert.True(t, uploadPartRet.MD5 == dataMd5)
	progresses = append(progresses, storage.UploadPartInfo{PartNumber: int64(1), Etag: uploadPartRet.Etag})
	// 同一个key 在force insert only的情况下，不能继续上传了
	completePartRet = &storage.PutRet{}
	completeExtra = &storage.RputV2Extra{Progresses: progresses}
	completePartsErr = resumeUploader.CompleteParts(context.Background(), upToken2, uphost, completePartRet, bucket,
		key, true, initRet.UploadID, completeExtra)
	assert.True(t, completePartsErr != nil)
	cerr, ok = completePartsErr.(*client.ErrorInfo)
	assert.True(t, ok && cerr.Code == 614) // file exists
}
