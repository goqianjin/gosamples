package main

import (
	"flag"
	"log"
	"net/http"
)

type CorsHandler struct {
	http.Handler
}

func (c *CorsHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, r.URL)
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	c.Handler.ServeHTTP(rw, r)
	return
}

func main() {
	addr := flag.String("addr", ":8000", "server listen addr")
	dir := flag.String("dir", ".", "dir")
	flag.Parse()
	log.Println("http fileserver running at", *addr)
	log.Panicln(http.ListenAndServe(*addr, &CorsHandler{http.FileServer(http.Dir(*dir))}))
}
