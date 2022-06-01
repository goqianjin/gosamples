package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func main() {
	// 从环境变量获取参数: fileserver -p 8080
	portAddress := flag.Int("p", 8080, "specify the port of your file server")
	flag.Parse()
	port := *portAddress
	// 获取当前目录
	wd, _ := os.Getwd()
	// 提示监听端口。可指定端口，默认：8000
	fmt.Printf("Serving HTTP on http://0.0.0.0:%d", port)
	http.HandleFunc("/get/", getByCode())
	http.HandleFunc("/get", get())
	http.Handle("/", http.FileServer(http.Dir(wd)))
	// 启动监听 HTTP 服务
	log.Fatal(http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), nil))
}

func getByCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		segs := strings.Split(r.URL.Path, "/")
		if len(segs) != 3 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if code, err := strconv.Atoi(segs[2]); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else if code >= 200 || code < 1000 {
			w.WriteHeader(code)
			w.Write([]byte("I am body.\n"))
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := "我是文件内容。\n" +
			"How are you?\n" +
			"结束了！\n" +
			"再见！"
		attname := r.URL.Query().Get("attname")
		filename := r.URL.Query().Get("filename")
		filenamex := r.URL.Query().Get("filenamex")
		if attname == "" {
			attname = "文件名1.txt"
		}
		cdValue := "attachment"
		if filename == "ex" {
			cdValue += "; filename=" + attname
		} else if filename == "en" {
			cdValue += "; filename=" + url.PathEscape(attname)
		}
		if filenamex == "ex" {
			cdValue += "; filename*=utf-8''" + attname
		} else if filenamex == "en" {
			cdValue += "; filename*=utf-8''" + url.PathEscape(attname)
		} else if filenamex == "ens" {
			cdValue += "; filename*=utf-8' '" + url.PathEscape(attname)
		}
		w.Header().Set("Content-Disposition", cdValue)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))
	}
}
