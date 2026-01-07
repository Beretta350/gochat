package redisclient

import (
	"context"
	"crypto/tls"

	"github.com/redis/go-redis/v9"

	"github.com/Beretta350/gochat/internal/config"
	"github.com/Beretta350/gochat/pkg/logger"
)

// Client wraps the Redis client
type Client struct {
	rdb *redis.Client
}

// NewClient creates a new Redis client (Fx provider)
func NewClient(cfg *config.Config) (*Client, error) {
	opts := &redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}

	// Enable TLS for cloud Redis providers (Upstash, etc.)
	if cfg.Redis.TLS {
		opts.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	rdb := redis.NewClient(opts)

	// Test connection
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	logger.Info("Redis client connected")
	return &Client{rdb: rdb}, nil
}

// Close closes the Redis client (Fx lifecycle)
func (c *Client) Close() error {
	logger.Info("Closing Redis connection")
	return c.rdb.Close()
}

// Publish publishes a message to a channel
func (c *Client) Publish(ctx context.Context, channel string, message interface{}) error {
	return c.rdb.Publish(ctx, channel, message).Err()
}

// Subscribe subscribes to a channel and returns a PubSub
func (c *Client) Subscribe(ctx context.Context, channel string) *redis.PubSub {
	return c.rdb.Subscribe(ctx, channel)
}

// AddToStream adds a message to the main stream
func (c *Client) AddToStream(ctx context.Context, values map[string]interface{}) (string, error) {
	return c.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: "messages:stream",
		Values: values,
	}).Result()
}

// AddToPending adds a message to a user's pending queue
func (c *Client) AddToPending(ctx context.Context, userID string, messageJSON string) error {
	return c.rdb.RPush(ctx, "pending:"+userID, messageJSON).Err()
}

// GetPendingMessages gets all pending messages for a user and removes them
func (c *Client) GetPendingMessages(ctx context.Context, userID string) ([]string, error) {
	key := "pending:" + userID

	messages, err := c.rdb.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	if len(messages) > 0 {
		c.rdb.Del(ctx, key)
	}

	return messages, nil
}

// CreateConsumerGroup creates a consumer group for the stream
func (c *Client) CreateConsumerGroup(ctx context.Context, group string) error {
	err := c.rdb.XGroupCreateMkStream(ctx, "messages:stream", group, "0").Err()
	if err != nil && err.Error() == "BUSYGROUP Consumer Group name already exists" {
		return nil
	}
	return err
}

// ReadStreamGroup reads from stream as part of a consumer group
func (c *Client) ReadStreamGroup(ctx context.Context, group, consumer string, count int64, block interface{}) ([]redis.XStream, error) {
	return c.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: consumer,
		Streams:  []string{"messages:stream", ">"},
		Count:    count,
		Block:    0,
	}).Result()
}

// AckMessage acknowledges a message
func (c *Client) AckMessage(ctx context.Context, group, id string) error {
	return c.rdb.XAck(ctx, "messages:stream", group, id).Err()
}

// ==================== Online Status Tracking ====================

const onlineUsersKey = "online:users"

// SetUserOnline marks a user as online
func (c *Client) SetUserOnline(ctx context.Context, userID string) error {
	return c.rdb.SAdd(ctx, onlineUsersKey, userID).Err()
}

// SetUserOffline marks a user as offline
func (c *Client) SetUserOffline(ctx context.Context, userID string) error {
	return c.rdb.SRem(ctx, onlineUsersKey, userID).Err()
}

// IsUserOnline checks if a user is online
func (c *Client) IsUserOnline(ctx context.Context, userID string) (bool, error) {
	return c.rdb.SIsMember(ctx, onlineUsersKey, userID).Result()
}

// GetOnlineUsers returns all online user IDs
func (c *Client) GetOnlineUsers(ctx context.Context) ([]string, error) {
	return c.rdb.SMembers(ctx, onlineUsersKey).Result()
}

// GetOnlineUsersFromList checks which users from the list are online
func (c *Client) GetOnlineUsersFromList(ctx context.Context, userIDs []string) ([]string, error) {
	if len(userIDs) == 0 {
		return []string{}, nil
	}

	// Use pipeline to check multiple users efficiently
	pipe := c.rdb.Pipeline()
	cmds := make([]*redis.BoolCmd, len(userIDs))

	for i, userID := range userIDs {
		cmds[i] = pipe.SIsMember(ctx, onlineUsersKey, userID)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	var onlineUsers []string
	for i, cmd := range cmds {
		if cmd.Val() {
			onlineUsers = append(onlineUsers, userIDs[i])
		}
	}

	return onlineUsers, nil
}
