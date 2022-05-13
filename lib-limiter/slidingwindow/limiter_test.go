package slidingwindow

import (
	"fmt"
	"testing"
	"time"
)

func TestLimiter(t *testing.T) {
	limiter := NewSliding(100*time.Millisecond, time.Second, 10)
	for i := 0; i < 5; i++ {
		fmt.Println(limiter.LimitTest())
	}
	time.Sleep(100 * time.Millisecond)
	for i := 0; i < 5; i++ {
		fmt.Println(limiter.LimitTest())
	}
	fmt.Println(limiter.LimitTest())
	for _, v := range limiter.windows {
		fmt.Println(v.timestamp, v.count)
	}

	fmt.Println("moments later...")
	time.Sleep(time.Second)
	for i := 0; i < 7; i++ {
		fmt.Println(limiter.LimitTest())
	}
	for _, v := range limiter.windows {
		fmt.Println(v.timestamp, v.count)
	}
}
