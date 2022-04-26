package main

import (
	"errors"
	"fmt"
)

func main() {
	ok1, err1 := test(1)
	ok2, err2 := test(2)
	fmt.Printf("ok: %v, err: %v\n", ok1, err1)
	fmt.Printf("ok: %v, err: %v\n", ok2, err2)
}

func test(i int) (ok string, err error) {
	if i % 2 == 0 {
		ok, err := testOK()
		fmt.Printf("ok: %v, err: %v\n", ok, err)
		return ok, err
	} else {
		ok, err := testErr()
		fmt.Printf("ok: %v, err: %v\n", ok, err)
		return ok, err
	}
}

func testOK() (ok string, err error){
	return "ok", nil
}
func testErr() (ok string, err error){
	return "no", errors.New("I'a an error.")
}



