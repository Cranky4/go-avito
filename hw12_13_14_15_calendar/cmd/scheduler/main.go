package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/logger"
	schedulerinternal "github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/scheduler"
	brokeradapters "github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/scheduler/broker_adapters"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	config := NewConfig(configFile)
	logg := logger.New(config.Logger.Level)

	// broker
	var producer schedulerinternal.Producer
	var err error

	switch config.Broker.Adapter {
	case KafkaAdapter:
		producer, err = brokeradapters.NewKafkaAdapter(config.Broker.Address, logg)
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

	scheduler, err := schedulerinternal.NewScheduler(ctx, config.Database.Dsn, &producer, logg)
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

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		scheduler.Stop(ctx)
	}()

	logg.Info("scheduler is running...")

	<-ctx.Done()
}
