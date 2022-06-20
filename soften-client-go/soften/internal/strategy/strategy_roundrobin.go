package strategy

import (
	"errors"
	"fmt"
)

// ------ RoundRobin ------

type roundRobinStrategy struct {
	indexes []int
	curr    int
}

func NewRoundRobinStrategy(destLen int) (*roundRobinStrategy, error) {
	if destLen <= 0 {
		return nil, errors.New(fmt.Sprintf("invalid destLen: %d", destLen))
	}
	indexes := make([]int, destLen)
	for index := 0; index < destLen; index++ {
		indexes[index] = index
	}
	return &roundRobinStrategy{indexes: indexes, curr: 0}, nil
}

func (s *roundRobinStrategy) Next(excludes ...int) int {
	next := s.indexes[s.curr]
	s.curr = (s.curr + 1) % len(s.indexes)
	return next
}
