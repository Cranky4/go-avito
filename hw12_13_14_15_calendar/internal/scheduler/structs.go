package internalscheduler

import (
	"time"

	iternalbroker "github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/broker"
)

type Notification struct {
	ID, Title string
	StartedAt time.Time
	UserID    string //?
}

type Adapter interface {
	InitProducer() error
	Produce(message iternalbroker.Message) error
}

type Config struct {
	Logger   LoggerConf
	Database DatabaseConf
	Broker   iternalbroker.BrokerConf
	DBWorker DBWorkerConf
}

type LoggerConf struct {
	Level string
}

type DatabaseConf struct {
	Dsn, ConnectionTryDelay string
	MaxConnectionTries      int
}

type DBWorkerConf struct {
	ScanPeriod      string
	ClearPeriodDays int
}
