package app

import (
	"time"

	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/common"  //nolint:depguard
	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
)

type App struct {
	logger  common.Logger
	storage storage.EventStorage
}

func New(logger common.Logger, storage storage.EventStorage) *App {
	return &App{logger, storage}
}

func (a *App) CreateEvent(ev storage.Event) error {
	return a.storage.Create(ev)
}

func (a *App) UpdateEvent(ev storage.Event) error {
	return a.storage.Update(ev)
}

func (a *App) DeleteEvent(id string) error {
	return a.storage.Delete(id)
}

func (a *App) ListEventsForDay(date time.Time) ([]storage.Event, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, date.Location())
	return a.storage.List(start, end)
}

func (a *App) ListEventsForWeek(date time.Time) ([]storage.Event, error) {
	start, end := getWeekBounds(date)
	return a.storage.List(start, end)
}

func (a *App) ListEventsForMonth(start time.Time) ([]storage.Event, error) {
	start, end := getMonthBounds(start)
	return a.storage.List(start, end)
}

func (a *App) ListForNotification() ([]storage.Event, error) {
	start := time.Now()
	end := start.Add(time.Minute)
	return a.storage.ListForNotification(start, end)
}

func (a *App) DeleteOutdated() error {
	olderThan := time.Now().AddDate(-1, 0, 0)
	return a.storage.DeleteOutdated(olderThan)
}

func getWeekBounds(t time.Time) (start, end time.Time) {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	start = time.Date(t.Year(), t.Month(), t.Day()-weekday+1, 0, 0, 0, 0, t.Location())
	end = start.AddDate(0, 0, 7).Add(-time.Nanosecond)
	return
}

func getMonthBounds(t time.Time) (start, end time.Time) {
	start = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	end = start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	return
}
