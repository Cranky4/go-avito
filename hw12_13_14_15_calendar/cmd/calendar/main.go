package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/app"
	"github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/server/grpc"
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

	config := NewConfig(configFile)
	logg := logger.New(config.Logger.Level, log.LstdFlags)

	var storage app.Storage
	if config.Storage.Driver == SQLStorage {
		storage = sqlstorage.New(ctx, config.Database, logg)
		s, ok := storage.(*sqlstorage.Storage)
		if ok {
			err := s.Connect(ctx)
			if err != nil {
				logg.Error(err.Error())
				cancel()
				os.Exit(1)
			}
		}
	} else {
		storage = memorystorage.New()
	}

	calendar := app.New(logg, storage)

	var server *internalhttp.Server
	go func() {
		server = internalhttp.NewServer(logg, calendar, config.HTTP.Addr, config.HTTP.RequestLogFile)
		if err := server.Start(ctx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
			os.Exit(1)
		}
	}()

	var grpcServer *internalgrpc.Server
	go func() {
		grpcServer = internalgrpc.New(calendar, logg, config.GRPC.RequestLogFile)
		if err := grpcServer.Start(ctx, config.GRPC.Addr); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
			cancel()
			os.Exit(1)
		}
	}()
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		// отключение вэб сервера
		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}

		// отключение grpc сервера
		grpcServer.Stop(ctx)

		// отключение от базы данных
		s, ok := storage.(*sqlstorage.Storage)
		if ok {
			s.Close(ctx)
		}
	}()

	logg.Info("calendar is running...")

	<-ctx.Done()
}
