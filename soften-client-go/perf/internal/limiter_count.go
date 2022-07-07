package internal

import (
	"fmt"
)

type ConcurrencyLimiter interface {
	TryAcquire() bool
	Acquire()
	Release()
}

type concurrencyLimiter struct {
	ch chan struct{}
}

func NewConcurrencyLimiter(n int) ConcurrencyLimiter {
	if n <= 0 {
		panic(fmt.Sprintf("concurrency should more than one. current: %d", n))
	}
	ch := make(chan struct{}, n)
	for i := 0; i < n; i++ {
		ch <- struct{}{}
	}
	return &concurrencyLimiter{ch: ch}
}

func (l *concurrencyLimiter) TryAcquire() bool {
	select {
	case <-l.ch:
		return true
	default:
		return false
	}
}

func (l *concurrencyLimiter) Acquire() {
	//l.chLock.Lock()
	//defer l.chLock.Unlock()
	<-l.ch
}

func (l *concurrencyLimiter) Release() {
	//l.chLock.Lock()
	//defer l.chLock.Unlock()
	l.ch <- struct{}{}
}
