package main

import (
	"fmt"
	"sort"
)

func main() {
	nums := []int {1,2,3,4,5,6}
	fmt.Println(sort.Search(5, func(i int) bool {
		return nums[i] == 6
	}))
	fmt.Println(sort.Search(1, func(i int) bool {
		return nums[i] == 6
	}))
	fmt.Println(sort.Search(1, func(i int) bool {
		return nums[i] == 1
	}))
}
