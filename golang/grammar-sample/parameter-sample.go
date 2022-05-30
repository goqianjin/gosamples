package main

import "fmt"

func main() {
	type User struct {
		Name   string
		getAge func() int
	}

	user := &User{Name: "zhang"}
	fmt.Println(user)
	setAgeFunc := func(getAge *func() int) {
		aa := func() int {
			return 12
		}
		getAge = &aa
	}
	setAgeFunc(&user.getAge)
	fmt.Println(user)
	fmt.Println(user.getAge())
}
