package strategy

import (
	"errors"
	"fmt"
	"math/rand"
)

// ------ EachRand ------

type eachRandStrategy struct {
	total      uint
	prefixSums []uint
}

func NewEachRandStrategy(weights []uint) (*eachRandStrategy, error) {
	totalWeight := uint(0)
	prefixSums := make([]uint, len(weights))
	for index, weight := range weights {
		if weight > 100 {
			return nil, errors.New(fmt.Sprintf("invalid weight: %d. weight cannot exceed 100", weight))
		}
		prefixSums[index] = totalWeight
		totalWeight += weight
	}
	return &eachRandStrategy{total: totalWeight, prefixSums: prefixSums}, nil
}

func (s *eachRandStrategy) Next() int {
	r := uint(rand.Intn(int(s.total)))
	for index := 0; index < len(s.prefixSums)-1; index++ {
		if r >= s.prefixSums[index] && r < s.prefixSums[index+1] {
			return index
		}
	}
	return 0
}
