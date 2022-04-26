package form

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/qianjin/kodo-security/kodokey"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type MyPutReturn struct {
	Key    string
	Hash   string
	Fsize  int
	Bucket string
	Name   string
}

func TestFormUpByPut(t *testing.T) {
	key := "my-file-key-01"
	accessKey := kodokey.Prod_AK_shenqianjin
	secretKey := kodokey.Prod_SK_shenqianjin
	mac := qbox.NewMac(accessKey, secretKey)

	bucket := "qj-huanan-01"
	putPolicy := storage.PutPolicy{
		Scope:   bucket,
		Expires: 7200, // 两小时失效， 默认1小时
	}
	upToken := putPolicy.UploadToken(mac)

	config := storage.Config{}
	config.Zone = &storage.ZoneHuanan
	// config.UpHost = "r19tfe1y3.hn-bkt.clouddn.com"

	ret := MyPutReturn{}
	formUploader := storage.NewFormUploader(&config)
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"extra-field-01": "extra-value-01",
		},
	}
	myFile, err := os.Open("./test-file-in-go-exmaple.json")
	//myFile, err := os.Open("/Users/kodo/GolandProjects/go-example/kodo/updown/test-file-in-go-exmaple.json")
	//myFile := strings.NewReader("{\n  \"id\": 1234,\n  \"name\": \"Go Language Design\",\n  \"date\": \"2021-10-20 18:55:00\"\n}")
	if err != nil {
		log.Fatal(err)
	}
	stat, err := myFile.Stat()
	if err != nil {
		log.Fatal(err)
	}
	err = formUploader.Put(context.Background(), &ret, upToken, key, myFile, stat.Size(), &putExtra)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ret, myFile)
	fmt.Println(os.Getenv("PWD"))

	fmt.Println(os.Environ())
	// 下载
	// 公开空间
	domain := "http://r19tfe1y3.hn-bkt.clouddn.com"
	publicAccessUrl := storage.MakePublicURL(domain, key)
	fmt.Println(publicAccessUrl)
	// 私有空间
	deadline := time.Now().Add(time.Second * 3600).Unix()
	privateAccessUrl := storage.MakePrivateURL(mac, domain, key, deadline)
	fmt.Println(privateAccessUrl)

	// http 下载
	publicResp, err := http.Get(publicAccessUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer publicResp.Body.Close()

	fmt.Println(publicResp.Body)

	publicData, err := ioutil.ReadAll(publicResp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("public data: ", publicData)
	ioutil.WriteFile(key, publicData, 0644)

}
