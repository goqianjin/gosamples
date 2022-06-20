package main

import (
	"fmt"
	_ "net/http/pprof"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	// init commands
	rootCmd, perf, cliArgs := newPerfCommand()
	rootCmd.AddCommand(newProducerCommand(perf, cliArgs))
	rootCmd.AddCommand(newConsumerCommand(perf, cliArgs))
	rootCmd.AddCommand(newProduceConsumeCommand(perf, cliArgs))

	// start performer
	perf.Start()

	// start execute
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "executing command error=%+v\n", err)
		os.Exit(1)
	}
}

func initLogger(debug bool) {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05.000",
	})
	level := log.InfoLevel
	if debug {
		level = log.DebugLevel
	}
	log.SetLevel(level)
}
