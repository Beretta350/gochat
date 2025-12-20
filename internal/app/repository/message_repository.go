package repository

import (
	"context"
	"sync"

	"github.com/Beretta350/gochat/internal/app/model"
	"github.com/Beretta350/gochat/pkg/logger"
)

// MessageRepository defines the interface for message persistence
type MessageRepository interface {
	Save(ctx context.Context, msg *model.ChatMessage) error
	SaveBatch(ctx context.Context, msgs []*model.ChatMessage) error
	GetConversation(ctx context.Context, user1, user2 string, limit, offset int) ([]*model.ChatMessage, error)
	GetUndelivered(ctx context.Context, userID string) ([]*model.ChatMessage, error)
	MarkAsDelivered(ctx context.Context, messageID string) error
	MarkAsRead(ctx context.Context, messageID string) error
}

// InMemoryMessageRepository is a temporary in-memory implementation
type InMemoryMessageRepository struct {
	messages []*model.ChatMessage
	mu       sync.RWMutex
}

// NewMessageRepository creates a new message repository (Fx provider)
// Currently returns in-memory, will return Postgres later
func NewMessageRepository() MessageRepository {
	logger.Info("Message repository initialized (in-memory)")
	return &InMemoryMessageRepository{
		messages: make([]*model.ChatMessage, 0),
	}
}

func (r *InMemoryMessageRepository) Save(ctx context.Context, msg *model.ChatMessage) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.messages = append(r.messages, msg)
	return nil
}

func (r *InMemoryMessageRepository) SaveBatch(ctx context.Context, msgs []*model.ChatMessage) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.messages = append(r.messages, msgs...)
	logger.Infof("Saved batch of %d messages", len(msgs))
	return nil
}

func (r *InMemoryMessageRepository) GetConversation(ctx context.Context, user1, user2 string, limit, offset int) ([]*model.ChatMessage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*model.ChatMessage
	for _, msg := range r.messages {
		if (msg.Sender == user1 && msg.Recipient == user2) ||
			(msg.Sender == user2 && msg.Recipient == user1) {
			result = append(result, msg)
		}
	}

	start := offset
	if start > len(result) {
		return []*model.ChatMessage{}, nil
	}
	end := start + limit
	if end > len(result) {
		end = len(result)
	}

	return result[start:end], nil
}

func (r *InMemoryMessageRepository) GetUndelivered(ctx context.Context, userID string) ([]*model.ChatMessage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*model.ChatMessage
	for _, msg := range r.messages {
		if msg.Recipient == userID && msg.ReceivedAt == nil {
			result = append(result, msg)
		}
	}
	return result, nil
}

func (r *InMemoryMessageRepository) MarkAsDelivered(ctx context.Context, messageID string) error {
	return nil
}

func (r *InMemoryMessageRepository) MarkAsRead(ctx context.Context, messageID string) error {
	return nil
}
