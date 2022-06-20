package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan int, 20)
	ch2 := make(chan int, 20)
	var ch3 chan int
	var cnt [3]int

	go func(c1 chan int) {
		for i := 0; i <= 1000; i++ {
			c1 <- 0
		}
	}(ch1)

	go func(c2 chan int) {
		for i := 0; i <= 1000; i++ {
			c2 <- 1
		}
	}(ch2)
	for i := 0; i < 1000; i++ {

		//等10毫秒，确保两个channel都已准备就绪
		//time.Sleep(10 * time.Millisecond)

		var index int
		time.Sleep(time.Millisecond)
		select {
		case index = <-ch1:
			cnt[index]++
		case index = <-ch1:
			cnt[index]++
		case index = <-ch1:
			cnt[index]++
		case index = <-ch2:
			cnt[index]++
		case <-ch3:
			cnt[2]++

		}
	}
	fmt.Printf("cnt=%v\n", cnt)
}
