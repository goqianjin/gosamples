package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shenqianjin/soften-client-go/soften"
	"github.com/shenqianjin/soften-client-go/soften/checker"
	"github.com/shenqianjin/soften-client-go/soften/config"

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
	cli, err := soften.NewClient(config.ClientConfig{URL: "pulsar://localhost:6650"})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	// start listening
	listener, err := cli.CreateListener(config.ConsumerConfig{Topics: []string{"XXX-JOB", "XXX-JOB-L1", "XXX-JOB-S1"}})
	if err != nil {
		log.Fatal(err)
	}
	err = listener.StartPremium(context.Background(), svc.handleMessage,
		checker.PrePendingChecker(func(message pulsar.Message) checker.CheckStatus {
			return checker.CheckStatusPassed
		}))
	if err != nil {
		log.Fatal(err)
	}

	return svc
}

func (s *service) handleMessage(message pulsar.Message) soften.HandleStatus {
	if message.RedeliveryCount() == 1 {
		return soften.HandleStatusOk
	} else if 1 < 3 {
		return soften.HandleStatusFail.GotoAction("message.MessageStatusRetrying")
	}
	return soften.HandleStatusFail.GotoAction("soam.MessageStatusRetrying")
}

func (s *service) checkRate(message pulsar.Message) bool {
	return true
}

func (s *service) checkQuota(message pulsar.Message) bool {
	return true
}
