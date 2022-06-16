package strategy

import (
	"errors"
	"fmt"
	"math/rand"
)

// ------ RoundRand ------

type roundRandStrategy struct {
	total      uint
	indexesMap map[int]int // index到owner的下标映射

	nextIndexes []int
}

func NewRoundRandStrategy(weights []uint) (*roundRandStrategy, error) {
	total := uint(0)
	count := 0
	indexesMap := make(map[int]int)
	for index, weight := range weights {
		if weight > 100 {
			return nil, errors.New(fmt.Sprintf("invalid weight: %d. weight cannot exceed 100", weight))
		}
		for i := 0; i < int(weight); i++ {
			indexesMap[count] = index
			count++
		}
		total += weight
	}
	return &roundRandStrategy{total: total, indexesMap: indexesMap}, nil
}

func (s *roundRandStrategy) Next() int {
	if len(s.nextIndexes) == 0 {
		s.nextIndexes = rand.Perm(int(s.total))
	}
	next := s.nextIndexes[0]
	s.nextIndexes = s.nextIndexes[1:]
	return s.indexesMap[next]
}
