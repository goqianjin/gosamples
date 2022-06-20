package main

import (
	"sync"

	"github.com/spf13/cobra"
)

func newPerfCommand() (*cobra.Command, *performer, *clientArgs) {
	var perf *performer
	pArgs := &perfArgs{}
	cliArgs := &clientArgs{}
	rootCmd := &cobra.Command{
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			initLogger(pArgs.flagDebug)
		},
		Use: "soften-perf-go",
	}
	flags := rootCmd.PersistentFlags()
	// load args here
	flagLoader.loadPerfFlags(flags, pArgs, cliArgs)
	perf = newPerformer(pArgs)
	return rootCmd, perf, cliArgs
}

func newConsumerCommand(perf *performer, cliArgs *clientArgs) *cobra.Command {
	var c *consumer
	cArgs := &consumeArgs{}
	cmd := &cobra.Command{
		Use:   "consume <topic>",
		Short: "Consume from topic",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cArgs.Topic = args[0]
			c.perfConsume(perf.stopCh)
		},
	}
	flags := cmd.Flags()
	// load flags here
	flagLoader.loadConsumeFlags(flags, cArgs)
	c = newConsumer(cliArgs, cArgs)

	return cmd
}

func newProducerCommand(perf *performer, cliArgs *clientArgs) *cobra.Command {
	var p *producer
	pArgs := &produceArgs{}
	cmd := &cobra.Command{
		Use:   "produce ",
		Short: "Produce on a topic and measure performance",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			pArgs.Topic = args[0]
			p.perfProduce(perf.stopCh)
		},
	}
	flags := cmd.Flags()
	// load flags here
	flagLoader.loadProduceFlags(flags, pArgs)
	p = newProducer(cliArgs, pArgs)
	return cmd
}

func newProduceConsumeCommand(perf *performer, cliArgs *clientArgs) *cobra.Command {
	var p *producer
	var c *consumer
	pArgs := &produceArgs{}
	cArgs := &consumeArgs{}
	cmd := &cobra.Command{
		Use:   "produce-consume ",
		Short: "Both produce and consume on a topic and measure performance",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			group := sync.WaitGroup{}
			group.Add(2)
			// produce
			go func() {
				pArgs.Topic = args[0]
				p.perfProduce(perf.stopCh)
				group.Done()
			}()
			// consume
			go func() {
				cArgs.Topic = args[0]
				c.perfConsume(perf.stopCh)
				group.Done()
			}()
			group.Wait()
		},
	}
	flags := cmd.Flags()
	// load flags here
	flagLoader.loadProduceFlags(flags, pArgs)
	flagLoader.loadConsumeFlags(flags, cArgs)
	p = newProducer(cliArgs, pArgs)
	c = newConsumer(cliArgs, cArgs)
	return cmd
}
