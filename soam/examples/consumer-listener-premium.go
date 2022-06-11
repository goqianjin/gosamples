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
	err = cli.SubscribePremium(
		soam.ConsumerConfig{Topics: []string{"XXX-JOB", "XXX-JOB-L1", "XXX-JOB-S1"}},
		svc.handleMessage,
		soam.PrePendingChecker(func(message pulsar.Message) (passed bool) {
			return false
		}))
	if err != nil {
		log.Fatal(err)
	}

	return svc
}

func (s *service) handleMessage(message pulsar.Message) soam.HandleStatus {
	if message.RedeliveryCount() == 1 {
		return soam.HandleStatusOk
	} else if 1 < 3 {
		return soam.HandleStatusFail.TransferTo(soam.MessageStatusRetrying)
	}
	return soam.HandleStatusFail.TransferTo(soam.MessageStatusRetrying)
}

func (s *service) checkRate(message pulsar.Message) bool {
	return true
}

func (s *service) checkQuota(message pulsar.Message) bool {
	return true
}
