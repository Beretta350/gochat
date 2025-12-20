package worker

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Beretta350/gochat/internal/app/model"
	"github.com/Beretta350/gochat/internal/app/repository"
	"github.com/Beretta350/gochat/pkg/logger"
	"github.com/Beretta350/gochat/pkg/redisclient"
)

const (
	consumerGroup = "message-workers"
	batchSize     = 100
	batchTimeout  = time.Second * 2
)

// MessageWorker processes messages from Redis Stream
type MessageWorker struct {
	redis      *redisclient.Client
	repo       repository.MessageRepository
	workerID   string
	batchSize  int
	buffer     []*model.ChatMessage
	bufferLock sync.Mutex
	lastFlush  time.Time
}

// NewMessageWorker creates a new message worker (Fx provider)
func NewMessageWorker(redis *redisclient.Client, repo repository.MessageRepository) *MessageWorker {
	logger.Info("Message worker initialized")
	return &MessageWorker{
		redis:     redis,
		repo:      repo,
		workerID:  "worker-1",
		batchSize: batchSize,
		buffer:    make([]*model.ChatMessage, 0, batchSize),
		lastFlush: time.Now(),
	}
}

// Start starts the worker
func (w *MessageWorker) Start(ctx context.Context) {
	if err := w.redis.CreateConsumerGroup(ctx, consumerGroup); err != nil {
		logger.Errorf("Failed to create consumer group: %v", err)
	}

	logger.Infof("Message worker %s started", w.workerID)

	go w.flushTicker(ctx)

	for {
		select {
		case <-ctx.Done():
			w.flush(ctx)
			logger.Infof("Message worker %s stopped", w.workerID)
			return
		default:
			w.processMessages(ctx)
		}
	}
}

func (w *MessageWorker) processMessages(ctx context.Context) {
	streams, err := w.redis.ReadStreamGroup(ctx, consumerGroup, w.workerID, int64(w.batchSize), time.Second)
	if err != nil {
		if err != redis.Nil {
			logger.Errorf("Error reading from stream: %v", err)
		}
		time.Sleep(time.Second)
		return
	}

	for _, stream := range streams {
		for _, msg := range stream.Messages {
			chatMsg := w.parseMessage(msg.Values)
			if chatMsg != nil {
				w.addToBuffer(ctx, chatMsg)
			}

			if err := w.redis.AckMessage(ctx, consumerGroup, msg.ID); err != nil {
				logger.Errorf("Error acknowledging message %s: %v", msg.ID, err)
			}
		}
	}
}

func (w *MessageWorker) parseMessage(values map[string]interface{}) *model.ChatMessage {
	msgJSON, ok := values["data"].(string)
	if !ok {
		return nil
	}

	var msg model.ChatMessage
	if err := json.Unmarshal([]byte(msgJSON), &msg); err != nil {
		logger.Errorf("Error parsing message: %v", err)
		return nil
	}

	return &msg
}

func (w *MessageWorker) addToBuffer(ctx context.Context, msg *model.ChatMessage) {
	w.bufferLock.Lock()
	defer w.bufferLock.Unlock()

	w.buffer = append(w.buffer, msg)

	if len(w.buffer) >= w.batchSize {
		w.flushLocked(ctx)
	}
}

func (w *MessageWorker) flushTicker(ctx context.Context) {
	ticker := time.NewTicker(batchTimeout)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.flush(ctx)
		}
	}
}

func (w *MessageWorker) flush(ctx context.Context) {
	w.bufferLock.Lock()
	defer w.bufferLock.Unlock()
	w.flushLocked(ctx)
}

func (w *MessageWorker) flushLocked(ctx context.Context) {
	if len(w.buffer) == 0 {
		return
	}

	if err := w.repo.SaveBatch(ctx, w.buffer); err != nil {
		logger.Errorf("Error saving batch: %v", err)
		return
	}

	logger.Infof("Worker %s flushed %d messages", w.workerID, len(w.buffer))

	w.buffer = make([]*model.ChatMessage, 0, w.batchSize)
	w.lastFlush = time.Now()
}
