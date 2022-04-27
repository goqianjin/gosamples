package auth

import (
	"net/http"
)

type SignType string

type TokenGenerator interface {
	GenerateToken(req *http.Request) string
}
