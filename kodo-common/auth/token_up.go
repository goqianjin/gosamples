package auth

import (
	"encoding/json"
	"time"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

const (
	SignTypeUp SignType = "Up"
)

type UpTokenGenerator struct {
	*authKey
	// extra for up token
	putPolicy   *storage.PutPolicy
	putPolicyV2 *PutPolicyV2
}

func NewUpTokenGenerator(ak, sk string) *UpTokenGenerator {
	return &UpTokenGenerator{authKey: &authKey{ak: ak, sk: sk}}
}

func (k *UpTokenGenerator) WithPutPolicy(putPolicy *storage.PutPolicy) *UpTokenGenerator {
	k.putPolicy = putPolicy
	return k
}

func (k *UpTokenGenerator) WithPutPolicyV2(putPolicyV2 *PutPolicyV2) *UpTokenGenerator {
	k.putPolicyV2 = putPolicyV2
	return k
}

func (k *UpTokenGenerator) GenerateToken() string {
	return "UpToken " + k.GenerateRawToken()
}

func (k *UpTokenGenerator) GenerateRawToken() string {
	mac := qbox.NewMac(k.ak, k.sk)
	if k.putPolicyV2 != nil {
		return k.putPolicyV2.UploadToken(mac)
	} else if k.putPolicy != nil {
		return k.putPolicy.UploadToken(mac)
	} else {
		panic("missing put policy0")
	}
}

// -------- Helper --------

type PutPolicyV2 struct {
	Scope           string `json:"scope"`
	Expires         uint64 `json:"deadline"` // 截止时间（以秒为单位）
	IsPrefixalScope int    `json:"isPrefixalScope,omitempty"`
	InsertOnly      uint16 `json:"insertOnly,omitempty"` // Exclusive 的别名

	DetectMime          uint8  `json:"detectMime,omitempty"` // 若非0, 则服务端根据内容自动确定 MimeType
	FsizeMin            int64  `json:"fsizeMin,omitempty"`
	FsizeLimit          int64  `json:"fsizeLimit,omitempty"`
	MimeLimit           string `json:"mimeLimit,omitempty"`
	ForceSaveKey        bool   `json:"forceSaveKey,omitempty"`
	SaveKey             string `json:"saveKey,omitempty"`
	CallbackFetchKey    uint8  `json:"callbackFetchKey,omitempty"`
	CallbackURL         string `json:"callbackUrl,omitempty"`
	CallbackHost        string `json:"callbackHost,omitempty"`
	CallbackBody        string `json:"callbackBody,omitempty"`
	CallbackBodyType    string `json:"callbackBodyType,omitempty"`
	ReturnURL           string `json:"returnUrl,omitempty"`
	ReturnBody          string `json:"returnBody,omitempty"`
	PersistentOps       string `json:"persistentOps,omitempty"`
	PersistentNotifyURL string `json:"persistentNotifyUrl,omitempty"`
	PersistentPipeline  string `json:"persistentPipeline,omitempty"`
	EndUser             string `json:"endUser,omitempty"`
	DeleteAfterDays     int    `json:"deleteAfterDays,omitempty"`
	FileType            int    `json:"fileType,omitempty"`
	Exclusive           uint16 `json:"exclusive,omitempty"`       // 若为非0, 即使Scope为"Bucket:key"的形式也是insert only
	ForceInsertOnly     bool   `json:"forceInsertOnly,omitempty"` // 若为true,即使上传hash相同的文件也会报文件已存在,优先级高于InsertOnly

	MimeType string `json:"mimeType,omitempty"`
}

// UploadToken 方法用来进行上传凭证的生成
// 该方法生成的过期时间是现对于现在的时间
func (p *PutPolicyV2) UploadToken(cred *auth.Credentials) string {
	return p.uploadToken(cred)
}

func (p PutPolicyV2) uploadToken(cred *auth.Credentials) (token string) {
	if p.Expires == 0 {
		p.Expires = 3600 // 默认一小时过期
	}
	p.Expires += uint64(time.Now().Unix())
	putPolicyJSON, _ := json.Marshal(p)
	token = cred.SignWithData(putPolicyJSON)
	return
}
