package main

import (
	"fmt"
	"log"
	"soam/soam"

	"github.com/apache/pulsar-client-go/pulsar"
)

func main() {
	svc := NewService()
	fmt.Println(svc)
}

type service struct {
}

func NewService() *service {
	// initialize
	svc := &service{}

	// new client
	cli, err := soam.NewClient(soam.ClientConfig{URL: "pulsar://localhost:6650"})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	// start listening
	err = cli.SubscribeInPremium(
		soam.ComsumerConfig{Topics: []string{"XXX-JOB", "XXX-JOB-L1", "XXX-JOB-S1"}},
		svc.handleMessage,
		soam.PrePendingChecker(func(message pulsar.Message) (passed bool) {
			return true
		}))
	if err != nil {
		log.Fatal(err)
	}

	return svc
}

func (s *service) handleMessage(message pulsar.Message) soam.HandledStatus {
	if message.RedeliveryCount() == 1 {
		return soam.HandleStatusDone
	} else if 1 < 3 {
		return soam.HandleStatusRetry
	}
	return soam.HandleStatusRetry
}

func (s *service) checkRate(message pulsar.Message) bool {
	return true
}

func (s *service) checkQuota(message pulsar.Message) bool {
	return true
}
