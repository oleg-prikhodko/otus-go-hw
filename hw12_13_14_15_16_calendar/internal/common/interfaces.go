package common

import (
	"time"

	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
)

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Application interface {
	CreateEvent(ev storage.Event) error
	UpdateEvent(ev storage.Event) error
	DeleteEvent(id string) error
	ListEventsForDay(date time.Time) ([]storage.Event, error)
	ListEventsForWeek(date time.Time) ([]storage.Event, error)
	ListEventsForMonth(date time.Time) ([]storage.Event, error)
	ListForNotification() ([]storage.Event, error)
}
