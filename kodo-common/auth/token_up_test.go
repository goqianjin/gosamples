package auth

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/qiniu/go-sdk/v7/storage"

	"github.com/qianjin/kodo-common/auth/config"
)

func TestGenerateUpToken_Dev(t *testing.T) {
	contentType := "application/json"
	request, _ := http.NewRequest(http.MethodPost, "/fetch", strings.NewReader(""))
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("Host", "10.200.20.23:9433")
	putPolicy := &storage.PutPolicy{
		Scope: "bucket02",
	}

	token := NewUpTokenGenerator(kodokey.Dev_AK_general_storage_002, kodokey.Dev_SK_general_torage_002).
		WithPutPolicy(putPolicy).
		GenerateToken(request)
	fmt.Println("Generated token: " + token)
}

func TestGenerateUpToken_Prod(t *testing.T) {
	contentType := "application/json"
	request, _ := http.NewRequest(http.MethodPost, "/sisyphus/fetch", strings.NewReader(""))
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("Host", "api-z0.qiniu.com")

	putPolicy := &storage.PutPolicy{
		Scope: "bucket02",
	}
	token := NewUpTokenGenerator(kodokey.Prod_AK_shenqianjin, kodokey.Prod_SK_shenqianjin).
		WithPutPolicy(putPolicy).
		GenerateToken(request)
	fmt.Println("Generated token: " + token)
}
