package main

import (
	"context"
	"log"

	"github.com/shenqianjin/soften-client-go/soften/checker"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/shenqianjin/soften-client-go/soften"
	"github.com/shenqianjin/soften-client-go/soften/config"
)

func main() {
	cli, err := soften.NewClient(config.ClientConfig{URL: "pulsar://localhost:6650"})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	consumer, err := cli.CreateListener(config.ConsumerConfig{})
	if err != nil {
		log.Fatal(err)
	}
	err = consumer.Start(context.Background(), handleBiz,
		checker.PreBlockingChecker(checkQuota), checker.PrePendingChecker(checkRate))
	if err != nil {
		log.Fatal(err)
	}
}

func handleBiz(message pulsar.Message) (bool, error) {
	return true, nil
}

func checkRate(message pulsar.Message) checker.CheckStatus {
	return checker.CheckStatusPassed
}

func checkQuota(message pulsar.Message) checker.CheckStatus {
	return checker.CheckStatusPassed
}
