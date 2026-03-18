package app

import (
	"context"

	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/common"
	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  common.Logger
	storage storage.EventStorage
}

func New(logger common.Logger, storage storage.EventStorage) *App {
	return &App{logger, storage}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
