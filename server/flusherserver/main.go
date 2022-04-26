package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {

		flusher, ok := writer.(http.Flusher)
		if !ok {
			panic("expected http.ResponseWriter to be an http.Flusher")
		}

		for i := 0; i < 5; i++ {
			fmt.Fprintf(writer, "chunk [%02d]: %v\n", i, time.Now())
			flusher.Flush()
			time.Sleep(time.Second)
		}
	})

	http.ListenAndServe(":8080", nil)
}
