package auth

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/qianjin/kodo-security/kodokey"
)

func TestGenerate_entryInfo_AdminProd(t *testing.T) {
	contentType := "application/x-www-form-urlencoded"
	BodyStr := "itbl=504699156&key=fragments/z1.wypd.wypd/1631977359914-1631977367102.ts"
	request, _ := http.NewRequest(http.MethodPost, "/entryinfo", strings.NewReader(BodyStr))
	request.Header.Set("Content-Type", contentType)
	//request.Header.Set("Host", "rs_sample-z1.qbox.me")
	token := NewManagedTokenGenerator(kodokey.Prod_AK_admin, kodokey.Prod_SK_admin).
		WithSignType(SignTypeQiniuAdmin).WithSuInfo(123, 0).
		GenerateToken(request)
	fmt.Println("Generated token: " + token)
}
func TestGenerate_entryInfo_AdminDev(t *testing.T) {
	contentType := "application/x-www-form-urlencoded"
	BodyStr := "itbl=504699156&key=fragments/z1.wypd.wypd/1631977359914-1631977367102.ts"
	request, _ := http.NewRequest(http.MethodPost, "/entryinfo", strings.NewReader(BodyStr))
	request.Header.Set("Content-Type", contentType)
	//request.Header.Set("Host", "rs_sample-z1.qbox.me")
	token := NewManagedTokenGenerator(kodokey.Dev_AK_admin, kodokey.Dev_SK_admin).
		WithSignType(SignTypeQiniuAdmin).WithSuInfo(123, 0).
		GenerateToken(request)
	fmt.Println("Generated token: " + token)
}

func TestGenerate_fetch_Dev(t *testing.T) {
	contentType := "application/json"
	BodyStr := "{\n    \"url\" : \"http://r1krij46a.hd-bkt.clouddn.com/test01.txt\",\n    \"bucket\": \"qj-test-fetch-tar\",\n    \"key\": \"test01.txt\",\n  \"callbackurl\": \"http://localhost:80/callbackmock\",\n  \"callbackbody\": \"{}\",\n    \"callbackbodytype\": \"\",\n    \"file_type\":0\n}"
	request, _ := http.NewRequest(http.MethodPost, "/fetch", strings.NewReader(BodyStr))
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("Host", "10.200.20.23:9433")

	token := NewManagedTokenGenerator(kodokey.Dev_AK_general_storage_002, kodokey.Dev_SK_general_torage_002).
		WithSignType(SignTypeQiniu).
		GenerateToken(request)
	fmt.Println("Generated token: " + token)
}

func TestGenerate_fetch_Prod(t *testing.T) {
	contentType := "application/json"
	BodyStr := "{\n    \"url\" : \"http://r1krij46a.hd-bkt.clouddn.com/test01.txt\",\n    \"bucket\": \"qj-test-fetch-tar\",\n    \"key\": \"test01.txt\",\n  \"callbackurl\": \"http://localhost:80/callbackmock\",\n  \"callbackbody\": \"{}\",\n    \"callbackbodytype\": \"\",\n    \"file_type\":0\n}"
	request, _ := http.NewRequest(http.MethodPost, "/sisyphus/fetch", strings.NewReader(BodyStr))
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("Host", "api-z0.qiniu.com")
	token := NewManagedTokenGenerator(kodokey.Prod_AK_shenqianjin, kodokey.Prod_SK_shenqianjin).
		WithSignType(SignTypeQiniu).
		GenerateToken(request)
	fmt.Println("Generated token: " + token)
}
