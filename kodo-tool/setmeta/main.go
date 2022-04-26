package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/client"
	"github.com/qiniu/go-sdk/v7/storage"
)

func main() {
	//ak := "557TpseUM8ovpfUhaw8gfa2DQ0104ZScM-BTIcBx"
	//sk := "d9xLPyreEG59pR01sRQcFywhm4huL-XEpHHcVa90"
	ak := "dxVQk8gyk3WswArbNhdKIwmwibJ9nFsQhMNUmtIM"
	sk := "s-95BDHpzyHPzHAe7WGdAeTy98vVwqdki-0U027j"
	bucketManager := getBucketManager(ak, sk)
	chmeta(bucketManager, "qj-bucket-z0", "test02.txt", nil)
	//TestListFiles(log.Default(), bucketManager, "fusionlog")

}

func getBucketManager(ak, sk string) *storage.BucketManager {
	mac := auth.New(ak, sk)
	clt := client.Client{
		Client: &http.Client{
			Timeout: time.Second * 30,
		},
	}
	cfg := storage.Config{}
	//cfg.UseCdnDomains = true
	return storage.NewBucketManagerEx(mac, &cfg, &clt)
}

func TestListFiles(t *log.Logger, bucketManager *storage.BucketManager, testBucket string) {
	limit := 100
	prefix := "" //listfiles/"
	entries, _, _, hasNext, err := bucketManager.ListFiles(testBucket, prefix, "", "", limit)
	if err != nil {
		t.Fatalf("ListFiles() error, %s", err)
	}
	t.Println(hasNext)

	/*if hasNext {
		t.Fatalf("ListFiles() failed, unexpected hasNext")
	}*/

	if len(entries) != limit {
		t.Fatalf("ListFiles() failed, unexpected items count, expected: %d, actual: %d", limit, len(entries))
	}

	for _, entry := range entries {
		t.Printf("ListItem:\n%s", entry.String())
	}
	//encodedUri := storage.EncodedEntry(testBucket, entries[0].Key)
	//fileInfo, err := bucketManager.Stat(testBucket, entries[0].Key)
	client.DebugMode = true
	fileInfo, err := StatWithOpts(bucketManager, testBucket, "v2/beta-media.hunliji.com_2022-02-28-11_part-00000.gz", nil)
	t.Printf("fileInfo: %+v\n", fileInfo)
}

// FileInfo 文件基本信息
type FileInfoExt struct {
	storage.FileInfo
	XQnMeta map[string]string `json:"x-qn-meta,omitempty" bson:"-"` // 对外字段

}

// StatWithParts 用来获取一个文件的基本信息以及分片信息
func StatWithOpts(m *storage.BucketManager, bucket, key string, opt *storage.StatOpts) (info FileInfoExt, err error) {
	reqHost, reqErr := m.RsReqHost(bucket)
	if reqErr != nil {
		err = reqErr
		return
	}

	reqURL := fmt.Sprintf("%s%s", reqHost, storage.URIStat(bucket, key))
	if opt != nil {
		if opt.NeedParts {
			reqURL += "?needparts=true"
		}
	}
	err = m.Client.CredentialedCall(context.Background(), m.Mac, auth.TokenQiniu, &info, "POST", reqURL, nil)
	return
}

// StatWithParts 用来获取一个文件的基本信息以及分片信息
func chmeta(m *storage.BucketManager, bucket, key string, opt *storage.StatOpts) (err error) {
	reqHost, reqErr := m.RsReqHost(bucket)
	if reqErr != nil {
		err = reqErr
		return
	}

	//reqURL := fmt.Sprintf("%s%s/%s/%s", reqHost, URIChgm(bucket, key),
	//	base64.URLEncoding.EncodeToString([]byte("A")), "12")
	reqURL := fmt.Sprintf("%s%s/%s/%s", reqHost, URIChgm(bucket, key),
		"x-qn-meta-A", base64.URLEncoding.EncodeToString([]byte("")))

	client.DebugMode = true
	err = m.Client.CredentialedCall(context.Background(), m.Mac, auth.TokenQiniu, nil, "POST", reqURL, nil)
	return
}

// URIStat 构建 stat 接口的请求命令
func URIChgm(bucket, key string) string {
	return fmt.Sprintf("/chgm/%s", storage.EncodedEntry(bucket, key))
}
