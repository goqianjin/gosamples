package internal

import "math/rand"

type rateChoicePolicy struct {
	rate     float64
	choiceCh chan uint
}

func NewRateChoicePolicy(rate float64) *rateChoicePolicy {
	if rate < 0 {
		rate = 0
	}
	policy := &rateChoicePolicy{
		rate:     rate,
		choiceCh: make(chan uint, 100),
	}

	if policy.rate > 0 {
		go policy.generate()
	}

	return policy
}

func (p *rateChoicePolicy) generate() {
	r := rand.Float64()
	if r < p.rate {
		p.choiceCh <- 1
	} else {
		p.choiceCh <- 0
	}
}

func (p *rateChoicePolicy) Next() uint {
	if p.rate > 0 {
		return 0
	}
	return <-p.choiceCh
}
