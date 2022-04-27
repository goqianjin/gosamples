package auth

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/qianjin/kodo-common/authkey"
	"github.com/qiniu/go-sdk/v7/auth"
)

const qiniuHeaderPrefix = "X-Qiniu-"

const (
	SignTypeQiniu      SignType = "Qiniu"
	SignTypeQiniuAdmin SignType = "QiniuAdmin"
	SignTypeQBox       SignType = "QBox"
	SignTypeQBoxAdmin  SignType = "QBoxAdmin"
)

type ManagedTokenGenerator struct {
	*authkey.AuthKey
	signType SignType

	// extra for admin sign
	suInfo string
}

func (k *ManagedTokenGenerator) WithSignType(signType SignType) *ManagedTokenGenerator {
	k.signType = signType
	return k
}

func (k *ManagedTokenGenerator) WithSuInfo(uid, appId uint32) *ManagedTokenGenerator {
	k.suInfo = fmt.Sprintf("%v/%v", uid, appId)
	return k
}

func NewManagedTokenGeneratorByKey(authKey *authkey.AuthKey) *ManagedTokenGenerator {
	return &ManagedTokenGenerator{AuthKey: authKey}
}

func NewManagedTokenGenerator(ak, sk string) *ManagedTokenGenerator {
	return &ManagedTokenGenerator{AuthKey: &authkey.AuthKey{AK: ak, SK: sk}}
}

func (k *ManagedTokenGenerator) GenerateToken(req *http.Request) string {
	token := k.generateToken(req)
	if k.signType == SignTypeQiniuAdmin || k.signType == SignTypeQBoxAdmin {
		token = k.suInfo + ":" + token
	}
	return string(k.signType) + " " + token
}

func (k *ManagedTokenGenerator) generateToken(req *http.Request) string {
	// step 1: connect data
	var data string
	// 签request method: Qiniu & QiniuAdmin
	if k.signType == SignTypeQiniu || k.signType == SignTypeQiniuAdmin {
		data += req.Method + " "
	}
	// 签Path and RawQuery: all
	data += req.URL.Path
	if req.URL.RawQuery != "" {
		data += "?" + req.URL.RawQuery
	}
	// 签Host and Content-Type headers: Qiniu & QiniuAdmin
	if k.signType == SignTypeQiniu || k.signType == SignTypeQiniuAdmin {
		data += "\nHost: " + req.Host
		if req.Header.Get("Content-Type") != "" {
			data += "\nContent-Type: " + req.Header.Get("Content-Type")
		}
	}
	// 签suInfo: QiniuAdmin & QBoxAdmin
	if k.signType == SignTypeQiniuAdmin {
		data += "\nAuthorization: QiniuAdmin " + k.suInfo
	} else if k.signType == SignTypeQBoxAdmin {
		data += "\nAuthorization: QBoxAdmin " + k.suInfo
	}
	// 签X-Qiniu-* headers: Qiniu & QiniuAdmin
	if k.signType == SignTypeQiniu || k.signType == SignTypeQiniuAdmin {
		data = k.signQiniuHeaderValues(req.Header, data)
	}
	// 签分割符：
	if k.signType == SignTypeQBox { // QBox 一个空行
		data += "\n"
	} else { // QBoxAdmin, Qiniu, QiniuAdmin 两个个空行
		data += "\n\n"
	}
	// 签Body
	if ((k.signType == SignTypeQiniu || k.signType == SignTypeQiniuAdmin) && k.incBody4Qiniu(req, req.Header.Get("Content-Type"))) ||
		((k.signType == SignTypeQBox || k.signType == SignTypeQBoxAdmin) && k.incBody4QBox(req)) {
		bbody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			_ = fmt.Errorf("failed to read req.Body")
		}
		sbody := string(bbody)
		if sbody != "" {
			data += sbody
		}
		bodyReader := strings.NewReader(sbody)
		req.Body = io.NopCloser(bodyReader)
	}
	// step 2:
	//hmac ,use sha1
	mac := auth.New(k.AK, k.SK)
	token := mac.Sign([]byte(data))
	// step 3:
	// fmt.Printf("generated token: " + token)
	return token
}

// 判断Body是否应该计入签名: QBox & QBoxAdmin
func (k *ManagedTokenGenerator) incBody4QBox(req *http.Request) bool {

	if req.Body == nil || req.ContentLength == 0 {
		return false
	}
	if ct, ok := req.Header["Content-Type"]; ok {
		switch ct[0] {
		case "application/x-www-form-urlencoded":
			return true
		}
	}
	return false
}

// 判断Body是否应该计入签名: Qiniu & QiniuAdmin
func (k *ManagedTokenGenerator) incBody4Qiniu(req *http.Request, ctType string) bool {

	return req.ContentLength != 0 && req.Body != nil && ctType != "" && ctType != "application/octet-stream"
}

func (k *ManagedTokenGenerator) signQiniuHeaderValues(header map[string][]string, data string) string {
	var keys []string
	for key := range header {
		if len(key) > len(qiniuHeaderPrefix) && key[:len(qiniuHeaderPrefix)] == qiniuHeaderPrefix {
			keys = append(keys, key)
		}
	}
	if len(keys) == 0 {
		return data
	}

	if len(keys) > 1 {
		sort.Sort(sort.StringSlice(keys))
	}
	for _, key := range keys {
		if len(header[key]) == 0 {
			continue
		}
		data = data + "\n" + key + ": " + strings.Join(header[key], "; ")
	}
	return data
}
