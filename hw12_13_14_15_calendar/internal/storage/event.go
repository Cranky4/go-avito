package storage

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrDateBusy      = errors.New("date is busy")
)

type EventID struct {
	uuid.UUID
}

type NotifyAfter struct {
	Time  time.Time
	IsSet bool
}

func NewEventID() EventID {
	return EventID{uuid.New()}
}

func NewEventIDFromString(id string) (EventID, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return EventID{}, err
	}

	return EventID{uuid}, nil
}

type Event struct {
	ID               EventID
	Title            string
	StartsAt, EndsAt time.Time // дата и время начала и завершения события
	NotifyAfter      NotifyAfter
}

type EventStorage interface {
	AddEvent(event Event) error
	UpdateEvent(id EventID, event Event) error
	DeleteEvent(id EventID) error
	GetDayEvents(date time.Time) ([]Event, error)
	GetWeekEvents(date time.Time) ([]Event, error)
	GetMonthEvents(date time.Time) ([]Event, error)
	IsPeriodBusy(dateFrom, dateTo time.Time) (bool, error)
}
