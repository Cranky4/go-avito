package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/storage"
	// init pgsql.
	_ "github.com/jackc/pgx/stdlib"
)

var ErrNotConnected = errors.New("not connected to database")

type Storage struct {
	dsn     string
	db      *sql.DB
	context context.Context
}

func New(ctx context.Context, dsn string) *Storage {
	return &Storage{dsn: dsn, context: ctx}
}

func (s *Storage) ensureConnected() error {
	if s.db == nil {
		if err := s.Connect(s.context); err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sql.Open("pgx", s.dsn)
	if err != nil {
		return err
	}
	s.db = db

	if err := s.db.Ping(); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	if s.db == nil {
		return ErrNotConnected
	}
	s.db.Close()

	return nil
}

func (s *Storage) CreateEvent(event storage.Event) error {
	err := s.ensureConnected()
	if err != nil {
		return err
	}

	isBusy, err := s.IsPeriodBusy(event.StartsAt, event.EndsAt, nil)
	if err != nil {
		return err
	}
	if isBusy {
		return storage.ErrDateBusy
	}

	var notifyAfter sql.NullTime
	if event.NotifyAfter.IsSet {
		notifyAfter.Time = event.NotifyAfter.Time
		notifyAfter.Valid = true
	}

	err = s.execTransactionally(
		"INSERT INTO events(id, title, starts_at, ends_at, notify_after) VALUES($1, $2, $3, $4, $5)",
		event.ID.String(),
		event.Title,
		event.StartsAt,
		event.EndsAt,
		notifyAfter,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) IsPeriodBusy(dateFrom, dateTo time.Time, excludeIds []string) (bool, error) {
	if err := s.ensureConnected(); err != nil {
		return false, err
	}

	// стартовый запрос
	var builder strings.Builder
	builder.WriteString("SELECT COUNT(id) FROM events WHERE (starts_at <= $2 AND ends_at >= $1)")

	params := make([]interface{}, 0, len(excludeIds)+2)
	params = append(params, dateFrom)
	params = append(params, dateTo)

	// сборка IN-услвовия на исключение идешников
	if len(excludeIds) > 0 {
		builder.WriteString("AND id NOT IN (")

		for i, id := range excludeIds {
			builder.WriteString("$")
			builder.WriteString(strconv.Itoa(i + 3))

			if i < len(excludeIds)-1 {
				builder.WriteString(", ")
			}

			params = append(params, id)
		}
		builder.WriteString(")")
	}

	stmt, err := s.db.Prepare(builder.String())
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(s.context, params...)
	if err != nil {
		return false, err
	}
	if rows.Err() != nil {
		return false, rows.Err()
	}

	var count int
	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return false, err
		}
		defer rows.Close()
	}

	return count > 0, nil
}

