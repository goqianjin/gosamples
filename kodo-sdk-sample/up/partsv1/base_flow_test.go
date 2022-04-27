package partsv1

import (
	"bytes"
	"context"
	"crypto/rand"
	"testing"
	"time"

	"github.com/qianjin/kodo-security/kodokey"

	"github.com/qianjin/kodo-common/auth"
	"github.com/qiniu/go-sdk/v7/client"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/stretchr/testify/assert"
)

func TestUpPartsV1_BaseFlow_Dev(t *testing.T) {
	// 打开sdk 日志
	client.DebugMode = true
	//client.DeepDebugInfo = true
	uphost := "http://10.200.20.23:5010"
	bucket := "file-exist-test-01"
	key := "test01.txt-" + time.Now().Format("20060102150405")
	putPolicy2 := &auth.PutPolicyV2{
		Scope:      bucket,
		InsertOnly: 1,
		//Exclusive: 1,
		//ForceInsertOnly: true,
	}
	upToken2 := auth.NewUpTokenGenerator(kodokey.Dev_AK_general_storage_011, kodokey.Dev_SK_general_torage_011).WithPutPolicyV2(putPolicy2).GenerateRawToken()
	resumeUploader := newResumeV1Uploader(uphost)

	// block1 - 初始化
	// 初始化： 初始化4M的块 & 上传第一个片(1M)
	size_1M := 1024 * 1024
	blockSize := size_1M * 4           // 每次上传1M+300B的随机数据(非最后一片)
	blkputRet1 := &storage.BlkputRet{} // 后面的请求会用到这个对象
	data1 := make([]byte, size_1M)
	n, err := rand.Read(data1)
	assert.True(t, err == nil)
	assert.True(t, n == size_1M)
	mkblkErr := resumeUploader.Mkblk(context.Background(), upToken2, uphost, blkputRet1, blockSize, bytes.NewReader(data1), len(data1))
	assert.True(t, mkblkErr == nil)
	// block1 - 上传片
	// 上传片: 上传第一个块后面的片(3M)
	data3 := make([]byte, size_1M*3)
	n, err = rand.Read(data3)
	assert.True(t, err == nil)
	assert.True(t, n == size_1M*3)
	blkputRetErr := resumeUploader.Bput(context.Background(), upToken2, blkputRet1, bytes.NewReader(data3), len(data3))
	assert.True(t, blkputRetErr == nil)

	// block2 - 初始化
	// 初始化: 初始化4M的块 & 上传整个片(4M)
	blkputRet2 := &storage.BlkputRet{} // 后面的请求会用到这个对象
	data4 := make([]byte, size_1M*4)
	n, err = rand.Read(data4)
	assert.True(t, err == nil)
	assert.True(t, n == size_1M*4)
	mkblkErr = resumeUploader.Mkblk(context.Background(), upToken2, uphost, blkputRet2, blockSize, bytes.NewReader(data4), len(data4))
	assert.True(t, mkblkErr == nil)

	// block3 - 初始化
	// 初始化: 初始化4M的块 & 上传第一个片(1M)
	tailData := []byte("I am The last part data!")
	blkputRet3 := &storage.BlkputRet{} // 后面的请求会用到这个对象
	n, err = rand.Read(data1)
	assert.True(t, err == nil)
	assert.True(t, n == size_1M)
	mkblkErr = resumeUploader.Mkblk(context.Background(), upToken2, uphost, blkputRet3, size_1M*3+len(tailData), bytes.NewReader(data1), len(data1))
	assert.True(t, mkblkErr == nil)
	// block3 - 上传片
	// 上传片: 一共传3片(1M + 1M + 尾部数据大小)
	partSize := 3
	// call upload part
	for i := 1; i <= partSize; i++ {
		//blkputRet := &storage.BlkputRet{}
		var data []byte
		if i != partSize {
			n, err = rand.Read(data1)
			assert.True(t, err == nil)
			assert.True(t, n == size_1M)
			data = data1
		} else {
			data = tailData
		}
		// blkputRet 这个即为前一次请求得到的响应对象
		blkputRetErr = resumeUploader.Bput(context.Background(), upToken2, blkputRet3, bytes.NewReader(data), len(data))
		assert.True(t, blkputRetErr == nil)
	}

	// 生成文件(4M + 4M +3M+尾部大小)
	mkfile_blkputRet := &storage.PutRet{}
	mkfile_extra := &storage.RputExtra{Progresses: []storage.BlkputRet{*blkputRet1, *blkputRet2, *blkputRet3}}
	mkfile_RetErr := resumeUploader.Mkfile(context.Background(), upToken2, uphost, mkfile_blkputRet, key, true, int64(size_1M*11+len(tailData)), mkfile_extra)
	assert.True(t, mkfile_RetErr == nil)
	assert.True(t, mkfile_blkputRet != nil)
	assert.True(t, mkfile_blkputRet.Key == key)

	// 二次生成：Body与第一次一致，幂等成功
	mkfile2_blkputRet := &storage.PutRet{} //&storage.BlkputRet{}
	mkfile2_RetErr := resumeUploader.Mkfile(context.Background(), upToken2, uphost, mkfile2_blkputRet, key, true, int64(size_1M*11+len(tailData)), mkfile_extra)
	assert.True(t, mkfile2_RetErr == nil)
	assert.True(t, mkfile2_blkputRet != nil)
	assert.True(t, mkfile2_blkputRet.Key != "")
	assert.True(t, mkfile2_blkputRet.Key == mkfile_blkputRet.Key)
	assert.True(t, mkfile2_blkputRet.Hash == mkfile_blkputRet.Hash)

	// 二次生成：Body与第一次不一致，期望报错
	mkfile3_blkputRet := &storage.PutRet{} //&storage.BlkputRet{}
	mkfile_extra = &storage.RputExtra{Progresses: []storage.BlkputRet{*blkputRet1, *blkputRet3}}
	mkfile_RetErr = resumeUploader.Mkfile(context.Background(), upToken2, uphost, mkfile3_blkputRet, key, true, int64(size_1M*7+len(tailData)), mkfile_extra)
	assert.True(t, mkfile_RetErr != nil)
	cerr, ok := mkfile_RetErr.(*client.ErrorInfo)
	assert.True(t, ok && cerr.Code == 614) // file exists
	assert.True(t, mkfile3_blkputRet != nil)
	assert.True(t, mkfile3_blkputRet.Key == "")
}
