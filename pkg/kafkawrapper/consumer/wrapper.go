package kafkawrapper

import (
	"sync"

	"github.com/Beretta350/gochat/pkg/logger"
)

var consumerOnce sync.Once
var consumerWrapperInstance *consumerWrapper

type ConsumerWrapper interface {
	NewConsumer(groupId string) (ConsumerInterface, error)
	NewConsumerWithConfig(groupId string, config map[string]interface{}) (ConsumerInterface, error)
}

type consumerWrapper struct {
	brokers string
}

func InitConsumer(brokers string) {
	consumerOnce.Do(func() {
		consumerWrapperInstance = &consumerWrapper{
			brokers: brokers,
		}
	})
	logger.Info("Kafka consumer wrapper initialized")
}

// For testing purposes
func SetConsumerWrapper(brokers string) {
	consumerWrapperInstance = &consumerWrapper{
		brokers: brokers,
	}
}

// NewConsumer creates a new Kafka consumer with default configuration
func NewConsumer(groupId string) (ConsumerInterface, error) {
	adapter, err := NewConsumerAdapter(consumerWrapperInstance.brokers, groupId)
	if err != nil {
		return nil, err
	}
	return adapter, nil
}

// NewConsumerWithConfig creates a new Kafka consumer with custom configuration
func NewConsumerWithConfig(groupId string, config map[string]interface{}) (ConsumerInterface, error) {
	adapter, err := NewConsumerAdapterWithConfig(consumerWrapperInstance.brokers, groupId, config)
	if err != nil {
		return nil, err
	}
	return adapter, nil
}

// CloseConsumer is a placeholder as consumers are created on demand
func CloseConsumer() {
	// Actual closing happens when each consumer instance is closed
	logger.Info("Kafka consumer wrapper resources released")
}
