package kafkawrapper

import (
	"context"
	"errors"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/Beretta350/gochat/pkg/logger"
)

var once sync.Once

var clientWrapperInstance *clientWrapper

type ClientWrapper interface {
	CreateTopic(ctx context.Context, name string) error
	DeleteTopic(ctx context.Context, name string) error
	Close()
}

type clientWrapper struct {
	adminClient AdminClientInterface
	brokers     string
}

func Init(brokers string) {
	once.Do(func() {
		adminClient, err := NewAdminClientAdapter(brokers)
		if err != nil {
			logger.Fatal(err)
		}

		clientWrapperInstance = &clientWrapper{
			adminClient: adminClient,
			brokers:     brokers,
		}
	})
	logger.Info("Kafka admin client wrapper initialized")
}

// For testing purposes
func SetAdminClient(client AdminClientInterface) {
	if clientWrapperInstance == nil {
		clientWrapperInstance = &clientWrapper{
			adminClient: client,
			brokers:     "test-broker",
		}
	} else {
		clientWrapperInstance.adminClient = client
	}
}

func CreateTopic(ctx context.Context, name string) error {
	// Check if the topic already exists
	describeTopics, err := clientWrapperInstance.adminClient.GetMetadata(&name, false, 5000)
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

	results, err := clientWrapperInstance.adminClient.CreateTopics(ctx, []kafka.TopicSpecification{topicSpec})
	if err != nil {
		return err
	}

	if len(results) > 0 && results[0].Error.Code() != kafka.ErrNoError {
		return errors.New(results[0].Error.String())
	}

	return nil
}

func DeleteTopic(ctx context.Context, name string) error {
	results, err := clientWrapperInstance.adminClient.DeleteTopics(ctx, []string{name})
	if err != nil {
		return err
	}

	if len(results) > 0 && results[0].Error.Code() != kafka.ErrNoError {
		return errors.New(results[0].Error.String())
	}

	return nil
}

func Close() {
	clientWrapperInstance.adminClient.Close()
}
