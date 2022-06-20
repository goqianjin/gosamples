package internal

import (
	"math/rand"
)

type GotoPolicy interface {
	Next() interface{}
}

type roundRandWeightGotoPolicy struct {
	weights  []uint
	choiceCh chan interface{}

	total      uint
	indexesMap map[int]interface{} // index到owner的下标映射
}

func NewRoundRandWeightGotoPolicy(weightMap map[string]uint) *roundRandWeightGotoPolicy {
	total := uint(0)
	count := 0
	indexesMap := make(map[int]interface{})
	weights := make([]uint, len(weightMap))
	for key, weight := range weights {
		for i := 0; i < int(weight); i++ {
			indexesMap[count] = key
			count++
		}
		total += weight
	}

	policy := &roundRandWeightGotoPolicy{
		weights:    weights,
		choiceCh:   make(chan interface{}, 100),
		total:      total,
		indexesMap: indexesMap,
	}

	if len(policy.weights) > 1 {
		go policy.generate()
	}

	return policy
}

func (p *roundRandWeightGotoPolicy) generate() {
	nextIndexes := rand.Perm(int(p.total))
	for _, index := range nextIndexes {
		p.choiceCh <- p.indexesMap[index]
	}
}

func (p *roundRandWeightGotoPolicy) Next() interface{} {
	if len(p.weights) > 1 {
		return <-p.choiceCh
	}
	return 0
}
