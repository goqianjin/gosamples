package main

import "fmt"

func main() {
	a := 1
	fmt.Printf("a = %v\n", a)
	defer func(a1 int) {
		fmt.Printf("a1 in defer = %v\n", a1)

	}(a)
	a = 2
	fmt.Printf("a = %v\n", a)

	fmt.Println("------------defer in for ----start")
	for i := 0; i < 3; i++ {
		defer func() {
			fmt.Printf("defer in for, i = %v\n", i)
		}()
		fmt.Printf("inner for, i = %v\n", i)

	}
	fmt.Println("------------defer in for ----end")
	deferTest03()
}

func deferTest03() {
	defer func() {
		fmt.Println("first defer ...")
	}()
	defer func() {
		fmt.Println("second defer ...")
	}()
	fmt.Println("------------defer test 03 ----start ---")
	fmt.Println("------------defer test 03 ----processing ---")
	fmt.Println("------------defer test 03 ----end ---")

}
