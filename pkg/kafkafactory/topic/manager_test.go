package kafkafactory

import (
	"context"
	"errors"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAdminClient is a mock implementation of AdminClientInterface
type MockAdminClient struct {
	mock.Mock
}

func (m *MockAdminClient) GetMetadata(topic *string, allTopics bool, timeoutMs int) (*kafka.Metadata, error) {
	args := m.Called(topic, allTopics, timeoutMs)
	return args.Get(0).(*kafka.Metadata), args.Error(1)
}

func (m *MockAdminClient) CreateTopics(ctx context.Context, topics []kafka.TopicSpecification, options ...kafka.CreateTopicsAdminOption) ([]kafka.TopicResult, error) {
	args := m.Called(ctx, topics)
	return args.Get(0).([]kafka.TopicResult), args.Error(1)
}

func (m *MockAdminClient) DeleteTopics(ctx context.Context, topics []string, options ...kafka.DeleteTopicsAdminOption) ([]kafka.TopicResult, error) {
	args := m.Called(ctx, topics)
	return args.Get(0).([]kafka.TopicResult), args.Error(1)
}

func (m *MockAdminClient) Close() {
	m.Called()
}

func TestCreateTopic_TopicAlreadyExists(t *testing.T) {
	// Arrange
	mockClient := new(MockAdminClient)
	SetAdminClient(mockClient)

	topicName := "test-topic"
	ctx := context.Background()

	// Mock the GetMetadata call to return a metadata that indicates the topic exists
	metadata := &kafka.Metadata{
		Topics: map[string]kafka.TopicMetadata{
			topicName: {},
		},
	}
	mockClient.On("GetMetadata", &topicName, false, 5000).Return(metadata, nil)

	// Act
	err := CreateTopic(ctx, topicName)

	// Assert
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestCreateTopic_Success(t *testing.T) {
	// Arrange
	mockClient := new(MockAdminClient)
	SetAdminClient(mockClient)

	topicName := "test-topic"
	ctx := context.Background()

	// Mock the GetMetadata call to return a metadata that indicates the topic doesn't exist
	metadata := &kafka.Metadata{
		Topics: map[string]kafka.TopicMetadata{},
	}
	mockClient.On("GetMetadata", &topicName, false, 5000).Return(metadata, nil)

	// Mock the CreateTopics call
	topicSpec := kafka.TopicSpecification{
		Topic:             topicName,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
	results := []kafka.TopicResult{
		{
			Topic: topicName,
			Error: kafka.NewError(kafka.ErrNoError, "", false),
		},
	}
	mockClient.On("CreateTopics", ctx, []kafka.TopicSpecification{topicSpec}).Return(results, nil)

	// Act
	err := CreateTopic(ctx, topicName)

	// Assert
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestCreateTopic_MetadataError(t *testing.T) {
	// Arrange
	mockClient := new(MockAdminClient)
	SetAdminClient(mockClient)

	topicName := "test-topic"
	ctx := context.Background()
	expectedErr := errors.New("metadata error")

	// Mock the GetMetadata call to return an error
	mockClient.On("GetMetadata", &topicName, false, 5000).Return(&kafka.Metadata{}, expectedErr)

	// Act
	err := CreateTopic(ctx, topicName)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockClient.AssertExpectations(t)
}

func TestCreateTopic_CreateError(t *testing.T) {
	// Arrange
	mockClient := new(MockAdminClient)
	SetAdminClient(mockClient)

	topicName := "test-topic"
	ctx := context.Background()
	expectedErr := errors.New("create error")

	// Mock the GetMetadata call to return a metadata that indicates the topic doesn't exist
	metadata := &kafka.Metadata{
		Topics: map[string]kafka.TopicMetadata{},
	}
	mockClient.On("GetMetadata", &topicName, false, 5000).Return(metadata, nil)

	// Mock the CreateTopics call to return an error
	topicSpec := kafka.TopicSpecification{
		Topic:             topicName,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
	mockClient.On("CreateTopics", ctx, []kafka.TopicSpecification{topicSpec}).Return([]kafka.TopicResult{}, expectedErr)

	// Act
	err := CreateTopic(ctx, topicName)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockClient.AssertExpectations(t)
}

func TestCreateTopic_TopicError(t *testing.T) {
	// Arrange
	mockClient := new(MockAdminClient)
	SetAdminClient(mockClient)

	topicName := "test-topic"
	ctx := context.Background()
	kafkaErr := kafka.NewError(kafka.ErrTopicAlreadyExists, "topic already exists", false)

	// Mock the GetMetadata call to return a metadata that indicates the topic doesn't exist
	metadata := &kafka.Metadata{
		Topics: map[string]kafka.TopicMetadata{},
	}
	mockClient.On("GetMetadata", &topicName, false, 5000).Return(metadata, nil)

	// Mock the CreateTopics call to return a topic-level error
	topicSpec := kafka.TopicSpecification{
		Topic:             topicName,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
	results := []kafka.TopicResult{
		{
			Topic: topicName,
			Error: kafkaErr,
		},
	}
	mockClient.On("CreateTopics", ctx, []kafka.TopicSpecification{topicSpec}).Return(results, nil)

	// Act
	err := CreateTopic(ctx, topicName)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, kafkaErr.String(), err.Error())
	mockClient.AssertExpectations(t)
}

func TestDeleteTopic_Success(t *testing.T) {
	// Arrange
	mockClient := new(MockAdminClient)
	SetAdminClient(mockClient)

	topicName := "test-topic"
	ctx := context.Background()

	// Mock the DeleteTopics call
	results := []kafka.TopicResult{
		{
			Topic: topicName,
			Error: kafka.NewError(kafka.ErrNoError, "", false),
		},
	}
	mockClient.On("DeleteTopics", ctx, []string{topicName}).Return(results, nil)

	// Act
	err := DeleteTopic(ctx, topicName)

	// Assert
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestDeleteTopic_DeleteError(t *testing.T) {
	// Arrange
	mockClient := new(MockAdminClient)
	SetAdminClient(mockClient)

	topicName := "test-topic"
	ctx := context.Background()
	expectedErr := errors.New("delete error")

	// Mock the DeleteTopics call to return an error
	mockClient.On("DeleteTopics", ctx, []string{topicName}).Return([]kafka.TopicResult{}, expectedErr)

	// Act
	err := DeleteTopic(ctx, topicName)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockClient.AssertExpectations(t)
}

func TestDeleteTopic_TopicError(t *testing.T) {
	// Arrange
	mockClient := new(MockAdminClient)
	SetAdminClient(mockClient)

	topicName := "test-topic"
	ctx := context.Background()
	kafkaErr := kafka.NewError(kafka.ErrUnknownTopic, "unknown topic", false)

	// Mock the DeleteTopics call to return a topic-level error
	results := []kafka.TopicResult{
		{
			Topic: topicName,
			Error: kafkaErr,
		},
	}
	mockClient.On("DeleteTopics", ctx, []string{topicName}).Return(results, nil)

	// Act
	err := DeleteTopic(ctx, topicName)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, kafkaErr.String(), err.Error())
	mockClient.AssertExpectations(t)
}
