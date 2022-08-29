package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	internalbroker "github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/broker"
	"github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/logger"
	internalsender "github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/sender"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	config := NewConfig(configFile)

	var adapter internalsender.Adapter
	var err error

	logg := logger.New(config.Logger.Level, log.LstdFlags)

	switch config.Broker.Adapter {
	case internalbroker.KafkaBrokerAdapter:
		adapter, err = internalbroker.NewKafkaAdapter(config.Broker, logg)
		if err != nil {
			logg.Error("failed to start broker adapter: " + err.Error())
			cancel()
			os.Exit(1)
		}
	default:
		logg.Error("unknown broker adapter")
		cancel()
		os.Exit(1)
	}

	sender := internalsender.NewSender(ctx, config, &adapter, logg)

	if err := sender.Start(); err != nil {
		logg.Error("failed to start sender: " + err.Error())
		cancel()
		os.Exit(1)
	}

	go func() {
		<-ctx.Done()

		if err := sender.Stop(); err != nil {
			logg.Error("failed to stop sender: " + err.Error())
			cancel()
			os.Exit(1)
		}
	}()

	logg.Info("sender is running...")

	<-ctx.Done()
	cancel()
}
