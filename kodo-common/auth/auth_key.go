package auth

import (
	"net/http"
)

type SignType string

type authKey struct {
	ak       string
	sk       string
	signType SignType
}

func (k *authKey) WithSignType(signType SignType) *authKey {
	k.signType = signType
	return k
}

func (k *authKey) GenerateToken(req *http.Request) string {
	panic("GenerateToken has not been implemented")
}
