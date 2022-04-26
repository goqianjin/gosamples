package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/qianjin/kodo-security/osskey"
)

func main() {
	endpoint := "https://oss-cn-hangzhou.aliyuncs.com"
	ak := osskey.OSS_AK_shenqianjin
	sk := osskey.OSS_SK_shenqianjin
	bucket := "qianjin-test-01"

	// 创建OSSClient实例。
	// yourEndpoint填写Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://oss-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
	// 阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维，请登录RAM控制台创建RAM用户。
	client, err := oss.New(endpoint, ak, sk, oss.SetLogLevel(oss.Debug))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 初始化标签。
	tag1 := oss.Tag{
		Key:   "key1",
		Value: "value1",
	}
	tag2 := oss.Tag{
		Key:   "key2",
		Value: "value2",
	}
	tagging := oss.Tagging{
		Tags: []oss.Tag{tag1, tag2},
	}
	// 填写Bucket名称，例如examplebucket。
	// 设置Bucket标签。
	err = client.SetBucketTagging(bucket, tagging)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 填写Bucket名称，例如examplebucket。
	// 获取Bucket标签信息。
	ret, err := get(client, bucket, nil) //client.GetBucketTagging(bucket)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	// 打印标签个数。
	fmt.Println("Tag length: ", ret)
}

func get(client *oss.Client, bucketName string, options ...oss.Option) (string, error) {
	var out string
	params := map[string]interface{}{}
	params["key1"] = "11"
	params["tagging"] = nil
	params["key3"] = "33"
	params["skey3"] = "44"
	params["zzkey3"] = "66"
	resp, err := client.Conn.Do("GET", bucketName, "", params, nil, nil, 0, nil)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("failed to readAll from resp.Body: %v\n", err)
	}
	//err = json.Unmarshal(bytes, &out)
	out = string(bytes)
	return out, err
}
