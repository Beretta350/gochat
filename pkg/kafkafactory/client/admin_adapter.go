package kafkafactory

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// AdminClientInterface is an interface to make Kafka AdminClient mockable
type AdminClientInterface interface {
	GetMetadata(topic *string, allTopics bool, timeoutMs int) (*kafka.Metadata, error)
	CreateTopics(ctx context.Context, topics []kafka.TopicSpecification, options ...kafka.CreateTopicsAdminOption) ([]kafka.TopicResult, error)
	DeleteTopics(ctx context.Context, topics []string, options ...kafka.DeleteTopicsAdminOption) ([]kafka.TopicResult, error)
	Close()
}

// AdminClientAdapter adapts the Kafka AdminClient to implement AdminClientInterface
type AdminClientAdapter struct {
	client *kafka.AdminClient
}

// NewAdminClientAdapter creates a new AdminClientAdapter
func NewAdminClientAdapter(brokers string) (*AdminClientAdapter, error) {
	client, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": brokers})
	if err != nil {
		return nil, err
	}
	return &AdminClientAdapter{client: client}, nil
}

// GetMetadata implements AdminClientInterface
func (a *AdminClientAdapter) GetMetadata(topic *string, allTopics bool, timeoutMs int) (*kafka.Metadata, error) {
	return a.client.GetMetadata(topic, allTopics, timeoutMs)
}

// CreateTopics implements AdminClientInterface
func (a *AdminClientAdapter) CreateTopics(ctx context.Context, topics []kafka.TopicSpecification, options ...kafka.CreateTopicsAdminOption) ([]kafka.TopicResult, error) {
	return a.client.CreateTopics(ctx, topics, options...)
}

// DeleteTopics implements AdminClientInterface
func (a *AdminClientAdapter) DeleteTopics(ctx context.Context, topics []string, options ...kafka.DeleteTopicsAdminOption) ([]kafka.TopicResult, error) {
	return a.client.DeleteTopics(ctx, topics, options...)
}

// Close implements AdminClientInterface
func (a *AdminClientAdapter) Close() {
	a.client.Close()
}
