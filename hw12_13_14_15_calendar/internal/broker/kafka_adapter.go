package iternalbroker

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
)

type KafkaAdapter struct {
	producer      *sarama.SyncProducer
	consumerGroup *sarama.ConsumerGroup
	config        BrokerConf
	logg          Logger
}

func NewKafkaAdapter(conf BrokerConf, logg Logger) (*KafkaAdapter, error) {
	return &KafkaAdapter{config: conf, logg: logg}, nil
}

func (a *KafkaAdapter) InitProducer() error {
	broker := sarama.NewBroker(a.config.Address)

	if err := broker.Open(nil); err != nil {
		return err
	}

	response, err := broker.GetMetadata(&sarama.MetadataRequest{Topics: []string{a.config.Topic}})
	if err != nil {
		return err
	}

	exists := false
	for _, t := range response.Topics {
		if t.Name == a.config.Topic {
			exists = true
			break
		}
	}

	if !exists {
		if err := createTopics(broker, []string{a.config.Topic}); err != nil {
			return err
		}
	}

	if err := broker.Close(); err != nil {
		return err
	}

	p, err := createProducer(a.config)
	if err != nil {
		return err
	}
	a.producer = p

	return nil
}

func (a *KafkaAdapter) Produce(message Message) error {
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

func (a *KafkaAdapter) InitConsumer() error {
	ver, err := sarama.ParseKafkaVersion(a.config.Version)
	if err != nil {
		return err
	}

	config := sarama.NewConfig()
	config.Version = ver
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumerGroup([]string{a.config.Address}, a.config.Topic, config)
	if err != nil {
		return err
	}
	a.consumerGroup = &consumer

	return nil
}

type ConsumerHandler struct {
	out    chan Message
	logger Logger
}

func (c *ConsumerHandler) Setup(s sarama.ConsumerGroupSession) error {
	return nil
}

func (c *ConsumerHandler) Cleanup(s sarama.ConsumerGroupSession) error {
	c.logger.Info(fmt.Sprintf("handler cleanup %#v", s))
	return nil
}

func (c *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case <-session.Context().Done():
			return nil
		case message := <-claim.Messages():
			c.logger.Info(
				fmt.Sprintf(
					"Message claimed: value = %s, timestamp = %v, topic = %s",
					string(message.Value),
					message.Timestamp,
					message.Topic,
				),
			)
			session.MarkMessage(message, "")
			c.out <- Message{Text: string(message.Value)}
		}
	}
}

func (a *KafkaAdapter) Consume(ctx context.Context, topic string) (<-chan Message, error) {
	// errors
	go func() {
		a.logg.Info("error log started")
		for {
			select {
			case <-ctx.Done():
				a.logg.Info("error log stopped")
				return
			case err := <-(*a.consumerGroup).Errors():
				a.logg.Error(err.Error())
			}
		}
	}()

	// setup
	out := make(chan Message)
	consumerHandler := ConsumerHandler{out: out, logger: a.logg}

	go func(ctx context.Context, topic string) {
	L:
		for {
			select {
			case <-ctx.Done():
				break L
			default:
				if err := (*a.consumerGroup).Consume(ctx, []string{topic}, &consumerHandler); err != nil {
					a.logg.Error(err.Error())
					break L
				}
				if ctx.Err() != nil {
					a.logg.Error(ctx.Err().Error())
					break L
				}
			}
		}

		close(consumerHandler.out)
		a.logg.Info("Consumer stopped")
	}(ctx, topic)

	a.logg.Info("Consumer up and running!")

	return out, nil
}

func (a *KafkaAdapter) CloseConsumer() error {
	if err := (*a.consumerGroup).Close(); err != nil {
		return err
	}
	a.logg.Info("Consumer closed...")

	return nil
}

func createProducer(conf BrokerConf) (*sarama.SyncProducer, error) {
	ver, err := sarama.ParseKafkaVersion(conf.Version)
	if err != nil {
		return nil, err
	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true
	config.Version = ver

	// TLS?

	producer, err := sarama.NewSyncProducer([]string{conf.Address}, config)
	if err != nil {
		return nil, err
	}

	return &producer, nil
}

func createTopics(broker *sarama.Broker, topicsToCreate []string) error {
	topics := make(map[string]*sarama.TopicDetail)
	for _, t := range topicsToCreate {
		topics[t] = &sarama.TopicDetail{
			NumPartitions:     1,
			ReplicationFactor: 1,
		}
	}

	_, err := broker.CreateTopics(&sarama.CreateTopicsRequest{TopicDetails: topics})
	if err != nil {
		return err
	}

	return nil
}
