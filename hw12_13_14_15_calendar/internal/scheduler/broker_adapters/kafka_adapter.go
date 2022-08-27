package brokeradapters

import (
	internalscheduler "github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/scheduler"
	"github.com/Shopify/sarama"
)

type KafkaAdapter struct {
	producer   *sarama.SyncProducer
	brokerAddr string
	logg       internalscheduler.Logger
}

func NewKafkaAdapter(brokerAddr string, logg internalscheduler.Logger) (*KafkaAdapter, error) {
	p, err := createProducer(brokerAddr)
	if err != nil {
		return nil, err
	}

	return &KafkaAdapter{producer: p, brokerAddr: brokerAddr, logg: logg}, nil
}

func (a *KafkaAdapter) Init() error {
	broker := sarama.NewBroker(a.brokerAddr)

	if err := broker.Open(nil); err != nil {
		return err
	}

	response, err := broker.GetMetadata(&sarama.MetadataRequest{Topics: []string{"notifications"}})
	if err != nil {
		return err
	}

	exists := false
	for _, t := range response.Topics {
		if t.Name == "notifications" {
			exists = true
			break
		}
	}

	if !exists {
		if err := createTopics(broker); err != nil {
			return err
		}
	}

	if err := broker.Close(); err != nil {
		return err
	}

	return nil
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

func createTopics(broker *sarama.Broker) error {
	topics := make(map[string]*sarama.TopicDetail)
	topics["notifications"] = &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	_, err := broker.CreateTopics(&sarama.CreateTopicsRequest{TopicDetails: topics})
	if err != nil {
		return err
	}

	return nil
}
