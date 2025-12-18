package worker

import (
	"context"
	"encoding/json"
	"sync"
	"time"

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

// MessageWorker processes messages from Redis Stream and saves to database
type MessageWorker struct {
	repo       repository.MessageRepository
	workerID   string
	batchSize  int
	buffer     []*model.ChatMessage
	bufferLock sync.Mutex
	lastFlush  time.Time
}

// NewMessageWorker creates a new message worker
func NewMessageWorker(repo repository.MessageRepository, workerID string) *MessageWorker {
	return &MessageWorker{
		repo:      repo,
		workerID:  workerID,
		batchSize: batchSize,
		buffer:    make([]*model.ChatMessage, 0, batchSize),
		lastFlush: time.Now(),
	}
}

// Start starts the worker
func (w *MessageWorker) Start(ctx context.Context) {
	// Create consumer group if not exists
	if err := redisclient.CreateConsumerGroup(ctx, consumerGroup); err != nil {
		logger.Errorf("Failed to create consumer group: %v", err)
	}

	logger.Infof("Message worker %s started", w.workerID)

	// Start flush ticker
	go w.flushTicker(ctx)

	// Process messages
	for {
		select {
		case <-ctx.Done():
			w.flush(ctx) // Final flush before shutdown
			logger.Infof("Message worker %s stopped", w.workerID)
			return
		default:
			w.processMessages(ctx)
		}
	}
}

func (w *MessageWorker) processMessages(ctx context.Context) {
	messages, err := redisclient.ReadStreamGroup(ctx, consumerGroup, w.workerID, int64(w.batchSize), time.Second)
	if err != nil {
		logger.Errorf("Error reading from stream: %v", err)
		time.Sleep(time.Second)
		return
	}

	for _, msg := range messages {
		chatMsg := w.parseMessage(msg.Values)
		if chatMsg != nil {
			w.addToBuffer(ctx, chatMsg)
		}

		// Acknowledge message
		if err := redisclient.AckMessage(ctx, consumerGroup, msg.ID); err != nil {
			logger.Errorf("Error acknowledging message %s: %v", msg.ID, err)
		}
	}
}

func (w *MessageWorker) parseMessage(values map[string]interface{}) *model.ChatMessage {
	// Parse from stream values
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

	// Flush if buffer is full
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

	// Save batch to repository (will be Postgres later)
	if err := w.repo.SaveBatch(ctx, w.buffer); err != nil {
		logger.Errorf("Error saving batch: %v", err)
		return
	}

	logger.Infof("Worker %s flushed %d messages to storage", w.workerID, len(w.buffer))

	// Clear buffer
	w.buffer = make([]*model.ChatMessage, 0, w.batchSize)
	w.lastFlush = time.Now()
}
