package main

import (
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/pflag"
)

var flagLoader = &flagLoaderImpl{}

type flagLoaderImpl struct {
	perfFlagLoadOnce    sync.Once
	consumeFlagLoadOnce sync.Once
	produceFlagLoadOnce sync.Once
}

func (l *flagLoaderImpl) loadPerfFlags(flags *pflag.FlagSet, pArgs *perfArgs, cliArgs *clientArgs) {
	flags.IntVar(&pArgs.profilePort, "profilePort", 6060, "Port to expose profiling info, use -1 to disable")
	flags.IntVar(&pArgs.PrometheusPort, "metrics", 8000, "Port to use to export metrics for Prometheus. Use -1 to disable.")
	flags.BoolVar(&pArgs.flagDebug, "debug", false, "enable debug output")

	flags.StringVarP(&cliArgs.ServiceURL, "service-url", "u", "pulsar://localhost:6650", "The Pulsar service URL")
	flags.StringVar(&cliArgs.TokenFile, "token-file", "", "file path to the Pulsar JWT file")
	flags.StringVar(&cliArgs.TLSTrustCertFile, "trust-cert-file", "", "file path to the trusted certificate file")
}

func (l *flagLoaderImpl) loadConsumeFlags(flags *pflag.FlagSet, cArgs *consumeArgs) {
	flags.StringVarP(&cArgs.SubscriptionName, "subscription", "S", "sub", "Subscription name")
	flags.IntVarP(&cArgs.ReceiverQueueSize, "receiver-queue-size", "R", 100000, "Receiver queue size")

	flags.Int64Var(&cArgs.costAverageInMs, "cost-average-in-ms", 0, "consume cost average in milliseconds, use -1 to disable")
	flags.Float64Var(&cArgs.costPositiveJitter, "cost-delay-jitter", 1, "cost delay jitter rate")
	flags.Float64Var(&cArgs.costNegativeJitter, "cost-precede-jitter", 1, "cost precede jitter rate")

	var limitStr string
	flags.StringVarP(&limitStr, "limit", "l", "10,10", "Consume limit. Set to 0 to go un-throttled, append with ',' for radicals Limits")
	cArgs.Limits = l.parseUint64Array(limitStr)
	var concurrencyStr string
	flags.StringVarP(&concurrencyStr, "concurrency", "c", "3,3", "Consume concurrency. Set to 0 to go un-throttled, append with ',' for radicals Limits")
	cArgs.RadicalConcurrences = l.parseUint64Array(concurrencyStr)

	flags.UintVar(&cArgs.handleDoneWeight, "done-weight", 100, "weight of handling message as done status")
	flags.UintVar(&cArgs.handleRetryingWeight, "retrying-weight", 60, "weight of handling message as retrying status")
	flags.UintVar(&cArgs.handlePendingWeight, "pending-weight", 30, "weight of handling message as pending status")
	flags.UintVar(&cArgs.handleBlockingWeight, "blocking-weight", 10, "weight of handling message as blocking status")
	flags.UintVar(&cArgs.handleDeadWeight, "dead-weight", 5, "weight of handling message as dead status")
	flags.UintVar(&cArgs.handleDiscardWeight, "discard-weight", 5, "weight of handling message as discard status")

	if cArgs.costAverageInMs < 0 {
		cArgs.costAverageInMs = 0
	}
	if cArgs.costNegativeJitter > 1 {
		cArgs.costNegativeJitter = 1
	}
	if cArgs.handleDoneWeight <= 0 {
		cArgs.handleDoneWeight = 1
	}

}

func (l *flagLoaderImpl) loadProduceFlags(flags *pflag.FlagSet, pArgs *produceArgs) {
	var rateStr string
	flags.StringVarP(&rateStr, "rate", "r", "80,20", "Publish rate. Set to 0 to go un-throttled, append with ',' for radicals rates")
	pArgs.Rates = l.parseUint64Array(rateStr)
	flags.IntVarP(&pArgs.BatchingTimeMillis, "batching-time", "b", 1, "Batching grouping time in millis")
	flags.IntVarP(&pArgs.BatchingMaxSize, "batching-max-size", "", 128, "Max size of a batch (in KB)")
	flags.IntVarP(&pArgs.MessageSize, "msg-size", "s", 1024, "Message size (int B)")
	flags.IntVarP(&pArgs.ProducerQueueSize, "queue-size", "q", 1000, "Produce queue size")
	flags.Float64Var(&pArgs.RadicalRate, "radical-rate", 0.8, "the rate of radical rate")

}

func (l *flagLoaderImpl) parseUint64Array(str string) []uint64 {
	nums := make([]uint64, 0)
	if str != "" {
		numsStr := strings.Split(str, ",")
		for _, ns := range numsStr {
			if n, err := strconv.ParseUint(ns, 10, 32); err != nil {
				panic(err)
			} else {
				nums = append(nums, n)
			}
		}
	}
	return nums
}
