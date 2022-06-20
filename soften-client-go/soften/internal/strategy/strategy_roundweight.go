package strategy

// ------ RoundWeight ------

type roundWeightStrategy struct {
	indexes []int
	curr    int
}

func NewRoundWeightStrategy(weights []uint) (*roundWeightStrategy, error) {
	indexes := make([]int, 0)
	for index, weight := range weights {
		for i := 0; i < int(weight); i++ {
			indexes = append(indexes, index)
		}
	}
	return &roundWeightStrategy{indexes: indexes, curr: 0}, nil
}

func (s *roundWeightStrategy) Next(excludes ...int) int {
	next := s.indexes[s.curr]
	s.curr = (s.curr + 1) % len(s.indexes)
	return next
}
