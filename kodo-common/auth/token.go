package auth

import (
	"net/http"
)

type TokenGenerator interface {
	GenerateToken(req *http.Request) string
}
