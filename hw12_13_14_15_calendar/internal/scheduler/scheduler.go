package internalscheduler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net"
	"time"

	iternalbroker "github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/broker"
	// init pgsql.
	_ "github.com/jackc/pgx/stdlib"
)

type Scheduler struct {
	conf    Config
	ctx     context.Context
	db      *sql.DB
	adapter *Adapter
	logg    *Logger
}

func NewScheduler(ctx context.Context, config Config, adapter *Adapter, logg Logger) (*Scheduler, error) {
	s := &Scheduler{
		ctx:     ctx,
		adapter: adapter,
		logg:    &logg,
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

	err = s.db.Ping()
	if err != nil {
		return err
	}

	return nil
}

func (s *Scheduler) Start() error {
	err := s.connectToDatabase()
	if err != nil {
		return err
	}

	err = s.connectToBroker()
	if err != nil {
		return err
	}

	notifications := make(chan Notification)
	go func(ctx context.Context, db *sql.DB, ch chan Notification, logg Logger) {
		startDBWorker(ctx, db, ch, logg, s.conf.DBWorker)

		<-ctx.Done()
	}(s.ctx, s.db, notifications, *s.logg)

O:
	for {
		select {
		case <-s.ctx.Done():
			break O
		case notification := <-notifications:
			n, err := json.Marshal(notification)
			if err != nil {
				(*s.logg).Error(err.Error())
			} else {
				(*s.adapter).Produce(iternalbroker.Message{
					Topic: s.conf.Broker.Topic,
					Text:  string(n),
				})
			}
		}
	}

	return nil
}

func (s *Scheduler) connectToBroker() error {
	for t := 0; t < s.conf.Broker.MaxConnectionTries; t++ {
		err := (*s.adapter).InitProducer()

		if err == nil {
			(*s.logg).Info("[Scheduler] Connected to broker")

			return nil
		}

		opError := new(net.OpError)
		if errors.As(err, &opError) {
			(*s.logg).Info("[Scheduler] Waiting for broker connection...")
			delay, err := time.ParseDuration(s.conf.Broker.ConnectionTryDelay)
			if err != nil {
				return err
			}
			time.Sleep(delay)

			continue
		} else {
			return err
		}
	}

	return errors.New("[Scheduler] Maximum tries reached")
}

func (s *Scheduler) connectToDatabase() error {
	for t := 0; t < s.conf.Database.MaxConnectionTries; t++ {
		err := s.ensureDBConnected()

		if err == nil {
			(*s.logg).Info("[Scheduler] Connected to database")

			return nil
		}

		opError := new(net.OpError)
		if errors.As(err, &opError) {
			(*s.logg).Info("[Scheduler] Waiting for database connection...")
			delay, err := time.ParseDuration(s.conf.Database.ConnectionTryDelay)
			if err != nil {
				return err
			}
			time.Sleep(delay)

			continue
		} else {
			return err
		}
	}

	return errors.New("[Scheduler] Maximum tries reached")
}

func (s *Scheduler) Stop() error {
	(*s.db).Close()

	return nil
}

func startDBWorker(ctx context.Context, db *sql.DB, ch chan Notification, logg Logger, conf DBWorkerConf) {
	worker := NewDBWorker(db, conf, logg)

	worker.Run(ctx, ch)
}
