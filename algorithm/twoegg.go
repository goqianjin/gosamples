package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println(twoEggDrop(100))
}

func twoEggDrop(n int) int {
	var dp [2][]int
	dp[0] = make([]int, n+1)
	dp[1] = make([]int, n+1)
	// 1个蛋, i层楼
	for j := 1; j <= n; j++ {
		dp[0][j] = j
	}
	//
	for j := 1; j <= n; j++ {
		dp[1][j] = math.MaxInt
		for k := 1; k <= j; k++ {
			dp[1][j] = Min(dp[1][j], Max(dp[0][k-1]+1, dp[1][j-k]+1))
		}
	}
	return dp[1][n]
}

func Max(v1, v2 int) int {
	if v1 > v2 {
		return v1
	} else {
		return v2
	}
}
func Min(v1, v2 int) int {
	if v1 < v2 {
		return v1
	} else {
		return v2
	}
}
