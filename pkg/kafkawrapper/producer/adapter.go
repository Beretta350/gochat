package kafkawrapper

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// ProducerInterface is an interface to make Kafka Producer mockable
type ProducerInterface interface {
	Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error
	Flush(timeoutMs int) int
	Close()
	Events() chan kafka.Event
}

// ProducerAdapter adapts the Kafka Producer to implement ProducerInterface
type ProducerAdapter struct {
	producer *kafka.Producer
}

// NewProducerAdapter creates a new ProducerAdapter
func NewProducerAdapter(brokers string) (*ProducerAdapter, error) {
	return NewProducerAdapterWithConfig(brokers, nil)
}

// NewProducerAdapterWithConfig creates a new ProducerAdapter with custom configuration
func NewProducerAdapterWithConfig(brokers string, customConfig map[string]interface{}) (*ProducerAdapter, error) {
	// Create base configuration
	config := &kafka.ConfigMap{
		"bootstrap.servers": brokers,
	}

	// Apply custom configuration if provided
	for key, value := range customConfig {
		if err := config.SetKey(key, value); err != nil {
			return nil, err
		}
	}

	producer, err := kafka.NewProducer(config)
	if err != nil {
		return nil, err
	}
	return &ProducerAdapter{producer: producer}, nil
}

// Produce implements ProducerInterface
func (p *ProducerAdapter) Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error {
	return p.producer.Produce(msg, deliveryChan)
}

// Flush implements ProducerInterface
func (p *ProducerAdapter) Flush(timeoutMs int) int {
	return p.producer.Flush(timeoutMs)
}

// Close implements ProducerInterface
func (p *ProducerAdapter) Close() {
	p.producer.Close()
}

// Events implements ProducerInterface
func (p *ProducerAdapter) Events() chan kafka.Event {
	return p.producer.Events()
}
