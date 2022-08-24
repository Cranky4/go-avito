package internalscheduler

import (
	"context"
	"database/sql"
	"fmt"

	// init pgsql.
	_ "github.com/jackc/pgx/stdlib"
)

type Scheduler struct {
	dsn      string
	db       *sql.DB
	producer *Producer
	logger   *Logger
}

func NewScheduler(ctx context.Context, dsn string, producer *Producer, logger Logger) (*Scheduler, error) {
	s := &Scheduler{
		producer: producer,
		logger:   &logger,
	}

	return s, nil
}

func (s *Scheduler) Start() error {
	// connect to db
	db, err := sql.Open("pgx", s.dsn)
	if err != nil {
		return err
	}
	s.db = db

	// connect to message broker
	msg := Message{
		Topic: "test",
		Text:  "Hello",
	}
	err = (*s.producer).Send(msg)
	(*s.logger).Info(fmt.Sprintf("%#v %#v", msg, err))
	if err != nil {
		return err
	}

	return nil
}

func (s *Scheduler) Stop(ctx context.Context) error {
	(*s.db).Close()

	return nil
}
