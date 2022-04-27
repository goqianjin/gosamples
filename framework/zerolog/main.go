package main

import (
	"flag"
	"fmt"
	"os"
	"time"
	"zerolog/logger"

	"github.com/rs/zerolog/diode"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	debug := flag.Bool("debug", false, "sets log level to debug")

	flag.Parse()

	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Debug().Msg("This message appears only when log level set to Debug")
	log.Info().Msg("This message appears when log level set to Debug or Info")

	hn, err := os.Hostname()
	log.Info().Msg(fmt.Sprintf("hostname: %s, err: %v", hn, err))
	if e := log.Debug(); e.Enabled() {
		// Compute log output only if enabled.
		value := "bar"
		e.Str("foo", value).Msg("some debug message")
	}

	// 输出多目标
	logger.InitLog()
	logger.Logger.Info().Msg(fmt.Sprintf("------hostname: %s, err: %v", hn, err))

	// 异步日志
	wr := diode.NewWriter(os.Stdout, 1000, 10*time.Millisecond, func(missed int) {
		fmt.Printf("Logger Dropped %d messages", missed)
	})
	log2 := zerolog.New(wr)
	//log.Print("test")
	log2.Info().Msg(fmt.Sprint   f("*** async *** hostname: %s, err: %v", hn, err))
	log.Info().Msg("--sync after async---")
	time.Sleep(time.Second * 2)
}
