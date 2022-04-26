package main

import (
	"fmt"
	"time"
)

func main() {
	nums := []int{11,12,13,14,15}
	noCacheChan := make(chan int)
	go func() {
		for i := range nums {
			fmt.Printf(" input nums[%d]=%d--\n", i, nums[i])
			noCacheChan <- nums[i]
		}
		time.Sleep(time.Second * 1800)

		//close(noCacheChan)
	}()
	time.Sleep(time.Second * 2)
	i := 0
	for ; ; {
		n1, ok := <- noCacheChan
		if !ok {
			fmt.Println("over...")
			break
		}
		fmt.Printf(" nums[%d]=%d-- %v\n", i, n1, ok)
		i++
	}
}
