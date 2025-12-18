package repository

import (
	"context"

	"github.com/Beretta350/gochat/internal/app/model"
)

// MessageRepository defines the interface for message persistence
// This will be implemented by PostgreSQL later
type MessageRepository interface {
	// Save persists a single message
	Save(ctx context.Context, msg *model.ChatMessage) error

	// SaveBatch persists multiple messages at once (for performance)
	SaveBatch(ctx context.Context, msgs []*model.ChatMessage) error

	// GetConversation retrieves messages between two users with pagination
	GetConversation(ctx context.Context, user1, user2 string, limit, offset int) ([]*model.ChatMessage, error)

	// GetUndelivered retrieves all undelivered messages for a user
	GetUndelivered(ctx context.Context, userID string) ([]*model.ChatMessage, error)

	// MarkAsDelivered marks a message as delivered
	MarkAsDelivered(ctx context.Context, messageID string) error

	// MarkAsRead marks a message as read
	MarkAsRead(ctx context.Context, messageID string) error
}

// InMemoryMessageRepository is a temporary in-memory implementation
// Will be replaced by PostgresMessageRepository
type InMemoryMessageRepository struct {
	messages []*model.ChatMessage
}

// NewInMemoryMessageRepository creates a new in-memory repository
func NewInMemoryMessageRepository() MessageRepository {
	return &InMemoryMessageRepository{
		messages: make([]*model.ChatMessage, 0),
	}
}

func (r *InMemoryMessageRepository) Save(ctx context.Context, msg *model.ChatMessage) error {
	r.messages = append(r.messages, msg)
	return nil
}

func (r *InMemoryMessageRepository) SaveBatch(ctx context.Context, msgs []*model.ChatMessage) error {
	r.messages = append(r.messages, msgs...)
	return nil
}

func (r *InMemoryMessageRepository) GetConversation(ctx context.Context, user1, user2 string, limit, offset int) ([]*model.ChatMessage, error) {
	var result []*model.ChatMessage
	for _, msg := range r.messages {
		if (msg.Sender == user1 && msg.Recipient == user2) ||
			(msg.Sender == user2 && msg.Recipient == user1) {
			result = append(result, msg)
		}
	}

	// Apply pagination
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
	var result []*model.ChatMessage
	for _, msg := range r.messages {
		if msg.Recipient == userID && msg.ReceivedAt == nil {
			result = append(result, msg)
		}
	}
	return result, nil
}

func (r *InMemoryMessageRepository) MarkAsDelivered(ctx context.Context, messageID string) error {
	// In-memory: no-op for now
	return nil
}

func (r *InMemoryMessageRepository) MarkAsRead(ctx context.Context, messageID string) error {
	// In-memory: no-op for now
	return nil
}
