package internalscheduler

import "time"

type Message struct {
	Text, Topic string
}

type Notification struct {
	ID, Title string
	StartedAt time.Time
	UserID    string //?
}

type Adapter interface {
	Init() error
	Send(message Message) error
}

type Config struct {
	Logger   LoggerConf
	Database DatabaseConf
	Broker   BrokerConf
	DBWorker DBWorkerConf
}

type LoggerConf struct {
	Level string
}

type DatabaseConf struct {
	Dsn string
}

type DBWorkerConf struct {
	ScanPeriod      string
	ClearPeriodDays int
}

type BrokerAdapter string

const (
	KafkaAdapter BrokerAdapter = "kafka"
)

type BrokerConf struct {
	Adapter BrokerAdapter
	Address string
}
