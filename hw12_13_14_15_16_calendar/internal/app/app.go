package app

import (
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

// TODO

func (a *App) CreateEvent(id, title string) error {
	return nil
	//return a.storage.Create(storage.Event{ID: id, Title: title})
}
