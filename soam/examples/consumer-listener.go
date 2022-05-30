package main

import (
	"log"
	"soam/soam"

	"github.com/apache/pulsar-client-go/pulsar"
)

func main() {
	cli, err := soam.NewClient(soam.ClientConfig{URL: "pulsar://localhost:6650"})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	err = cli.SubscribeInRegular(soam.ComsumerConfig{}, handleBiz,
		soam.PreBlockingChecker(checkQuota), soam.PrePendingChecker(checkRate))
	if err != nil {
		log.Fatal(err)
	}
}

func handleBiz(message pulsar.Message) bool {
	return true
}

func checkRate(message pulsar.Message) bool {
	return true
}

func checkQuota(message pulsar.Message) bool {
	return true
}
