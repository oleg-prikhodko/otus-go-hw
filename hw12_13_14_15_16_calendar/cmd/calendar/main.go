package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/server/http"
	eventstorage "github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig(configFile)
	logg := logger.New(config.Logger.Level)

	var storage eventstorage.EventStorage
	switch config.Storage.Type {
	case Memory:
		storage = memorystorage.New()
	case SQL:
		s := sqlstorage.New(logg, config.Storage.Addr)
		if err := s.Connect(); err != nil {
			logg.Error("failed to connect to db: " + err.Error())
			os.Exit(1)
		}
	default:
		logg.Error("unknown storage type: " + string(config.Storage.Type))
		os.Exit(1)
	}
	defer storage.Close()

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar, config.Server.Addr)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
