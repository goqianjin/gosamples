package main

import (
	"fmt"
	"net/url"
)

func main() {
	var bb1 bool
	fmt.Println(bb1)

	// TODO move to URL

	parse, err := url.Parse("http://node1.qiniu.com:8980/root/v1/user?a=1234&b=345")
	fmt.Println(err)
	fmt.Println(parse.Host)
	fmt.Println(parse.Path)

	fmt.Println("-=========")
	a1 := 1
	fmt.Printf("a = %v (%v)\n", a1, &a1)
	b1, a1 := 2, 2
	fmt.Printf("a = %v (%v), b = %v (%v)\n", a1, &a1, b1, &b1)
	fmt.Println("-=========")
	a2 := 1
	b2 := 2
	fmt.Printf("a = %v (%v)\n", a2, &a2)
	b2, a2 = 2, 2 // b2, a2 := 2, 2 --> compile error
	fmt.Printf("a = %v (%v), b = %v (%v)\n", a2, &a2, b2, &b2)
}
