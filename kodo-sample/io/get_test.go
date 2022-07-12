package io

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	timeout := time.Second * 1
	ctx, cancelFun := context.WithTimeout(context.Background(), timeout)
	defer cancelFun()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://video-origin.yunjilink.com/64008c85d5a042cb9b1fe12c9757d5f8.mp4", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Range", "bytes=0-")
	// 关闭整数验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	do, err := client.Do(req)
	assert.Nil(t, err)
	fmt.Println(do.Status)
	fmt.Println(do.Header)
}

func TestGetByCancel(t *testing.T) {
	timeout := time.Second * 1
	ctx, cancelFun := context.WithCancel(context.Background())
	go func() {
		time.Sleep(timeout)
		cancelFun()
	}()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://video-origin.yunjilink.com/64008c85d5a042cb9b1fe12c9757d5f8.mp4", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Range", "bytes=0-")
	// 关闭整数验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	do, err := client.Do(req)
	assert.Nil(t, err)
	fmt.Println(do.Status)
	fmt.Println(do.Header)
}
