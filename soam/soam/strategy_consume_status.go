package soam

import "math/rand"

// ------ consume strategy type ------

type ConsumeStrategy string

const (
	ConsumeStrategyRandEach  = ConsumeStrategy("RandEach")
	ConsumeStrategyRandRound = ConsumeStrategy("RandRound")
	ConsumeStrategyMainFirst = ConsumeStrategy("MainFirst")
)

// ------ consume strategy ------

type consumeStrategy interface {
	NextConsumer() int
}

// ------ RandEach ------

type randEachConsumeStrategy struct {
	total      uint
	prefixSums []uint
}

func (s *randEachConsumeStrategy) NextConsumer() int {
	r := uint(rand.Intn(int(s.total)))
	for index, prefixSum := range s.prefixSums {
		if r >= prefixSum {
			return index
		}
	}
	return 0
}

// ------ RandRound ------

type randRoundConsumeStrategy struct {
	total            int
	indexConsumerMap map[int]int
	nextIndexes      []int
}

func (s *randRoundConsumeStrategy) NextConsumer() int {
	if len(s.nextIndexes) == 0 {
		s.nextIndexes = rand.Perm(s.total)
	}
	next := s.nextIndexes[0]
	s.nextIndexes = s.nextIndexes[1:]
	return s.indexConsumerMap[next]
}
