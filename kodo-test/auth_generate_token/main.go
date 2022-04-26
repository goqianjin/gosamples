package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
)

var Method = "POST"
var Path = "/sisyphus/fetch"
var RawQuery = ""
var Host = "api-z0.qiniu.com"
var contentType = "application/json"
var bodyStr = "{\n    \"url\" : \"http://r1krij46a.hd-bkt.clouddn.com/test01.txt\",\n    \"bucket\": \"qj-test-fetch-tar\",\n    \"key\": \"test01.txt\",\n  \"callbackurl\": \"http://localhost:80/callbackmock\",\n  \"callbackbody\": \"{}\",\n    \"callbackbodytype\": \"\",\n    \"file_type\":0\n}"

// prod: kodo
var accessKey = "hg5QoXO5n3pdoEkQYmtO2WnhkfI6DO794w9dLv2T"
var secretKey = "4FezG-jBQnT3W33szRiLP-yoRVf5kgSzehvZQfAa"
// test: general_storage_002@test.qiniu.io
//var accessKey = "hg5QoXO5n3pdoEkQYmtO2WnhkfI6DO794w9dLv2T"
//var secretKey = "4FezG-jBQnT3W33szRiLP-yoRVf5kgSzehvZQfAa"

func main() {
	// step 1: connect data
	data := Method + " " + Path
	if RawQuery != "" {
		data += "?" + RawQuery
	}
	data += "\nHost: " + Host
	if contentType != "" {
		data += "\nContent-Type: " + contentType
	}
	data += "\n\n"
	if bodyStr != "" && contentType != "" && contentType != "application/octet-stream" {
		data += bodyStr
	}
	// step 2:
	//hmac ,use sha1
	key := []byte(secretKey)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(data))
	sign := mac.Sum(nil)
	encodedSign := base64.URLEncoding.EncodeToString(sign)
	// step 3:
	fmt.Printf("Qiniu " + accessKey + ":" + encodedSign)
}
