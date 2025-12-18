package redisclient

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// MainStream is the main message stream for persistence
	MainStream = "messages:stream"
	// PendingPrefix is the prefix for pending messages per user
	PendingPrefix = "pending:"
)

// AddToStream adds a message to the main stream
func AddToStream(ctx context.Context, values map[string]interface{}) (string, error) {
	return client.XAdd(ctx, &redis.XAddArgs{
		Stream: MainStream,
		Values: values,
	}).Result()
}

// AddToPending adds a message to a user's pending queue
func AddToPending(ctx context.Context, userID string, messageJSON string) error {
	key := PendingPrefix + userID
	return client.RPush(ctx, key, messageJSON).Err()
}

// GetPendingMessages gets all pending messages for a user and removes them
func GetPendingMessages(ctx context.Context, userID string) ([]string, error) {
	key := PendingPrefix + userID

	// Get all pending messages
	messages, err := client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	// Delete the pending queue after reading
	if len(messages) > 0 {
		client.Del(ctx, key)
	}

	return messages, nil
}

// ReadStream reads messages from the stream (for workers)
func ReadStream(ctx context.Context, lastID string, count int64, block time.Duration) ([]redis.XMessage, error) {
	streams, err := client.XRead(ctx, &redis.XReadArgs{
		Streams: []string{MainStream, lastID},
		Count:   count,
		Block:   block,
	}).Result()

	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if len(streams) > 0 {
		return streams[0].Messages, nil
	}
	return nil, nil
}

// AckMessage acknowledges a message (for consumer groups)
func AckMessage(ctx context.Context, group, id string) error {
	return client.XAck(ctx, MainStream, group, id).Err()
}

// CreateConsumerGroup creates a consumer group for the stream
func CreateConsumerGroup(ctx context.Context, group string) error {
	err := client.XGroupCreateMkStream(ctx, MainStream, group, "0").Err()
	// Ignore error if group already exists
	if err != nil && err.Error() == "BUSYGROUP Consumer Group name already exists" {
		return nil
	}
	return err
}

// ReadStreamGroup reads from stream as part of a consumer group
func ReadStreamGroup(ctx context.Context, group, consumer string, count int64, block time.Duration) ([]redis.XMessage, error) {
	streams, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: consumer,
		Streams:  []string{MainStream, ">"},
		Count:    count,
		Block:    block,
	}).Result()

	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if len(streams) > 0 {
		return streams[0].Messages, nil
	}
	return nil, nil
}
