package main

import "fmt"

func main() {
	type User struct {
		Name           string
		tryRefreshFunc func(err error) bool
	}
	user := &User{Name: "zhangsan"}
	fmt.Printf("user: %+v\n", user)
	_ = user.tryRefreshFunc != nil && user.tryRefreshFunc(nil)
	fmt.Printf("user: %+v\n", user)
	user.tryRefreshFunc = func(err error) bool {
		fmt.Printf("executed tryRefreshFunc. parm: %+v\n", err)
		return false
	}
	_ = user.tryRefreshFunc != nil && user.tryRefreshFunc(nil)
	fmt.Printf("user: %+v\n", user)
}
