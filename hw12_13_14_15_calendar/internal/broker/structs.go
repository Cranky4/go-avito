package iternalbroker

type Message struct {
	Text, Topic string
}

type BrokerAdapter string

const (
	KafkaBrokerAdapter BrokerAdapter = "kafka"
)

type BrokerConf struct {
	Adapter                 BrokerAdapter
	Address, Topic, Version string
}
