package kafkafactory

import (
	"sync"

	"github.com/Beretta350/gochat/pkg/logger"
)

var producerOnce sync.Once
var producerFactoryInstance *producerFactory

type ProducerFactory interface {
	NewProducer() (ProducerInterface, error)
	NewProducerWithConfig(config map[string]interface{}) (ProducerInterface, error)
	Close()
}

type producerFactory struct {
	brokers string
}

func Init(brokers string) {
	producerOnce.Do(func() {
		producerFactoryInstance = &producerFactory{
			brokers: brokers,
		}
	})
	logger.Info("Kafka producer factory initialized")
}

// For testing purposes
func SetProducerFactory(brokers string) {
	producerFactoryInstance = &producerFactory{
		brokers: brokers,
	}
}

// NewProducer creates a new Kafka producer with default configuration
func NewProducer() (ProducerInterface, error) {
	adapter, err := NewProducerAdapter(producerFactoryInstance.brokers)
	if err != nil {
		return nil, err
	}
	return adapter, nil
}

// NewProducerWithConfig creates a new Kafka producer with custom configuration
func NewProducerWithConfig(config map[string]interface{}) (ProducerInterface, error) {
	adapter, err := NewProducerAdapterWithConfig(producerFactoryInstance.brokers, config)
	if err != nil {
		return nil, err
	}
	return adapter, nil
}

// Close closes the producer
func CloseProducer() {
	// This is a placeholder as producers are created on demand
	// Actual closing happens when each producer instance is closed
	logger.Info("Kafka producer factory resources released")
}
