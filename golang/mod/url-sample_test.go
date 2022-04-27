package main

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/url"
	"strings"
	"testing"
)

func TestAbc(t *testing.T) {
	//b := nil == nil
	//fmt.Println(b)
	rawUrl := "http://127.0.0.1:8080/url/123"
	changeHost := "http://192.168.1.1:8080"

	fmt.Println(url.ParseRequestURI(changeHost))
	fmt.Println(url.Parse(changeHost))
	fmt.Println(govalidator.IsURL(rawUrl))
	fmt.Println(govalidator.IsURL(changeHost))
	fmt.Println(url.Parse(changeHost))
	newUrl, _ := url.Parse(rawUrl)
	fmt.Println(newUrl.Scheme + "://" + newUrl.Host)
	fmt.Println(newUrl.Hostname())
	newUrlHost := newUrl.Hostname()
	newUrlPort := newUrl.Port()
	newUrlPath := newUrl.Path
	stringUrl := newUrl.String()
	newUrl.Host = changeHost + ":" + newUrl.Port()
	fmt.Println(newUrlHost, newUrlPort, newUrlPath, stringUrl, newUrl)

	fmt.Println(strings.Replace("a...b", ".", "", -1))

	//

}

func isURL(s string) bool {
	// pattern := "(https?)://"
	return true
}

func TestUrlParse(t *testing.T) {
	url1, _ := url.Parse("https://origin.qiniu.com:8080/user/list?id=12345&name=zhang")
	fmt.Println(url1.String())
	fmt.Println(url1.Host)
	fmt.Println(url1.Hostname())
	fmt.Println(url1.Scheme)
	url1.Host = "www.baidu.com"
	fmt.Println(url1.String())
}
