package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestDownloadChunk(t *testing.T) {
	resp, err := http.Get("http://down.ws.eebbk.net/xzza/bbk.wy/a.log")
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(len(body))
	body, err = ioutil.ReadAll(resp.Body)
	fmt.Println(len(body))
	//sync.Map{}
}

func TestGetRefe(t *testing.T) {
	url := "http://localhost:6000/v2/reference?fh=CpYAABAAAAAAAD2YDl6no_UO43ctIXnKV9ojo6FVAQAAAAAAEAAAAAAABpb_fyUDAAAp92H3QzXoFjW9lyQFAAAAAAAQAAAAAABTW1buAAAAAD2YDl6no_UO43ctIXnKV9ojo6FV"

	for i := 0; i < 10000; i++ {

		resp, err := http.Get(url)
		if err != nil {
			// handle error
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		fmt.Println("result: " + string(body))
		time.Sleep(5 * time.Millisecond)
	}
}
