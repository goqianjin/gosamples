package soam

import (
	"errors"
	"strconv"
)

// ------ abbr backoff policy ------

type abbrBackoffPolicy struct {
	backoffDelays []uint
}

func (p abbrBackoffPolicy) Next(redeliveryCount int) uint {
	if redeliveryCount < 0 {
		redeliveryCount = 0
	}
	if redeliveryCount >= len(p.backoffDelays) {
		return p.backoffDelays[len(p.backoffDelays)-1]
	}
	return p.backoffDelays[redeliveryCount]
}

func newAbbrBackoffPolicy(delays []string) (BackoffPolicy, error) {
	if len(delays) == 0 {
		return nil, errors.New("backoffDelays is empty")
	}
	backoffDelays := make([]uint, len(delays))
	for _, delay := range delays {
		last := delay[len(delay)-1]
		if unit, ok := DelayUnitMap[string(last)]; !ok {
			return nil, errors.New("invalid unit in backOffDelays")
		} else if d, err := strconv.Atoi(delay[0 : len(delay)-1]); err != nil {
			return nil, errors.New("invalid in backOffDelays")
		} else {
			backoffDelays = append(backoffDelays, uint(d*unit))
		}
	}
	return abbrBackoffPolicy{backoffDelays: backoffDelays}, nil
}

// ------ abbr status backoff policy ------

type abbrStatusBackoffPolicy struct {
	backoffPolicy BackoffPolicy
}

func (p abbrStatusBackoffPolicy) Next(statusReconsumeTimes int) uint {
	return p.backoffPolicy.Next(statusReconsumeTimes)
}

func newAbbrStatusBackoffPolicy(delays []string) (StatusBackoffPolicy, error) {
	backoffPolicy, err := newAbbrBackoffPolicy(delays)
	if err != nil {
		return nil, err
	}
	return abbrStatusBackoffPolicy{backoffPolicy: backoffPolicy}, nil
}
