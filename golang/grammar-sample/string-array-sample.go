package main

import "fmt"

func main() {
	str1 := "上海"
	str2 := "SH"
	fmt.Printf("len(%s)=%d\n", str1, len(str1))
	fmt.Printf("len(%s)=%d\n", str2, len(str2))
	arr1 := []rune(str1)
	arr2 := []rune(str2)
	fmt.Printf("len([]rune(%s))=%d\n", str1, len(arr1))
	fmt.Printf("len([]rune(%s))=%d\n", str2, len(arr2))
	for index := range arr1 {
		c1 := arr1[index]
		fmt.Printf("[]rune(%s)[%d]=%v ", str1, index, string(c1))
	}
	fmt.Println()
	for index := range arr2 {
		c1 := arr2[index]
		fmt.Printf("[]rune(%s)[%d]=%v ", str2, index, string(c1))
	}
	fmt.Println()
	for index := range str1 {
		c1 := str1[index]
		fmt.Printf("%v ", string(c1))
	}
	fmt.Println()
	for index := range str2 {
		c1 := str2[index]
		fmt.Printf("%v ", string(c1))
	}
	fmt.Println()
}
