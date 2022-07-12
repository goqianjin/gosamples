package internal

import (
	"fmt"
	"testing"
)

func TestAAA(t *testing.T) {
	type User struct {
		Name string
	}
	up1 := &User{Name: "shenqianjin"}
	fmt.Println(up1)
	updateFunc := func(up *User) {
		up = &User{Name: "shenqianjin2222"}
	}
	updateFunc(up1)
	fmt.Println(up1)

}
