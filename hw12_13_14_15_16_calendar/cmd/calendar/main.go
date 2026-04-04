package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/app"                          //nolint:depguard
	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/logger"                       //nolint:depguard
	internalgrpc "github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/server/grpc"     //nolint:depguard
	internalhttp "github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/server/http"     //nolint:depguard
	eventstorage "github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/storage"         //nolint:depguard
	memorystorage "github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/storage/memory" //nolint:depguard
	sqlstorage "github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/storage/sql"       //nolint:depguard
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

	httpServer := internalhttp.NewServer(logg, calendar, config.Server.Addr)
	grpcServer := internalgrpc.NewServer(logg, calendar, config.GRPC.Addr)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		if err := httpServer.Start(ctx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
		}
	}()

	go func() {
		if err := grpcServer.Start(ctx); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	<-ctx.Done()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer shutdownCancel()

	if err := httpServer.Stop(shutdownCtx); err != nil {
		logg.Error("failed to stop http server: " + err.Error())
	}
	grpcServer.Stop()
}
