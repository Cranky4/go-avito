package internalscheduler

import (
	"context"
	"database/sql"
	"encoding/json"

	iternalbroker "github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/broker"
	// init pgsql.
	_ "github.com/jackc/pgx/stdlib"
)

type Scheduler struct {
	conf    Config
	ctx     context.Context
	db      *sql.DB
	adapter *Adapter
	logger  *Logger
}

func NewScheduler(ctx context.Context, config Config, adapter *Adapter, logger Logger) (*Scheduler, error) {
	s := &Scheduler{
		ctx:     ctx,
		adapter: adapter,
		logger:  &logger,
		conf:    config,
	}

	return s, nil
}

func (s *Scheduler) ensureDBConnected() error {
	db, err := sql.Open("pgx", s.conf.Database.Dsn)
	if err != nil {
		return err
	}
	s.db = db

	return nil
}

func (s *Scheduler) Start() error {
	if err := (*s.adapter).InitProducer(); err != nil {
		return err
	}
	if err := s.ensureDBConnected(); err != nil {
		return err
	}

	notifications := make(chan Notification)
	go func(ctx context.Context, db *sql.DB, ch chan Notification, logg Logger) {
		startDBWorker(ctx, db, ch, logg, s.conf.DBWorker)

		<-ctx.Done()
	}(s.ctx, s.db, notifications, *s.logger)

O:
	for {
		select {
		case <-s.ctx.Done():
			break O
		case notification := <-notifications:
			n, err := json.Marshal(notification)
			if err != nil {
				(*s.logger).Error(err.Error())
			} else {
				(*s.adapter).Produce(iternalbroker.Message{
					Topic: "notifications",
					Text:  string(n),
				})
			}
		}
	}

	return nil
}

func (s *Scheduler) Stop() error {
	(*s.db).Close()

	return nil
}

func startDBWorker(ctx context.Context, db *sql.DB, ch chan Notification, logg Logger, conf DBWorkerConf) {
	worker := NewDBWorker(db, conf, logg)

	worker.Run(ctx, ch)
}
