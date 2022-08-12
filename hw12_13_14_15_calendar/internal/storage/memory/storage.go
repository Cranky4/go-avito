package memorystorage

import (
	"sync"
	"time"

	"github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	_      sync.RWMutex // TODO: накрутить локов
	events map[storage.EventID]storage.Event
}

func New() *Storage {
	return &Storage{
		events: make(map[storage.EventID]storage.Event),
	}
}

func (s *Storage) CreateEvent(event storage.Event) error {
	isBusy, err := s.IsPeriodBusy(event.StartsAt, event.EndsAt, nil)
	if err != nil {
		return err
	}
	if isBusy {
		return storage.ErrDateBusy
	}
	s.events[event.ID] = event

	return nil
}

func (s *Storage) UpdateEvent(id storage.EventID, event storage.Event) error {
	event, exists := s.events[id]

	if !exists {
		return storage.ErrEventNotFound
	}

	isBusy, err := s.IsPeriodBusy(event.StartsAt, event.EndsAt, []string{id.String()})
	if err != nil {
		return err
	}
	if isBusy {
		return storage.ErrDateBusy
	}

	event.ID = id
	s.events[id] = event

	return nil
}

func (s *Storage) DeleteEvent(id storage.EventID) error {
	event, exists := s.events[id]

	if !exists {
		return storage.ErrEventNotFound
	}
	delete(s.events, event.ID)

	return nil
}

func (s *Storage) GetEvents(dateFrom, dateTo time.Time) []storage.Event {
	var list []storage.Event

	for _, event := range s.events {
		if event.StartsAt.After(dateFrom) && event.StartsAt.Before(dateTo) {
			list = append(list, event)
		}
	}

	return list
}

func (s *Storage) GetEvent(id storage.EventID) (storage.Event, error) {
	for eventID, event := range s.events {
		if eventID == id {
			return event, nil
		}
	}

	return storage.Event{}, storage.ErrEventNotFound
}

func (s *Storage) GetDayEvents(date time.Time) ([]storage.Event, error) {
	year, month, day := date.Date()
	fromDate := time.Date(year, month, day, 0, 0, 0, 0, date.Location())
	toDate := time.Date(year, month, day, 23, 59, 59, 0, date.Location())

	return s.GetEvents(fromDate, toDate), nil
}

func (s *Storage) GetWeekEvents(date time.Time) ([]storage.Event, error) {
	year, month, day := date.Date()

	fromDate := time.Date(year, month, day, 0, 0, 0, 0, date.Location())
	toDate := fromDate.AddDate(0, 0, 7)

	return s.GetEvents(fromDate, toDate), nil
}

func (s *Storage) GetMonthEvents(date time.Time) ([]storage.Event, error) {
	year, month, day := date.Date()

	fromDate := time.Date(year, month, day, 0, 0, 0, 0, date.Location())
	toDate := fromDate.AddDate(0, 1, 0)

	return s.GetEvents(fromDate, toDate), nil
}

func (s *Storage) IsPeriodBusy(dateFrom, dateTo time.Time, excludeIds []string) (bool, error) {
	excludeIdsMap := make(map[string]int)

	for _, id := range excludeIds {
		excludeIdsMap[id] = 1
	}

	for _, event := range s.events {
		_, excluded := excludeIdsMap[event.ID.String()]

		if !(event.EndsAt.Before(dateFrom) || event.StartsAt.After(dateTo)) && !excluded {
			return true, nil
		}
	}
	return false, nil
}
