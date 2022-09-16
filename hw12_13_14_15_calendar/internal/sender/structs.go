package internalsender

import (
	"context"

	iternalbroker "github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/broker"
)

type Adapter interface {
	InitConsumer() error
	Consume(ctx context.Context, topic string) (<-chan iternalbroker.Message, error)
	CloseConsumer() error
}

type Config struct {
	Logger LoggerConf
	Broker iternalbroker.BrokerConf
	File   FileConf
}

type LoggerConf struct {
	Level string
}

type FileConf struct {
	Path string
}
