package main

import (
	"fmt"

	"github.com/qianjin/kodo-security/kodokey"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

func main() {
	accessKey := kodokey.Prod_AK_shenqianjin
	secretKey := kodokey.Prod_SK_shenqianjin
	mac := qbox.NewMac(accessKey, secretKey)
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: true,
	}
	// 指定空间所在的区域，如果不指定将自动探测
	// 如果没有特殊需求，默认不需要指定
	//cfg.Zone=&storage.ZoneHuabei
	bucketManager := storage.NewBucketManager(mac, &cfg)

	bucket := "qj-test-fetch-tar"
	resURL := "http://r19tfe1y3.hn-bkt.clouddn.com/test_up2.json"

	// 指定保存的key
	//fetch(bucketManager, resURL, bucket, "test_up.json")
	// 不指定保存的key，默认用文件hash作为文件名
	//fetchWithoutKey(bucketManager, resURL, bucket)
	// 异步拉取
	fetchAsync(bucketManager, resURL, bucket, "test_up2.json")

}

func fetchAsync(bucketManager *storage.BucketManager, resURL string, bucket string, key string) {
	param := storage.AsyncFetchParam{
		Url:         resURL,
		Bucket:      bucket,
		Key:         key,
		CallbackURL: "http://185.12.311.12",
		//CallbackBody: "{}",
	}
	fetchRet, err := bucketManager.AsyncFetch(param)
	if err != nil {
		fmt.Println("fetch error,", err)
	} else {
		fmt.Printf("return of async fetch: %+n", fetchRet)
	}

}

func fetch(bucketManager *storage.BucketManager, resURL string, bucket string, key string) {
	// 指定保存的key
	//bucketManager.AsyncFetch()
	fetchRet, err := bucketManager.Fetch(resURL, bucket, "test_up.json")
	if err != nil {
		fmt.Println("fetch error,", err)
	} else {
		fmt.Println(fetchRet.String())
	}
}

func fetchWithoutKey(bucketManager *storage.BucketManager, resURL string, bucket string) {
	fetchRet, err := bucketManager.FetchWithoutKey(resURL, bucket)
	if err != nil {
		fmt.Println("fetch error,", err)
	} else {
		fmt.Println(fetchRet.String())
	}
}