func (s *Storage) GetEvent(id storage.EventID) (storage.Event, error) {
	err := s.ensureConnected()
	if err != nil {
		return storage.Event{}, err
	}

	stmt, err := s.db.Prepare("SELECT title, starts_at, ends_at, notify_after FROM events WHERE id=$1")
	if err != nil {
		return storage.Event{}, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(s.context, id.String())
	if err != nil {
		return storage.Event{}, err
	}
	if rows.Err() != nil {
		return storage.Event{}, rows.Err()
	}

	var title string
	var startsAt, endsAt time.Time
	var notifyAfterHandler sql.NullTime

	found := false

	if rows.Next() {
		err = rows.Scan(&title, &startsAt, &endsAt, &notifyAfterHandler)
		if err != nil {
			return storage.Event{}, err
		}
		defer rows.Close()
		found = true
	}

	if !found {
		return storage.Event{}, storage.ErrEventNotFound
	}

	var notifyAfter storage.NotifyAfter
	if notifyAfterHandler.Valid {
		notifyAfter.Time = notifyAfterHandler.Time
		notifyAfter.IsSet = true
	}

	return storage.Event{
		ID:          id,
		Title:       title,
		StartsAt:    startsAt,
		EndsAt:      endsAt,
		NotifyAfter: notifyAfter,
	}, nil
}

func (s *Storage) UpdateEvent(id storage.EventID, event storage.Event) error {
	err := s.ensureConnected()
	if err != nil {
		return err
	}

	isBusy, err := s.IsPeriodBusy(
		event.StartsAt,
		event.EndsAt,
		[]string{id.String()},
	)
	if err != nil {
		return err
	}
	if isBusy {
		return storage.ErrDateBusy
	}

	var notifyAfter sql.NullTime
	if event.NotifyAfter.IsSet {
		notifyAfter.Time = event.NotifyAfter.Time
		notifyAfter.Valid = true
	}

	err = s.execTransactionally(
		"UPDATE events SET title=$2, starts_at=$3, ends_at=$4, notify_after=$5 WHERE id=$1",
		event.ID.String(),
		event.Title,
		event.StartsAt,
		event.EndsAt,
		notifyAfter,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetEvents(dateFrom, dateTo time.Time) ([]storage.Event, error) {
	err := s.ensureConnected()
	if err != nil {
		return []storage.Event{}, err
	}

	stmt, err := s.db.Prepare(
		"SELECT id, title, starts_at, ends_at, notify_after  FROM events WHERE starts_at >= $1 AND starts_at < $2",
	)
	if err != nil {
		return []storage.Event{}, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(s.context, dateFrom, dateTo)
	if err != nil {
		return []storage.Event{}, err
	}
	if rows.Err() != nil {
		return []storage.Event{}, rows.Err()
	}

	events := make([]storage.Event, 0, 10) // TODO cap

	var id string
	var title string
	var startsAt, endsAt time.Time
	var notifyAfterHandler sql.NullTime

	for rows.Next() {
		err = rows.Scan(&id, &title, &startsAt, &endsAt, &notifyAfterHandler)
		if err != nil {
			return []storage.Event{}, err
		}
		defer rows.Close()

		eventID, err := storage.NewEventIDFromString(id)
		if err != nil {
			return []storage.Event{}, err
		}

		var notifyAfter storage.NotifyAfter
		if notifyAfterHandler.Valid {
			notifyAfter.IsSet = true
			notifyAfter.Time = notifyAfterHandler.Time
		}

		events = append(events, storage.Event{
			ID:          eventID,
			Title:       title,
			StartsAt:    startsAt,
			EndsAt:      endsAt,
			NotifyAfter: notifyAfter,
		})
	}

	return events, nil
}

func (s *Storage) GetDayEvents(date time.Time) ([]storage.Event, error) {
	year, month, day := date.Date()
	fromDate := time.Date(year, month, day, 0, 0, 0, 0, date.Location())
	toDate := time.Date(year, month, day, 23, 59, 59, 0, date.Location())

	return s.GetEvents(fromDate, toDate)
}

func (s *Storage) GetWeekEvents(date time.Time) ([]storage.Event, error) {
	year, month, day := date.Date()
	fromDate := time.Date(year, month, day, 0, 0, 0, 0, date.Location())
	toDate := fromDate.AddDate(0, 0, 7)

	return s.GetEvents(fromDate, toDate)
}

func (s *Storage) GetMonthEvents(date time.Time) ([]storage.Event, error) {
	year, month, day := date.Date()
	fromDate := time.Date(year, month, day, 0, 0, 0, 0, date.Location())
	toDate := fromDate.AddDate(0, 1, 0)

	return s.GetEvents(fromDate, toDate)
}

func (s *Storage) DeleteEvent(id storage.EventID) error {
	err := s.ensureConnected()
	if err != nil {
		return err
	}

	if _, err = s.GetEvent(id); err != nil {
		return err
	}

	err = s.execTransactionally("DELETE FROM events WHERE id=$1", id.String())

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) execTransactionally(query string, params ...interface{}) error {
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	tx, err := s.db.BeginTx(s.context, nil)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(s.context, params...)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
