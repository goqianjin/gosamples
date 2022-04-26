package partsv2

import (
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/qiniu/go-sdk/v7/client"

	"github.com/qiniu/go-sdk/v7/storage"
)

func newResumeV2Uploader(uphost string) *storage.ResumeUploaderV2 {
	// compose uploader
	cfg := storage.Config{
		Zone:          &storage.ZoneHuadong, // 空间对应的机房
		UseHTTPS:      false,                // 是否使用https域名
		UseCdnDomains: false,                // 上传是否使用CDN上传加速
	}
	//构建代理client对象
	urlParser, _ := url.Parse(uphost)
	tr := http.Transport{
		Proxy:                 http.ProxyURL(urlParser),
		ResponseHeaderTimeout: 3000 * time.Millisecond,
		Dial: (&net.Dialer{
			Timeout:   3000 * time.Millisecond,
			KeepAlive: 30 * time.Second,
		}).Dial,
	}
	client1 := http.Client{
		Transport: &tr,
	}
	resumeUploader := storage.NewResumeUploaderV2Ex(&cfg, &client.Client{Client: &client1})
	return resumeUploader
}
