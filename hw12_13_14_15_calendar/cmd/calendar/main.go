package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/app"
	"github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	config := NewConfig(configFile)
	logg := logger.New(config.Logger.Level)

	var storage app.Storage
	if config.Storage.Driver == SQLStorage {
		storage = sqlstorage.New(ctx, config.Database.Dsn)
		s, ok := storage.(*sqlstorage.Storage)
		if ok {
			s.Connect(ctx)
		}
	} else {
		storage = memorystorage.New()
	}

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar, config.HTTP.Addr, config.HTTP.RequestLogFile)
	server.Start(ctx)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		// отключение вэб сервера
		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}

		// отключение от базы данных
		s, ok := storage.(*sqlstorage.Storage)
		if ok {
			s.Close(ctx)
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
