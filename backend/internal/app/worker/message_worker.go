package worker

import (
	"context"
	"strconv"
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
	batchTimeout  = time.Millisecond * 500 // Flush every 500ms
)

// MessageWorker processes messages from Redis Stream and persists to PostgreSQL
type MessageWorker struct {
	redis      *redisclient.Client
	repo       repository.MessageRepository
	workerID   string
	batchSize  int
	buffer     []*model.Message
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
		buffer:    make([]*model.Message, 0, batchSize),
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
			message := w.parseMessage(msg.Values)
			if message != nil {
				w.addToBuffer(ctx, message)
			}

			if err := w.redis.AckMessage(ctx, consumerGroup, msg.ID); err != nil {
				logger.Errorf("Error acknowledging message %s: %v", msg.ID, err)
			}
		}
	}
}

func (w *MessageWorker) parseMessage(values map[string]interface{}) *model.Message {
	// Debug: log raw values
	logger.Infof("Worker parsing message: %+v", values)

	// Get required fields
	id, _ := values["id"].(string)
	conversationID, _ := values["conversation_id"].(string)
	senderID, _ := values["sender_id"].(string)
	content, _ := values["content"].(string)

	// Debug: log parsed fields
	logger.Infof("Parsed - id: %s, conv: %s, sender: %s, content: %s", id, conversationID, senderID, content)

	// Skip if missing required fields
	if id == "" || conversationID == "" || senderID == "" || content == "" {
		logger.Warnf("Skipping message - missing required fields")
		return nil
	}

	// Get optional fields
	msgType, _ := values["type"].(string)
	if msgType == "" {
		msgType = "text"
	}

	// Parse sent_at
	var sentAt time.Time
	if sentAtVal, ok := values["sent_at"]; ok {
		switch v := sentAtVal.(type) {
		case int64:
			sentAt = time.UnixMilli(v)
		case string:
			if ts, err := strconv.ParseInt(v, 10, 64); err == nil {
				sentAt = time.UnixMilli(ts)
			}
		}
	}
	if sentAt.IsZero() {
		sentAt = time.Now()
	}

	return &model.Message{
		ID:             id,
		ConversationID: conversationID,
		SenderID:       senderID,
		Content:        content,
		Type:           model.MessageType(msgType),
		SentAt:         sentAt,
	}
}

func (w *MessageWorker) addToBuffer(ctx context.Context, msg *model.Message) {
	w.bufferLock.Lock()
	defer w.bufferLock.Unlock()

	logger.Infof("Adding message to buffer: conv=%s, sender=%s, content=%s", msg.ConversationID, msg.SenderID, msg.Content)
	w.buffer = append(w.buffer, msg)
	logger.Infof("Buffer size: %d", len(w.buffer))

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

	if err := w.repo.CreateBatch(ctx, w.buffer); err != nil {
		logger.Errorf("Error saving batch: %v", err)
		return
	}

	logger.Infof("Worker %s flushed %d messages to PostgreSQL", w.workerID, len(w.buffer))

	w.buffer = make([]*model.Message, 0, w.batchSize)
	w.lastFlush = time.Now()
}
