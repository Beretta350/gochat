package redisclient

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"

	"github.com/Beretta350/gochat/pkg/logger"
)

var (
	client *redis.Client
	once   sync.Once
)

// Config holds Redis configuration
type Config struct {
	Addr     string
	Password string
	DB       int
}

// Init initializes the Redis client
func Init(cfg Config) error {
	var initErr error

	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     cfg.Addr,
			Password: cfg.Password,
			DB:       cfg.DB,
		})

		// Test connection
		ctx := context.Background()
		if err := client.Ping(ctx).Err(); err != nil {
			initErr = err
			return
		}

		logger.Info("Redis client connected successfully")
	})

	return initErr
}

// GetClient returns the Redis client instance
func GetClient() *redis.Client {
	return client
}

// Publish publishes a message to a channel
func Publish(ctx context.Context, channel string, message interface{}) error {
	return client.Publish(ctx, channel, message).Err()
}

// Subscribe subscribes to a channel and returns a PubSub
func Subscribe(ctx context.Context, channel string) *redis.PubSub {
	return client.Subscribe(ctx, channel)
}

// Close closes the Redis client
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}
