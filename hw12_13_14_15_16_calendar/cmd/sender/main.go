package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/logger"   //nolint:depguard
	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/rabbitmq" //nolint:depguard
	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/storage"  //nolint:depguard
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config := NewConfig(configFile)
	logg := logger.New(config.Logger.Level)

	queueClient, err := rabbitmq.NewQueueClient(
		config.RabbitMQ.Addr,
		config.RabbitMQ.Username,
		config.RabbitMQ.Password,
		config.RabbitMQ.Queue,
	)
	if err != nil {
		logg.Error("failed to create RabbitMQ consumer: " + err.Error())
		os.Exit(1)
	}
	defer queueClient.Close()

	msgs, err := queueClient.Consume()
	if err != nil {
		logg.Error("failed to start consuming: " + err.Error())
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logg.Info("sender started, waiting for messages...")

	for {
		select {
		case <-ctx.Done():
			logg.Info("sender shutting down")
			return
		case msg := <-msgs:
			var ev storage.Event
			if err := json.Unmarshal(msg.Body, &ev); err != nil {
				logg.Error("failed to unmarshal event: " + err.Error())
				continue
			}
			logg.Info(fmt.Sprintf("received event: %+v", ev))

			if err := queueClient.PublishTo(ctx, "sent_events", ev); err != nil {
				logg.Error("failed to publish event: " + err.Error())
				continue
			}
		}
	}
}
