package config

import (
	"errors"
	"fmt"

	"github.com/shenqianjin/soften-client-go/soften/topic"

	"github.com/apache/pulsar-client-go/pulsar/log"
	"github.com/shenqianjin/soften-client-go/soften/internal"
	"github.com/shenqianjin/soften-client-go/soften/internal/backoff"
	"github.com/sirupsen/logrus"
)

// ------ configuration validator ------

var Validator = &validator{}

type validator struct {
}

func (v *validator) ValidateAndDefaultClientConfig(conf *ClientConfig) error {
	if conf.Logger == nil {
		//conf.Logger = sogrus.NewLoggerWithLogrus(logrus.StandardLogger())
		conf.Logger = log.NewLoggerWithLogrus(logrus.StandardLogger())

	}
	return nil
}

func (v *validator) ValidateAndDefaultConsumerConfig(conf *ConsumerConfig) error {
	// default Level
	if len(conf.Levels) == 0 && conf.Level == "" {
		conf.Levels = []internal.TopicLevel{topic.L1}
	} else if conf.Level != "" {
		conf.Levels = []internal.TopicLevel{conf.Level}
	}
	if len(conf.Levels) > 0 && conf.Level != "" {
		if conf.Levels[0] != conf.Level {
			return errors.New("core level is not match between level and levels configuration")
		}
	}
	// default leveled balance
	if conf.LevelBalanceStrategy == "" {
		conf.LevelBalanceStrategy = BalanceStrategyRoundRand
	}
	// default status balance strategy
	if conf.BalanceStrategy == "" {
		conf.BalanceStrategy = BalanceStrategyRoundRand
	}

	// default topics
	if len(conf.Topics) == 0 && conf.Topic == "" {
		return errors.New("no topic found in your configuration")
	} else if conf.Topic != "" {
		conf.Topics = []string{conf.Topic}
	}
	if len(conf.Topics) >= 1 && conf.Topic != "" {
		if conf.Topics[0] != conf.Topic {
			return errors.New("core topic is not match between topic and topics configuration")
		}
	}
	// default leveled policy when consume more than one level
	if len(conf.Levels) > 1 {
		if err := v.validateAndDefaultLeveledPolicy(conf.Levels, &conf.LevelPolicies, defaultLeveledPolicy); err != nil {
			return err
		}
	}
	// default status Policy
	if conf.Ready == nil {
		conf.Ready = defaultStatusPolicyReady
	} else {
		if err := v.validateAndDefaultStatusPolicy(conf.Ready, defaultStatusPolicyReady); err != nil {
			return err
		}
	}
	if conf.PendingEnable {
		if conf.Pending == nil {
			conf.Pending = defaultStatusPolicyPending
		} else if err := v.validateAndDefaultStatusPolicy(conf.Pending, defaultStatusPolicyPending); err != nil {
			return err
		}
	}
	if conf.BlockingEnable {
		if conf.Blocking == nil {
			conf.Blocking = defaultStatusPolicyBlocking
		} else if err := v.validateAndDefaultStatusPolicy(conf.Blocking, defaultStatusPolicyBlocking); err != nil {
			return err
		}
	}
	if conf.RetryingEnable {
		if conf.Retrying == nil {
			conf.Retrying = defaultStatusPolicyRetrying
		} else if err := v.validateAndDefaultStatusPolicy(conf.Retrying, defaultStatusPolicyRetrying); err != nil {
			return err
		}
	}
	if conf.UpgradeEnable {
		if err := v.baseValidateTopicLevel(conf.UpgradeTopicLevel); err != nil {
			return nil
		}
		if conf.UpgradeTopicLevel.OrderOf() <= conf.Level.OrderOf() {
			return errors.New(fmt.Sprintf("upgrade level [%v] cannot be lower or equal than the consume level [%v]",
				conf.UpgradeTopicLevel, conf.Level))
		}
	}
	if conf.DegradeEnable {
		if err := v.baseValidateTopicLevel(conf.DegradeTopicLevel); err != nil {
			return nil
		}
		if conf.DegradeTopicLevel.OrderOf() >= conf.Level.OrderOf() {
			return errors.New(fmt.Sprintf("degrade level [%v] cannot be higher or equal than the consume level [%v]",
				conf.DegradeTopicLevel, conf.Level))
		}
	}
	// default concurrency policy
	if conf.Concurrency == nil {
		conf.Concurrency = defaultConcurrencyPolicy
	} else if err := v.validateAndDefaultConcurrencyPolicy(conf.Concurrency, defaultConcurrencyPolicy); err != nil {
		return err
	}
	return nil
}

func (v *validator) baseValidateTopicLevel(level internal.TopicLevel) error {
	if level == "" {
		return errors.New("missing upgrade/degrade TopicLevel configuration")
	}
	if !topic.Exists(level) {
		return errors.New(fmt.Sprintf("not supported topic level: %v", level))
	}
	if level.OrderOf() > topic.HighestLevel().OrderOf() {
		return errors.New("upgrade/degrade topic level is too high")
	}
	if level.OrderOf() < topic.LowestLevel().OrderOf() {
		return errors.New(fmt.Sprintf("upgrade/degrade topic level [%v] is too low", level))
	}
	return nil
}

func (v *validator) validateAndDefaultStatusPolicy(configuredPolicy *StatusPolicy, defaultPolicy *StatusPolicy) error {
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

func (v *validator) validateAndDefaultLeveledPolicy(configuredLevels []internal.TopicLevel, configuredPolicies *map[internal.TopicLevel]*LevelPolicy, defaultPolicy *LevelPolicy) error {
	if *configuredPolicies == nil {
		*configuredPolicies = make(map[internal.TopicLevel]*LevelPolicy, len(configuredLevels))
	}
	for _, level := range configuredLevels {
		configuredPolicy, ok := (*configuredPolicies)[level]
		if !ok {
			(*configuredPolicies)[level] = defaultPolicy
			continue
		}
		if configuredPolicy.ConsumeWeight == 0 {
			configuredPolicy.ConsumeWeight = defaultPolicy.ConsumeWeight
		} else if configuredPolicy.UpgradeLevel == "" {

		} else if configuredPolicy.DegradeLevel == "" {

		}
	}
	return nil
}

func (v *validator) validateAndDefaultConcurrencyPolicy(configuredPolicy *ConcurrencyPolicy, defaultPolicy *ConcurrencyPolicy) error {
	if configuredPolicy == nil {
		configuredPolicy = defaultPolicy
		return nil
	}
	if configuredPolicy.CorePoolSize <= 0 {
		configuredPolicy.CorePoolSize = defaultPolicy.CorePoolSize
	}
	if configuredPolicy.MaximumPoolSize <= 0 {
		configuredPolicy.MaximumPoolSize = defaultPolicy.MaximumPoolSize
	}
	if configuredPolicy.KeepAliveTime <= 0 {
		configuredPolicy.KeepAliveTime = defaultPolicy.KeepAliveTime
	}
	return nil
}
