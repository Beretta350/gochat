package kafka_wrapper

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/Beretta350/gochat/pkg/logger"
)

var once sync.Once
var kafkaWrapperInstance *kafkaWrapper

type KafkaWrapper interface {
	CreateTopic(ctx context.Context, name string) error
	DeleteTopic(ctx context.Context, name string) error
	NewConsumer(groupId string) (*kafka.Consumer, error)
	NewProducer() (*kafka.Producer, error)
	Close()
}

type kafkaWrapper struct {
	adminClient *kafka.AdminClient
	brokers     string
}

func Init(brokers string) {
	once.Do(func() {
		adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": brokers})
		if err != nil {
			logger.Fatal(err)
		}

		kafkaWrapperInstance = &kafkaWrapper{
			adminClient: adminClient,
			brokers:     brokers,
		}
	})
	logger.Info("Kafka wrapper and admin client initialized")
}

func CreateTopic(ctx context.Context, name string) error {
	// Check if the topic already exists
	describeTopics, err := kafkaWrapperInstance.adminClient.GetMetadata(&name, false, 5000)
	if err != nil {
		return err
	}

	if _, exists := describeTopics.Topics[name]; exists {
		return nil // Topic already exists
	}

	topicSpec := kafka.TopicSpecification{
		Topic:             name,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	results, err := kafkaWrapperInstance.adminClient.CreateTopics(ctx, []kafka.TopicSpecification{topicSpec})
	if err != nil {
		return err
	}

	if len(results) > 0 && results[0].Error.Code() != kafka.ErrNoError {
		return errors.New(results[0].Error.String())
	}

	return nil
}

func NewProducer() (*kafka.Producer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaWrapperInstance.brokers})
	if err != nil {
		return nil, err
	}
	return producer, nil
}

func NewConsumer(groupId string) (*kafka.Consumer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": kafkaWrapperInstance.brokers,
		"group.id":          groupId,
		"auto.offset.reset": "earliest",
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		Close()
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}
	return consumer, nil
}

func Close() {
	kafkaWrapperInstance.adminClient.Close()
}
