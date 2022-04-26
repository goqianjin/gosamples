package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// 从环境变量获取参数: fileserver -p 8080
	portAddress := flag.Int("p", 8080, "specify the port of your file server")
	flag.Parse()
	port := *portAddress
	// 提示监听端口。可指定端口，默认：8000
	r := mux.NewRouter().SkipClean(true)
	fmt.Printf("Serving HTTP on http://0.0.0.0:%d", port)
	r.Handle("/get", get()).Methods(http.MethodGet)
	r.Handle("/", get()).Methods(http.MethodGet)
	// 启动监听 HTTP 服务
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}

func get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := "Hello Golang!"
		vars := mux.Vars(r)
		v1, ok1 := vars["key1"]
		fmt.Printf("key1=%v, exists=%v", v1, ok1)
		v2, ok2 := vars["key2"]
		fmt.Printf("key2=%v, exists=%v", v2, ok2)

		fmt.Println("---------")
		getkey := "x-qn-status1"
		fmt.Printf("---get header: key: %s : %s", getkey, r.Header.Get(getkey))
		for k, v := range r.Header {
			fmt.Printf("header: %s: %s\n", k, v)
		}
		fmt.Println("---------end")

		w.Header().Add("X-QN-Status1", "11")
		w.Header().Add("X-qn-Status2", "22")
		w.Header().Add("X-Qn-Status3", "33")
		w.Header().Add("X-Qn-StatusCode", "4455")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))
	}
}
