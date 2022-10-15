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
	schedulerinternal "github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/scheduler"
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
	logg := logger.New(config.Logger.Level, log.LstdFlags)

	var adapter schedulerinternal.Adapter
	var err error

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

	scheduler, err := schedulerinternal.NewScheduler(ctx, config, &adapter, logg)
	if err != nil {
		logg.Error("unknown broker adapter")
		cancel()
		os.Exit(1)
	}

	if err := scheduler.Start(); err != nil {
		logg.Error("cannot start scheduler " + err.Error())
		cancel()
		os.Exit(1)
	}

	go func() {
		<-ctx.Done()

		scheduler.Stop()
	}()

	logg.Info("scheduler is running...")

	<-ctx.Done()
	cancel()
}
