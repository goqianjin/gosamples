package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	fmt.Println("main start....")

	for index := 1; index <= 5; index++ {
		go mockPanicByCondition(index)
	}
	time.Sleep(time.Second)
	fmt.Println("main sleeping....")
	fmt.Println("main end....")
}

func mockPanicByCondition(index int) {
	fmt.Println("goroutime " + strconv.Itoa(index) + " start...")
	if index == 3 {
		panic("panic happened here...")
	}
	fmt.Println("goroutime " + strconv.Itoa(index) + " end...")
}
