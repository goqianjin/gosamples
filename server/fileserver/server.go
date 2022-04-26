package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
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
	// 启动监听 HTTP 服务
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), http.FileServer(http.Dir(wd))))
}
