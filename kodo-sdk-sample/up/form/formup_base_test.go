package form

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/qianjin/kodo-common/auth"
	"github.com/qianjin/kodo-security/kodokey"
	"github.com/qiniu/go-sdk/v7/client"
	"github.com/qiniu/go-sdk/v7/storage"
)

func TestFormUp(t *testing.T) {
	// 打开sdk 日志
	client.DebugMode = true
	client.DeepDebugInfo = true

	bucket := "qj-test-2022"
	key := ""
	localFile := "/Users/shenqianjin/mydata/data1/test05.txt"
	// generate token
	/*putPolicy := storage.PutPolicy{
		Scope:      bucket,
		InsertOnly: 1,
	}*/
	putPolicy2 := &auth.PutPolicyV2{
		Scope:      bucket,
		InsertOnly: 1,
		Exclusive:  1,
		//ForceInsertOnly: true,
	}

	upToken2 := auth.NewUpTokenGenerator(kodokey.Dev_AK_general_storage_011, kodokey.Dev_SK_general_torage_011).WithPutPolicyV2(putPolicy2).GenerateRawToken()

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
	resumeUploader := storage.NewFormUploaderEx(&cfg, &client.Client{Client: &client1})
	// call upload
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}
	fmt.Println(" ---- " + base64.URLEncoding.EncodeToString([]byte(key)))
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
