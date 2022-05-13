package uputil

import (
	"bytes"
	"crypto/rand"
	"io"
)

func NewRandomBody(size int) io.Reader {
	data := make([]byte, size)
	n, err := rand.Read(data)
	if err != nil || n != size {
		panic("failed to generate random body")
	}
	return bytes.NewReader(data)
}
