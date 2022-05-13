package main

import (
	"fmt"
	"sync"
)

func main() {
	type User struct {
		Name     string
		initOnce sync.Once
	}
	user := &User{}
	user.initOnce.Do(func() {
		fmt.Println("aaa")
	})
}
