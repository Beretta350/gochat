package kafkafactory

import (
	"sync"

	"github.com/Beretta350/gochat/pkg/logger"
)

var consumerOnce sync.Once
var consumerFactoryInstance *consumerFactory

type ConsumerFactory interface {
	NewConsumer(groupId string) (ConsumerInterface, error)
	NewConsumerWithConfig(groupId string, config map[string]interface{}) (ConsumerInterface, error)
}

type consumerFactory struct {
	brokers string
}

func Init(brokers string) {
	consumerOnce.Do(func() {
		consumerFactoryInstance = &consumerFactory{
			brokers: brokers,
		}
	})
	logger.Info("Kafka consumer factory initialized")
}

// For testing purposes
func SetConsumerFactory(brokers string) {
	consumerFactoryInstance = &consumerFactory{
		brokers: brokers,
	}
}

// NewConsumer creates a new Kafka consumer with default configuration
func NewConsumer(groupId string) (ConsumerInterface, error) {
	adapter, err := NewConsumerAdapter(consumerFactoryInstance.brokers, groupId)
	if err != nil {
		return nil, err
	}
	return adapter, nil
}

// NewConsumerWithConfig creates a new Kafka consumer with custom configuration
func NewConsumerWithConfig(groupId string, config map[string]interface{}) (ConsumerInterface, error) {
	adapter, err := NewConsumerAdapterWithConfig(consumerFactoryInstance.brokers, groupId, config)
	if err != nil {
		return nil, err
	}
	return adapter, nil
}

// CloseConsumer is a placeholder as consumers are created on demand
func CloseConsumer() {
	// Actual closing happens when each consumer instance is closed
	logger.Info("Kafka consumer factory resources released")
}
