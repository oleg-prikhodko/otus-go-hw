package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/app"                          //nolint:depguard
	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/common"                       //nolint:depguard
	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/logger"                       //nolint:depguard
	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/rabbitmq"                     //nolint:depguard
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

	queueClient, err := rabbitmq.NewQueueClient(
		config.RabbitMQ.Addr,
		config.RabbitMQ.Username,
		config.RabbitMQ.Password,
		config.RabbitMQ.Queue,
	)
	if err != nil {
		logg.Error("failed to create RabbitMQ publisher: " + err.Error())
		return
	}
	defer queueClient.Close()

	calendar := app.New(logg, storage)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	logg.Info("scheduler started")

	for {
		select {
		case <-ticker.C:
			sendNotifications(ctx, calendar, queueClient, logg)
			deleteOutdated(calendar, logg)
		case <-ctx.Done():
			logg.Info("scheduler shutting down")
			return
		}
	}
}

func sendNotifications(ctx context.Context, calendar *app.App, queueClient *rabbitmq.QueueClient, logg common.Logger) {
	events, err := calendar.ListForNotification()
	if err != nil {
		logg.Error("ListForNotification error: " + err.Error())
		return
	}

	if len(events) == 0 {
		logg.Info("no events to notify")
		return
	}

	for _, ev := range events {
		eventJSON, err := json.Marshal(ev)
		if err != nil {
			logg.Error("failed to marshal event: " + err.Error())
			continue
		}
		logg.Info(fmt.Sprintf("sending event to RabbitMQ: %s", string(eventJSON)))

		if err := queueClient.Publish(ctx, ev); err != nil {
			logg.Error("failed to publish event to RabbitMQ: " + err.Error())
		}
	}
}

func deleteOutdated(calendar *app.App, logg common.Logger) {
	if err := calendar.DeleteOutdated(); err != nil {
		logg.Error("DeleteOutdated error: " + err.Error())
	}
}
