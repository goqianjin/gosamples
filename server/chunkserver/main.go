package main

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter().SkipClean(true)
	r.HandleFunc("/get/fsize/{fsize}", get).Methods(http.MethodGet)
	http.ListenAndServe(":8888", r)
}

func get(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	// parse params
	fsize, err := strconv.ParseInt(vars["fsize"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("invalid fsize: %v", vars["fsize"])))
		return
	}
	body := make([]byte, fsize)
	n, err := rand.Read(body)
	fmt.Printf("n = %d, err: %v\n", n, err)
	w.Header().Set("Content-Disposition", "attachment;")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
