package brokeradapters

import (
	internalscheduler "github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/scheduler"
	"github.com/Shopify/sarama"
)

type KafkaAdapter struct {
	producer *sarama.SyncProducer
	logg     internalscheduler.Logger
}

func NewKafkaAdapter(broker string, logg internalscheduler.Logger) (*KafkaAdapter, error) {
	p, err := createProducer(broker)
	if err != nil {
		return nil, err
	}

	return &KafkaAdapter{producer: p, logg: logg}, nil
}

func (a *KafkaAdapter) Send(message internalscheduler.Message) error {
	msg := &sarama.ProducerMessage{
		Topic: message.Topic,
		Value: sarama.StringEncoder(message.Text),
	}

	_, _, err := (*a.producer).SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}

func createProducer(broker string) (*sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true

	// TLS?

	producer, err := sarama.NewSyncProducer([]string{broker}, config)
	if err != nil {
		return nil, err
	}

	return &producer, nil
}
