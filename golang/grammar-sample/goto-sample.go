package main

import (
	"fmt"
	"time"
)

func main() {
	msgDoneCh := make(chan int, 1024)
	fmt.Printf("%v: main starting...\n", time.Now())
	tm := time.NewTimer(time.Second * 3)
	defer tm.Stop()
	msgDoneCh <- 3
	msgDoneCh <- 6
CleanLoop:
	for {
		select {
		case msg, ok := <-msgDoneCh:
			if ok {
				fmt.Printf("received msg: %v\n", msg)
			} else {
				fmt.Printf("received nothing. ok: %v\n", ok)
			}
		case <-tm.C:
			break CleanLoop
		}
	}
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		fmt.Printf("%v: main-later...\n", time.Now())
	}
	k := 0
	go func() {
		time.Sleep(time.Second * 2)
		k = k + 1
		msgDoneCh <- k
	}()

	fmt.Printf("%v: main ending...\n", time.Now())
}
