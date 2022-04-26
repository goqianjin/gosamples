package up_direct

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/qianjin/kodo-security/kodokey"

	"github.com/qianjin/kodo-common/auth"

	"github.com/qiniu/go-sdk/v7/client"
	"github.com/qiniu/go-sdk/v7/storage"
)

func TestMimeUp(t *testing.T) {
	// 打开sdk 日志
	client.DebugMode = true
	client.DeepDebugInfo = true

	bucket := "qj-bucket-z0"
	key := "assets/static/20220301/_nuxt/css/license.txt"
	localFile := "/Users/shenqianjin/license.txt"
	// generate token
	/*putPolicy := storage.PutPolicy{
		Scope:      bucket,
		InsertOnly: 1,
	}*/
	putPolicy2 := &auth.PutPolicyV2{
		Scope: bucket,
		//InsertOnly: 1,
		//Exclusive: 1,
		DetectMime: 1,
		MimeType:   "text/javascript",
	}
	//upToken2 := auth.NewUpTokenGenerator(kodokey.Dev_AK_general_storage_011, kodokey.Dev_SK_general_torage_011).WithPutPolicyV2(putPolicy2).GenerateRawToken()
	upToken2 := auth.NewUpTokenGenerator(kodokey.Prod_AK_shenqianjin, kodokey.Prod_SK_shenqianjin).WithPutPolicyV2(putPolicy2).GenerateRawToken()
	// compose uploader
	cfg := storage.Config{
		Zone:          &storage.ZoneHuadong, // 空间对应的机房
		UseHTTPS:      false,                // 是否使用https域名
		UseCdnDomains: false,                // 上传是否使用CDN上传加速
	}
	//构建代理client对象
	//urlParser, _ := url.Parse("http://10.200.20.23:5010")
	tr := http.Transport{
		//Proxy:                 http.ProxyURL(urlParser),
		ResponseHeaderTimeout: 1000 * time.Millisecond,
		Dial: (&net.Dialer{
			Timeout:   3000 * time.Millisecond,
			KeepAlive: 30 * time.Second,
		}).Dial,
	}
	client1 := http.Client{
		Transport: &tr,
	}
	resumeUploader := storage.NewResumeUploaderV2Ex(&cfg, &client.Client{Client: &client1})
	// call upload
	ret := storage.PutRet{}
	putExtra := storage.RputV2Extra{}
	err := resumeUploader.PutFile(context.Background(), &ret, upToken2, key, localFile, &putExtra)
	// validate result
	if err != nil {
		fmt.Println(err)
		return
	}
	//
	fmt.Println(ret.Key, ret.Hash)
	//
}
