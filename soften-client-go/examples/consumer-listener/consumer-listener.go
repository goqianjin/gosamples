package main

import (
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
	_, err = cli.SubscribeRegular(config.ConsumerConfig{}, handleBiz,
		checker.PreBlockingChecker(checkQuota), checker.PrePendingChecker(checkRate))
	if err != nil {
		log.Fatal(err)
	}
}

func handleBiz(message pulsar.Message) (bool, error) {
	return true, nil
}

func checkRate(message pulsar.Message) bool {
	return true
}

func checkQuota(message pulsar.Message) bool {
	return true
}
