package kafkawrapper

import (
	"sync"

	"github.com/Beretta350/gochat/pkg/logger"
)

var producerOnce sync.Once
var producerWrapperInstance *producerWrapper

type ProducerWrapper interface {
	NewProducer() (ProducerInterface, error)
	NewProducerWithConfig(config map[string]interface{}) (ProducerInterface, error)
	Close()
}

type producerWrapper struct {
	brokers string
}

func InitProducer(brokers string) {
	producerOnce.Do(func() {
		producerWrapperInstance = &producerWrapper{
			brokers: brokers,
		}
	})
	logger.Info("Kafka producer wrapper initialized")
}

// For testing purposes
func SetProducerWrapper(brokers string) {
	producerWrapperInstance = &producerWrapper{
		brokers: brokers,
	}
}

// NewProducer creates a new Kafka producer with default configuration
func NewProducer() (ProducerInterface, error) {
	adapter, err := NewProducerAdapter(producerWrapperInstance.brokers)
	if err != nil {
		return nil, err
	}
	return adapter, nil
}

// NewProducerWithConfig creates a new Kafka producer with custom configuration
func NewProducerWithConfig(config map[string]interface{}) (ProducerInterface, error) {
	adapter, err := NewProducerAdapterWithConfig(producerWrapperInstance.brokers, config)
	if err != nil {
		return nil, err
	}
	return adapter, nil
}

// Close closes the producer
func CloseProducer() {
	// This is a placeholder as producers are created on demand
	// Actual closing happens when each producer instance is closed
	logger.Info("Kafka producer wrapper resources released")
}
