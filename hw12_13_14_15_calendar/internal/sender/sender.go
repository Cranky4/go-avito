package internalsender

import (
	"context"
	"fmt"

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
	if err := (*s.adapter).InitConsumer(); err != nil {
		return err
	}

	notifications, err := (*s.adapter).Consume(s.ctx, s.config.Broker.Topic)
	if err != nil {
		return err
	}

	s.logger.Info("sender started")

	go func(notifications <-chan iternalbroker.Message) {
	L:
		for {
			select {
			case <-s.ctx.Done():
				break L
			case msg := <-notifications:
				// send here
				fmt.Printf("[NOTIFICATION SENT] %s\n", msg.Text)
			}
		}
	}(notifications)

	return nil
}

func (s *Sender) Stop() error {
	if err := (*s.adapter).CloseConsumer(); err != nil {
		return err
	}

	return nil
}
