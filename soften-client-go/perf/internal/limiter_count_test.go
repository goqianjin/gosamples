package internal

import (
	"fmt"
	"sync"
	"testing"
)

func TestConcurrencyLimiter(t *testing.T) {
	acquire := 0
	var acquireLock sync.RWMutex
	release := 0
	var releaseLock sync.RWMutex
	var group sync.WaitGroup
	loop := 10000
	group.Add(loop)
	limiter := NewConcurrencyLimiter(3)

	for i := 0; i < loop; i++ {
		go func(i int) {
			if !limiter.TryAcquire() {
				group.Done()
				return
			}
			acquireLock.Lock()
			acquire++
			acquireLock.Unlock()
			if i%100 == 0 {
				fmt.Println("acquire", acquire)
			}
			defer func() {
				limiter.Release()
				releaseLock.Lock()
				release++
				releaseLock.Unlock()
				if i%100 == 0 {
					fmt.Println("release", release)
				}
			}()
			// do biz

			group.Done()

		}(i)
	}
	group.Wait()

}
