package main

import (
	"fmt"
	"reflect"
)

func main() {
	type user struct {
		id   int
		name string
	}
	ch1 := make(chan *user, 20)
	ch2 := make(chan *user, 20)

	user1 := &user{0, "zhang"}
	go func(c1 chan *user) {
		for i := 0; i <= 1000; i++ {
			c1 <- user1
		}
	}(ch1)

	user2 := &user{1, "li"}
	go func(c2 chan *user) {
		for i := 0; i <= 1000; i++ {
			c2 <- user2
		}
	}(ch2)

	chs := []<-chan *user{ch1, ch2}
	cases := make([]reflect.SelectCase, len(chs))
	for i, ch := range chs {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}
	chosen, value, ok := reflect.Select(cases)
	fmt.Println(chosen, value, ok)
	ivalue := value.Interface()
	rvalue, ok := ivalue.(*user)
	if !ok {
		panic(fmt.Sprintf("convert %v to user failed", value))
	}
	fmt.Printf("user1: %p, user2: %p, ivalue: %p, value: %p", user1, user2, ivalue, rvalue)
}
