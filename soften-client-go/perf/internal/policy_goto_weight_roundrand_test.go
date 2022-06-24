package internal

import (
	"log"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/shenqianjin/soften-client-go/soften/message"
)

func TestRoundRandWeightGotoPolicy(t *testing.T) {
	weightMap := make(map[string]uint)
	weightMap[string(message.GotoDone)] = 19
	weightMap[string(message.GotoRetrying)] = 5
	weightMap[string(message.GotoPending)] = 5
	weightMap[string(message.GotoBlocking)] = 5
	weightMap[string(message.GotoDiscard)] = 1

	chooseMap := make(map[string]int)
	loop := 500
	policy := NewRoundRandWeightGotoPolicy(weightMap)
	for i := 0; i < loop; i++ {
		next := policy.Next()
		chooseMap[next.(string)]++
		assert.True(t, next != nil)
	}
	for status, weight := range weightMap {
		expectedRate := float64(weight) / float64(policy.total)
		chooseRate := float64(chooseMap[status]) / float64(loop)
		log.Printf("weighted round rand goto policy - expect rate: %v, cost rate: %v", expectedRate, chooseRate)
		assert.True(t, math.Abs(expectedRate-chooseRate) < expectedRate*0.1)
	}

}
