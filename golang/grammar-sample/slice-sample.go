package main

import "fmt"

func main() {
	var s1 []string
	fmt.Printf("s1: %v, len(s1): %v, cap(s1): %v\n", s1, len(s1), cap(s1))
	fmt.Printf("%v\n", s1 == nil)
	for index, ele := range s1 {
		fmt.Printf("index: %v, element: %v\n", index, ele)
	}
	s2 := []string{"a", "b", "c"}
	fmt.Printf("s2: %v, len(s2): %v, cap(s2): %v\n", s2, len(s2), cap(s2))
	s1 = s2
	fmt.Printf("s1: %v, len(s1): %v, cap(s1): %v\n", s1, len(s1), cap(s1))
	s3 := make([]string, 32, 32)
	fmt.Printf("s3: %v, len(s3): %v, cap(s3): %v\n", s3, len(s3), cap(s3))
	s3 = make([]string, 0, 32)
	fmt.Printf("s3: %v, len(s3): %v, cap(s3): %v\n", s3, len(s3), cap(s3))
	// 切片
	s1 = []string{"a", "b", "c", "d", "e"}
	fmt.Printf("s1: %v, len(s1): %v, cap(s1): %v\n", s1, len(s1), cap(s1))
	s2 = s1[0:0]
	fmt.Printf("s2: %v, len(s2): %v, cap(s2): %v\n", s2, len(s2), cap(s2))

	strSlice := []string{"1", "2"}
	fmt.Println(strSlice)

	// update slice in func
	updateStrMapFunc := func(m []string) {
		m[1] = "3"
	}
	updateStrMapFunc(strSlice)
	fmt.Println(strSlice)

	// append slice in func
	appendStrSliceFunc := func(m []string, ele string) {
		m = append(m, "10")
	}
	appendStrSliceFunc(strSlice, "append1")
	fmt.Println(strSlice)

	// clear slice in func
	clearStrSliceFunc := func(m []string) {
		m = nil
	}
	clearStrSliceFunc(strSlice)
	fmt.Println(strSlice)
}
