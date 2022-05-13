package up

import (
	"bytes"
	"crypto/rand"
	"net/http"
	"testing"
	"time"

	"github.com/qianjin/kodo-common/auth"
	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-sample/bucket/bucketconfig"
	"github.com/qianjin/kodo-sample/bucket/bucketcrud"
	"github.com/qianjin/kodo-sample/up/upconfig"
	"github.com/qianjin/kodo-sample/up/upcrud_partsv1"
	"github.com/qianjin/kodo-sample/up/upmodel"
	"github.com/qianjin/kodo-security/kodokey"
	"github.com/stretchr/testify/assert"
)

func TestUpPartsV1_BaseFlow_Dev(t *testing.T) {
	upconfig.SetupEnv("10.200.20.23:5010", "10.200.20.23:5010")
	bucketconfig.SetupEnv("10.200.20.25:10221", "10.200.20.25:10221")
	// prepare bucket data
	bucketCli := client.NewClientWithHost(bucketconfig.Env.Domain).
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
		Scope:      bucket,
		InsertOnly: 1,
		//Exclusive: 1,
		//ForceInsertOnly: true,
	}
	upCli := client.NewUpClientWithHost(upconfig.Env.Domain).
		WithAuthKey(kodokey.Dev_AK_general_storage_011, kodokey.Dev_SK_general_torage_011).
		WithPutPolicyV2(putPolicy2)

	// block1 - 初始化
	// 初始化： 初始化4M的块 & 上传第一个片(1M)
	size_1M := 1024 * 1024
	blockSize := size_1M * 4 // 每次上传1M+300B的随机数据(非最后一片)
	//blkputRet1 := &storage.BlkputRet{} // 后面的请求会用到这个对象
	data1M := make([]byte, size_1M)
	n, err := rand.Read(data1M)
	assert.True(t, err == nil)
	assert.True(t, n == size_1M)
	mkblk1RespBody, mkblk1Resp := upcrud_partsv1.Mkblk(upCli, bytes.NewReader(data1M), upmodel.MkblkReq{
		BlockSize:  int64(blockSize),
		BodyLength: int64(size_1M),
	})
	assert.True(t, mkblk1Resp.Err == nil)
	// block1 - 上传片
	// 上传片: 上传第一个块后面的片(3M)
	data3M := make([]byte, size_1M*3)
	n, err = rand.Read(data3M)
	assert.True(t, err == nil)
	assert.True(t, n == size_1M*3)
	bput1RespBody, bput1Resp := upcrud_partsv1.Bput(upCli, bytes.NewReader(data3M), upmodel.BputReq{
		Ctx:    mkblk1RespBody.Ctx,
		Offset: mkblk1RespBody.Offset,
	})
	assert.True(t, bput1Resp.Err == nil)
	assert.True(t, bput1RespBody.Ctx != "")
	mkblk1RespBody.BlkputRet = bput1RespBody.BlkputRet

	// block2 - 初始化
	// 初始化: 初始化4M的块 & 上传整个片(4M)
	//blkputRet2 := &storage.BlkputRet{} // 后面的请求会用到这个对象
	data4M := make([]byte, size_1M*4)
	n, err = rand.Read(data4M)
	assert.True(t, err == nil)
	assert.True(t, n == size_1M*4)
	mkblk2RespBody, mkblk2Resp := upcrud_partsv1.Mkblk(upCli, bytes.NewReader(data4M), upmodel.MkblkReq{
		BlockSize: int64(blockSize),
	})
	assert.True(t, mkblk2Resp.Err == nil)

	// block3 - 初始化
	// 初始化: 初始化4M的块 & 上传第一个片(1M)
	tailData := []byte("I am The last part data!")
	//blkputRet3 := &storage.BlkputRet{} // 后面的请求会用到这个对象
	n, err = rand.Read(data1M)
	assert.True(t, err == nil)
	assert.True(t, n == size_1M)
	mkblk3RespBody, mkblk3Resp := upcrud_partsv1.Mkblk(upCli, bytes.NewReader(data1M), upmodel.MkblkReq{
		BlockSize: int64(size_1M*3 + len(tailData)),
	})
	assert.True(t, mkblk3Resp.Err == nil)
	// block3 - 上传片
	// 上传片: 一共传3片(1M + 1M + 尾部数据大小)
	partSize := 3
	// call upload part
	for i := 1; i <= partSize; i++ {
		//blkputRet := &storage.BlkputRet{}
		var data []byte
		if i != partSize {
			n, err = rand.Read(data1M)
			assert.True(t, err == nil)
			assert.True(t, n == size_1M)
			data = data1M
		} else {
			data = tailData
		}
		// blkputRet 这个即为前一次请求得到的响应对象
		bput3RespBody, bput3Resp := upcrud_partsv1.Bput(upCli, bytes.NewReader(data), upmodel.BputReq{
			Ctx:    mkblk3RespBody.Ctx,
			Offset: mkblk3RespBody.Offset,
		})
		assert.True(t, bput3Resp.Err == nil)
		assert.True(t, bput3RespBody.Ctx != "")
		mkblk3RespBody.BlkputRet = bput3RespBody.BlkputRet
	}

	// 生成文件(4M + 4M +3M+尾部大小)
	mkfile_extra := &upmodel.RputExtra{Progresses: []upmodel.BlkputRet{mkblk1RespBody.BlkputRet, mkblk2RespBody.BlkputRet, mkblk3RespBody.BlkputRet}}
	mkfile1RespBody, mkfile1Resp := upcrud_partsv1.Mkfile(upCli, upmodel.MkfileReq{
		Fsize: int64(size_1M*11 + len(tailData)),
		Key:   key,
		Extra: mkfile_extra,
	})
	assert.True(t, mkfile1Resp.Err == nil)
	assert.True(t, &mkfile1RespBody != nil)
	assert.True(t, mkfile1RespBody.Key == key)

	// 二次生成：Body与第一次一致，幂等成功
	mkfile2RespBody, mkfile2Resp := upcrud_partsv1.Mkfile(upCli, upmodel.MkfileReq{
		Fsize: int64(size_1M*11 + len(tailData)),
		Key:   key,
		Extra: mkfile_extra,
	})
	assert.True(t, mkfile2Resp.Err == nil)
	assert.True(t, &mkfile2RespBody != nil)
	assert.True(t, mkfile2RespBody.Key == "")
	assert.True(t, mkfile2RespBody.Key == mkfile1RespBody.Key)
	assert.True(t, mkfile2RespBody.Hash == mkfile1RespBody.Hash)

	// 二次生成：Body与第一次不一致，期望报错
	mkfile_extra = &upmodel.RputExtra{Progresses: []upmodel.BlkputRet{mkblk1RespBody.BlkputRet, mkblk3RespBody.BlkputRet}}
	mkfile3RespBody, mkfile3Resp := upcrud_partsv1.Mkfile(upCli, upmodel.MkfileReq{
		Fsize: int64(size_1M*7 + len(tailData)),
		Key:   key,
		Extra: mkfile_extra,
	})
	assert.True(t, mkfile3Resp.Err != nil)
	assert.True(t, mkfile3Resp.StatusCode == 614) // file exists
	assert.True(t, &mkfile3RespBody != nil)
	assert.True(t, mkfile3RespBody.Key == "")
}
