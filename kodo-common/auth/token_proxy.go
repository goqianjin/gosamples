package auth

import (
	"strconv"

	"github.com/qianjin/kodo-common/proxyuser"
)

const (
	SignTypeQiniuProxy SignType = "QiniuProxy"
)

type QiniuProxyTokenGenerator struct {
	//signType SignType
	user *proxyuser.ProxyUser
}

func NewQiniuProxyTokenGenerator(user proxyuser.ProxyUser) *QiniuProxyTokenGenerator {
	copiedUser := user
	return &QiniuProxyTokenGenerator{user: &copiedUser}
}

func (k *QiniuProxyTokenGenerator) GenerateToken() string {
	token := k.generateToken()
	return string(SignTypeQiniuProxy) + " " + token
}

func (k *QiniuProxyTokenGenerator) generateToken() string {
	user := k.user
	form := make([]byte, 0, 64)
	form = appendUint32(form, "uid=", user.Uid)
	form = appendUint32(form, "&ut=", user.Utype)
	if user.Sudoer != 0 {
		form = appendUint32(form, "&suid=", user.Sudoer)
	}
	if user.UtypeSu != 0 {
		form = appendUint32(form, "&sut=", user.UtypeSu)
	}
	if user.Devid != 0 {
		form = appendUint32(form, "&dev=", user.Devid)
	}
	if user.Appid != 0 {
		form = appendUint32(form, "&app=", user.Appid)
	}
	if user.Expires != 0 {
		form = appendUint32(form, "&e=", user.Expires)
	}
	return string(form)
}

func appendUint32(form []byte, k string, v uint32) []byte {
	str := strconv.FormatUint(uint64(v), 10)
	form = append(form, k...)
	return append(form, str...)
}
