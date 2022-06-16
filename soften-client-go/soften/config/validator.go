package config

import (
	"errors"
	"fmt"

	"github.com/apache/pulsar-client-go/pulsar/log"
	"github.com/sirupsen/logrus"

	"github.com/shenqianjin/soften-client-go/soften/internal"
	"github.com/shenqianjin/soften-client-go/soften/internal/backoff"
	"github.com/shenqianjin/soften-client-go/soften/topic"
)

// ------ configuration validator ------

var Validator = &confValidator{}

type confValidator struct {
}

func (v *confValidator) ValidateAndDefaultClientConfig(conf *ClientConfig) error {
	if conf.Logger == nil {
		conf.Logger = log.NewLoggerWithLogrus(logrus.StandardLogger())
	}
	return nil
}

func (v *confValidator) ValidateAndDefaultMultiLevelConsumerConfig(conf *MultiLevelConsumerConfig) error {
	if len(conf.Levels) == 0 {
		return errors.New("levels is empty")
	}
	if conf.LevelBalanceStrategy == "" {
		conf.LevelBalanceStrategy = BalanceStrategyRoundRand
	}
	return v.ValidateAndDefaultConsumerConfig(conf.ConsumerConfig)
}

func (v *confValidator) ValidateAndDefaultConsumerConfig(conf *ConsumerConfig) error {
	// default topics
	if conf.Topic == "" && len(conf.Topics) == 0 {
		return errors.New("no topic found in your configuration")
	}
	if len(conf.Topics) == 0 && conf.Topic != "" {
		conf.Topics = []string{conf.Topic}
	}
	// default Level
	if conf.Level == "" {
		conf.Level = topic.L1
	}
	// default balance strategy
	if conf.BalanceStrategy == "" {
		conf.BalanceStrategy = BalanceStrategyRoundRand
	}
	// default Policy
	if conf.Ready == nil {
		conf.Ready = defaultStatusPolicyReady
	} else {
		if err := v.validateAndDefaultStatusPolicy(conf.Ready, defaultStatusPolicyReady); err != nil {
			return err
		}
	}
	if conf.PendingEnable {
		if err := v.validateAndDefaultStatusPolicy(conf.Pending, defaultStatusPolicyPending); err != nil {
			return err
		}
	}
	if conf.BlockingEnable {
		if err := v.validateAndDefaultStatusPolicy(conf.Blocking, defaultStatusPolicyBlocking); err != nil {
			return err
		}
	}
	if conf.RetryingEnable {
		if err := v.validateAndDefaultStatusPolicy(conf.Retrying, defaultStatusPolicyRetrying); err != nil {
			return err
		}
	}
	if conf.UpgradeEnable {
		if err := v.baseValidateTopicLevel(conf.UpgradeTopicLevel); err != nil {
			return nil
		}
		if topic.OrderOf(conf.UpgradeTopicLevel) <= topic.OrderOf(conf.Level) {
			return errors.New(fmt.Sprintf("upgrade level [%v] cannot be lower or equal than the consume level [%v]",
				conf.UpgradeTopicLevel, conf.Level))
		}
	}
	if conf.DegradeEnable {
		if err := v.baseValidateTopicLevel(conf.DegradeTopicLevel); err != nil {
			return nil
		}
		if topic.OrderOf(conf.DegradeTopicLevel) >= topic.OrderOf(conf.Level) {
			return errors.New(fmt.Sprintf("degrade level [%v] cannot be higher or equal than the consume level [%v]",
				conf.DegradeTopicLevel, conf.Level))
		}
	}
	return nil
}

func (v *confValidator) baseValidateTopicLevel(l internal.TopicLevel) error {
	if l == "" {
		return errors.New("missing UpgradeTopicLevel configuration")
	}
	if !topic.Exists(l) {
		return errors.New(fmt.Sprintf("not supported topic level: %v", l))
	}
	if topic.OrderOf(l) > topic.OrderOf(topic.HighestLevel()) {
		return errors.New("upgrade topic level is too high")
	}
	if topic.OrderOf(l) < topic.OrderOf(topic.LowestLevel()) {
		return errors.New(fmt.Sprintf("upgrade topic level [%v] is too low", l))
	}
	return nil
}

func (v *confValidator) validateAndDefaultStatusPolicy(configuredPolicy *StatusPolicy, defaultPolicy *StatusPolicy) error {
	if configuredPolicy == nil {
		configuredPolicy = defaultPolicy
		return nil
	}
	if configuredPolicy.ConsumeWeight == 0 {
		configuredPolicy.ConsumeWeight = defaultPolicy.ConsumeWeight
	}
	if configuredPolicy.ConsumeMaxTimes == 0 {
		configuredPolicy.ConsumeMaxTimes = defaultPolicy.ConsumeMaxTimes
	}
	if configuredPolicy.BackoffDelays == nil && configuredPolicy.BackoffPolicy == nil {
		configuredPolicy.BackoffDelays = defaultPolicy.BackoffDelays
		configuredPolicy.BackoffPolicy = defaultPolicy.BackoffPolicy
	}
	if configuredPolicy.ReentrantDelay == 0 {
		configuredPolicy.ReentrantDelay = defaultPolicy.ReentrantDelay
	}
	if configuredPolicy.ReentrantMaxTimes == 0 {
		configuredPolicy.ReentrantMaxTimes = defaultPolicy.ReentrantMaxTimes
	}
	// default policy
	if configuredPolicy.BackoffPolicy == nil && configuredPolicy.BackoffDelays != nil {
		if backoffPolicy, err := backoff.NewAbbrStatusBackoffPolicy(configuredPolicy.BackoffDelays); err != nil {
			return err
		} else {
			configuredPolicy.BackoffDelays = nil // release unnecessary reference
			configuredPolicy.BackoffPolicy = backoffPolicy
		}
	}
	return nil
}
