package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/Beretta350/gochat/internal/app/model"
	"github.com/Beretta350/gochat/pkg/kafka_wrapper"
	"github.com/Beretta350/gochat/pkg/logger"
)

type KafkaClient interface {
	ProduceMessage(message model.ChatMessage) error
	ConsumeMessage(ctx context.Context, handler func(model.ChatMessage)) error
	CloseConnection()
}

type kafkaClient struct {
	userToken        string
	userTopicName    string
	userConsumerName string
	producer         *kafka.Producer
	consumer         *kafka.Consumer
}

func NewKafkaClient(ctx context.Context, token string) (KafkaClient, error) {
	userTopicName := token + "-topic"
	userConsumerName := token + "-consumer"

	err := kafka_wrapper.CreateTopic(ctx, userTopicName)
	if err != nil {
		return nil, fmt.Errorf("user topic %s creation error", userTopicName)
	}

	producer, err := kafka_wrapper.NewProducer()
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	consumer, err := kafka_wrapper.NewConsumer(userConsumerName)
	if err != nil {
		producer.Close()
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return &kafkaClient{
		userToken:        token,
		userTopicName:    userTopicName,
		userConsumerName: userConsumerName,
		producer:         producer,
		consumer:         consumer,
	}, nil
}

func (u *kafkaClient) ProduceMessage(message model.ChatMessage) error {
	recipientTopic := message.Recipient + "-topic"

	bytesMsg, err := message.Bytes()
	if err != nil {
		return err
	}

	return u.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &recipientTopic, Partition: kafka.PartitionAny},
		Value:          bytesMsg,
	}, nil)
}

func (u *kafkaClient) ConsumeMessage(ctx context.Context, handler func(model.ChatMessage)) error {
	//TODO: Have a list of user's consumer topic subscriptions (e.g user's topic and groups the user are in)

	//Get topics that the user are subscribed on the database

	//Subscribe to the fetched topics
	err := u.consumer.Subscribe(u.userTopicName, nil)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			logger.Infof("%s context canceled", u.userToken)
			return nil
		default:
			msg, readErr := u.consumer.ReadMessage(3) // Blocking Kafka read
			if readErr != nil {

				var kafkaErr kafka.Error
				if errors.As(readErr, &kafkaErr) && kafkaErr.Code() == kafka.ErrTimedOut {
					continue
				}

				logger.Error("Failed to read Kafka message: ", readErr)
				return readErr
			}

			// Parse and send message to the channel
			var chatMessage model.ChatMessage
			if jsonErr := json.Unmarshal(msg.Value, &chatMessage); jsonErr != nil {
				logger.Error("Failed to unmarshal Kafka message:", jsonErr)
				continue
			}

			handler(chatMessage)
		}
	}
}

func (u *kafkaClient) CloseConnection() {
	u.producer.Close()
	err := u.consumer.Close()
	if err != nil {
		logger.Fatal("Error closing consumer:", err)
	}
	logger.Infof("%s kafka's connection closed", u.userToken)
}
