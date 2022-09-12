package internalsender

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	iternalbroker "github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/broker"
)

type Sender struct {
	ctx     context.Context
	config  Config
	adapter *Adapter
	logger  Logger
}

func NewSender(ctx context.Context, config Config, adapter *Adapter, logger Logger) *Sender {
	return &Sender{ctx: ctx, config: config, adapter: adapter, logger: logger}
}

func (s *Sender) Start() error {
	if err := s.connectToBroker(); err != nil {
		return err
	}

	notifications, err := (*s.adapter).Consume(s.ctx, s.config.Broker.Topic)
	if err != nil {
		return err
	}

	go func(notifications <-chan iternalbroker.Message) {
	L:
		for {
			select {
			case <-s.ctx.Done():
				break L
			case msg, opened := <-notifications:
				if !opened {
					break L
				}

				time.Sleep(3 * time.Second)
				// send here
				fmt.Printf("[NOTIFICATION SENT] %s\n", msg.Text)
			}
		}

		s.logger.Info("sender stopped")
	}(notifications)

	s.logger.Info("sender started")

	return nil
}

func (s *Sender) connectToBroker() error {
	for t := 0; t < s.config.Broker.MaxConnectionTries; t++ {
		err := (*s.adapter).InitConsumer()

		if err == nil {
			s.logger.Info("[Sender] Connected to broker")

			return nil
		}

		opError := new(net.OpError)
		if errors.As(err, &opError) {
			s.logger.Info("[Sender] Waiting for broker connection...")
			delay, err := time.ParseDuration(s.config.Broker.ConnectionTryDelay)
			if err != nil {
				return err
			}
			time.Sleep(delay)

			continue
		} else {
			return err
		}
	}

	return errors.New("[Sender] Maximum tries reached")
}

func (s *Sender) Stop() error {
	if err := (*s.adapter).CloseConsumer(); err != nil {
		return err
	}

	return nil
}
