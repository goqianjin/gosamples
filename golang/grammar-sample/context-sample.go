package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	loop := 10
	//wg := sync.WaitGroup{}
	//wg.Add(loop)
	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < loop; i++ {
		go func(i int) {
			mockContext(ctx, i)
			//wg.Done()
		}(i)

	}
	time.Sleep(time.Second)
	cancel()
	time.Sleep(time.Second)
}

func mockContext(ctx context.Context, i int) {
	select {
	case <-ctx.Done():
		fmt.Println("completed .... %d", i)
	}

}
