package internalsender

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	iternalbroker "github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/broker"
)

type Sender struct {
	ctx     context.Context
	conf    Config
	adapter *Adapter
	logg    *Logger
}

func NewSender(ctx context.Context, config Config, adapter *Adapter, logger Logger) *Sender {
	return &Sender{ctx: ctx, conf: config, adapter: adapter, logg: &logger}
}

func (s *Sender) Start() error {
	if err := s.connectToBroker(); err != nil {
		return err
	}

	notifications, err := (*s.adapter).Consume(s.ctx, s.conf.Broker.Topic)
	if err != nil {
		return err
	}

	go func(notifications <-chan iternalbroker.Message) {
		f, err := os.Create(s.conf.File.Path)
		if err != nil {
			(*s.logg).Error(err.Error())
			return
		}
		defer f.Close()

	L:
		for {
			select {
			case <-s.ctx.Done():
				break L
			case msg, opened := <-notifications:
				if !opened {
					break L
				}

				f.WriteString(fmt.Sprintf("[NOTIFICATION SENT] %s \n", msg.Text))

				log.Printf("[NOTIFICATION SENT] %s\n", msg.Text)
			}
		}

		(*s.logg).Info("sender stopped")
	}(notifications)

	(*s.logg).Info("sender started")

	return nil
}

func (s *Sender) connectToBroker() error {
	delay, err := time.ParseDuration(s.conf.Broker.ConnectionTryDelay)
	if err != nil {
		return err
	}
	pause := time.NewTicker(delay)
	defer pause.Stop()

	currentTry := 1

	for {
		currentTry++
		err := (*s.adapter).InitConsumer()

		if err == nil {
			(*s.logg).Info("[Sender] Connected to broker")

			return nil
		}

		opError := new(net.OpError)
		if !errors.As(err, &opError) {
			return err
		}

		(*s.logg).Info("[Sender] Waiting for broker connection...")

		select {
		case <-s.ctx.Done():
			return nil
		case <-pause.C:
			if currentTry > s.conf.Broker.MaxConnectionTries {
				return errors.New("[Sender] Maximum tries reached")
			}
		}
	}
}

func (s *Sender) Stop() error {
	if err := (*s.adapter).CloseConsumer(); err != nil {
		return err
	}

	return nil
}
