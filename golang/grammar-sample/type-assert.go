package main

import "fmt"

func main() {
	a := 12
	var aI interface{}
	aI = a
	b := aI.(int64)
	fmt.Println(b)
}
