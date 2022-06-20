package internal

import (
	"math/rand"
	"time"
)

type CostPolicy interface {
	Next() time.Duration
}

type avgCostPolicy struct {
	avg            int64
	positiveJitter float64
	negativeJitter float64

	costCh chan int64
}

func NewAvgCostPolicy(avg int64, positiveJitter, negativeJitter float64) *avgCostPolicy {
	if avg < 0 {
		avg = 0
	}
	if positiveJitter < 0 {
		positiveJitter = 0
	}
	if negativeJitter < 0 {
		negativeJitter = 0
	}
	policy := &avgCostPolicy{
		avg:            avg,
		positiveJitter: positiveJitter,
		negativeJitter: negativeJitter,
		costCh:         make(chan int64, 100),
	}

	go policy.generate()

	return policy
}

func (p *avgCostPolicy) Next() time.Duration {
	cost := <-p.costCh
	return time.Duration(cost) * time.Millisecond
}

func (p *avgCostPolicy) generate() {
	base := float64(p.avg) * (1 - p.negativeJitter)
	jitter := rand.Float64() * (p.positiveJitter + p.negativeJitter)
	cost := base + jitter
	if cost < 0 {
		cost = 0
	}
	p.costCh <- int64(cost)
}
