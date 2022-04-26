package main

import (
	"fmt"
	"time"
)

func main() {
	// time.Tick test
	t1, err := time.Parse(time.RFC3339, "2021-11-16T00:00:00Z")
	if err != nil {
		panic(err)
	}
	fmt.Println(time.Since(t1))
	fmt.Println(time.Since(t1).Hours())


	//testTimeTick()


}

func testTimeTick() {
	fmt.Println("time tick test - starting ...")
	c := time.Tick(time.Duration(1000) * time.Millisecond)
	var count int = 1
	for range c {
		fmt.Printf("%v: do job [%d]\n", time.Now().String(), count)
		count++
	}
	fmt.Println("time tick test - ended ...")
}