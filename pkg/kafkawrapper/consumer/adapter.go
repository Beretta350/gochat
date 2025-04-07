package kafkawrapper

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// ConsumerInterface is an interface to make Kafka Consumer mockable
type ConsumerInterface interface {
	Subscribe(topic string, rebalanceCb kafka.RebalanceCb) error
	Unsubscribe() error
	ReadMessage(timeout time.Duration) (*kafka.Message, error)
	Close() error
	CommitMessage(m *kafka.Message) ([]kafka.TopicPartition, error)
	Commit() ([]kafka.TopicPartition, error)
	SubscribeTopics(topics []string, rebalanceCb kafka.RebalanceCb) error
}

// ConsumerAdapter adapts the Kafka Consumer to implement ConsumerInterface
type ConsumerAdapter struct {
	consumer *kafka.Consumer
}

// NewConsumerAdapter creates a new ConsumerAdapter
func NewConsumerAdapter(brokers string, groupId string) (*ConsumerAdapter, error) {
	return NewConsumerAdapterWithConfig(brokers, groupId, nil)
}

// NewConsumerAdapterWithConfig creates a new ConsumerAdapter with custom configuration
func NewConsumerAdapterWithConfig(brokers string, groupId string, customConfig map[string]interface{}) (*ConsumerAdapter, error) {
	// Create base configuration
	config := &kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          groupId,
		"auto.offset.reset": "earliest",
	}

	// Apply custom configuration if provided
	for key, value := range customConfig {
		if err := config.SetKey(key, value); err != nil {
			return nil, err
		}
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, err
	}
	return &ConsumerAdapter{consumer: consumer}, nil
}

// Subscribe implements ConsumerInterface
func (c *ConsumerAdapter) Subscribe(topic string, rebalanceCb kafka.RebalanceCb) error {
	return c.consumer.Subscribe(topic, rebalanceCb)
}

// Unsubscribe implements ConsumerInterface
func (c *ConsumerAdapter) Unsubscribe() error {
	return c.consumer.Unsubscribe()
}

// ReadMessage implements ConsumerInterface
func (c *ConsumerAdapter) ReadMessage(timeout time.Duration) (*kafka.Message, error) {
	return c.consumer.ReadMessage(timeout)
}

// Close implements ConsumerInterface
func (c *ConsumerAdapter) Close() error {
	return c.consumer.Close()
}

// CommitMessage implements ConsumerInterface
func (c *ConsumerAdapter) CommitMessage(m *kafka.Message) ([]kafka.TopicPartition, error) {
	return c.consumer.CommitMessage(m)
}

// Commit implements ConsumerInterface
func (c *ConsumerAdapter) Commit() ([]kafka.TopicPartition, error) {
	return c.consumer.Commit()
}

// SubscribeTopics implements ConsumerInterface
func (c *ConsumerAdapter) SubscribeTopics(topics []string, rebalanceCb kafka.RebalanceCb) error {
	return c.consumer.SubscribeTopics(topics, rebalanceCb)
}
