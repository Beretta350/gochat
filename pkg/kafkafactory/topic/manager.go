package kafkafactory

import (
	"context"
	"errors"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	client "github.com/Beretta350/gochat/pkg/kafkafactory/client"
	"github.com/Beretta350/gochat/pkg/logger"
)

var once sync.Once

var topicManagerInstance *topicManager

type TopicManager interface {
	CreateTopic(ctx context.Context, name string) error
	DeleteTopic(ctx context.Context, name string) error
	Close()
}

type topicManager struct {
	adminClient client.AdminClientInterface
	brokers     string
}

func Init(brokers string) {
	once.Do(func() {
		adminClient, err := client.NewAdminClientAdapter(brokers)
		if err != nil {
			logger.Fatal(err)
		}

		topicManagerInstance = &topicManager{
			adminClient: adminClient,
			brokers:     brokers,
		}
	})
	logger.Info("Kafka admin client factory initialized")
}

// For testing purposes
func SetAdminClient(client client.AdminClientInterface) {
	if topicManagerInstance == nil {
		topicManagerInstance = &topicManager{
			adminClient: client,
			brokers:     "test-broker",
		}
	} else {
		topicManagerInstance.adminClient = client
	}
}

func CreateTopic(ctx context.Context, name string) error {
	// Check if the topic already exists
	describeTopics, err := topicManagerInstance.adminClient.GetMetadata(&name, false, 5000)
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

	results, err := topicManagerInstance.adminClient.CreateTopics(ctx, []kafka.TopicSpecification{topicSpec})
	if err != nil {
		return err
	}

	if len(results) > 0 && results[0].Error.Code() != kafka.ErrNoError {
		return errors.New(results[0].Error.String())
	}

	return nil
}

func DeleteTopic(ctx context.Context, name string) error {
	results, err := topicManagerInstance.adminClient.DeleteTopics(ctx, []string{name})
	if err != nil {
		return err
	}

	if len(results) > 0 && results[0].Error.Code() != kafka.ErrNoError {
		return errors.New(results[0].Error.String())
	}

	return nil
}

func Close() {
	topicManagerInstance.adminClient.Close()
}
