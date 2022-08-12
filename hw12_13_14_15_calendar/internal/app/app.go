package app

import (
	"context"
	"fmt"
	"time"

	"github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Storage interface {
	CreateEvent(storage.Event) error
	UpdateEvent(storage.EventID, storage.Event) error
	GetEvent(storage.EventID) (storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{logger: logger, storage: storage}
}

func (a *App) CreateEvent(ctx context.Context, id storage.EventID, title string, startsAt, endsAt time.Time) error {
	event := storage.Event{
		ID:       id,
		Title:    title,
		StartsAt: startsAt,
		EndsAt:   endsAt,
	}

	err := a.storage.CreateEvent(event)
	if err != nil {
		a.logger.Error(err.Error())

		return err
	}

	a.logger.Info(
		fmt.Sprintf(
			"%s from %s to %s successfully created",
			event.Title,
			event.StartsAt.Format(time.UnixDate),
			event.EndsAt.Format(time.UnixDate),
		),
	)
	return nil
}

func (a *App) UpdateEvent(ctx context.Context, id storage.EventID, title string, startsAt, endsAt time.Time) error {
	event := storage.Event{
		ID:       id,
		Title:    title,
		StartsAt: startsAt,
		EndsAt:   endsAt,
	}

	if err := a.storage.UpdateEvent(id, event); err != nil {
		a.logger.Error(err.Error())

		return err
	}

	a.logger.Info(
		fmt.Sprintf(
			"%s from %s to %s successfully created",
			event.Title,
			event.StartsAt.Format(time.UnixDate),
			event.EndsAt.Format(time.UnixDate),
		),
	)
	return nil
}

func (a *App) GetEvent(ctx context.Context, id storage.EventID) (storage.Event, error) {
	event, err := a.storage.GetEvent(id)
	if err != nil {
		a.logger.Error(err.Error())

		return storage.Event{}, err
	}

	return event, nil
}
