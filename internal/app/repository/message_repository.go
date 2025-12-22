package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/Beretta350/gochat/internal/app/model"
	"github.com/Beretta350/gochat/pkg/logger"
	"github.com/Beretta350/gochat/pkg/postgres"
)

var (
	ErrMessageNotFound = errors.New("message not found")
)

// MessageRepository defines the interface for message persistence
type MessageRepository interface {
	Create(ctx context.Context, msg *model.Message) error
	CreateBatch(ctx context.Context, msgs []*model.Message) error
	GetByID(ctx context.Context, id string) (*model.Message, error)
	GetByConversation(ctx context.Context, conversationID string, cursor *time.Time, limit int) (*model.MessagesPage, error)
}

// PostgresMessageRepository implements MessageRepository with PostgreSQL
type PostgresMessageRepository struct {
	db *postgres.Client
}

// NewMessageRepository creates a new message repository (Fx provider)
func NewMessageRepository(db *postgres.Client) MessageRepository {
	logger.Info("Message repository initialized")
	return &PostgresMessageRepository{db: db}
}

func (r *PostgresMessageRepository) Create(ctx context.Context, msg *model.Message) error {
	query := `
		INSERT INTO messages (conversation_id, sender_id, content, type, sent_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	return r.db.Pool.QueryRow(ctx, query,
		msg.ConversationID,
		msg.SenderID,
		msg.Content,
		msg.Type,
		msg.SentAt,
	).Scan(&msg.ID)
}

func (r *PostgresMessageRepository) CreateBatch(ctx context.Context, msgs []*model.Message) error {
	if len(msgs) == 0 {
		return nil
	}

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	for _, msg := range msgs {
		_, err = tx.Exec(ctx, `
			INSERT INTO messages (conversation_id, sender_id, content, type, sent_at)
			VALUES ($1, $2, $3, $4, $5)
		`,
			msg.ConversationID,
			msg.SenderID,
			msg.Content,
			msg.Type,
			msg.SentAt,
		)
		if err != nil {
			return err
		}
	}

	logger.Infof("Saved batch of %d messages", len(msgs))
	return tx.Commit(ctx)
}

func (r *PostgresMessageRepository) GetByID(ctx context.Context, id string) (*model.Message, error) {
	query := `
		SELECT m.id, m.conversation_id, m.sender_id, u.username, m.content, m.type, m.sent_at
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.id = $1
	`
	var msg model.Message

	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&msg.ID,
		&msg.ConversationID,
		&msg.SenderID,
		&msg.SenderUsername,
		&msg.Content,
		&msg.Type,
		&msg.SentAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrMessageNotFound
	}
	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func (r *PostgresMessageRepository) GetByConversation(ctx context.Context, conversationID string, cursor *time.Time, limit int) (*model.MessagesPage, error) {
	// Default limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	// Build query based on cursor
	var query string
	var args []interface{}

	if cursor != nil {
		query = `
			SELECT m.id, m.conversation_id, m.sender_id, u.username, m.content, m.type, m.sent_at
			FROM messages m
			JOIN users u ON m.sender_id = u.id
			WHERE m.conversation_id = $1 AND m.sent_at < $2
			ORDER BY m.sent_at DESC
			LIMIT $3
		`
		args = []interface{}{conversationID, cursor, limit + 1} // +1 to check if there are more
	} else {
		query = `
			SELECT m.id, m.conversation_id, m.sender_id, u.username, m.content, m.type, m.sent_at
			FROM messages m
			JOIN users u ON m.sender_id = u.id
			WHERE m.conversation_id = $1
			ORDER BY m.sent_at DESC
			LIMIT $2
		`
		args = []interface{}{conversationID, limit + 1}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []model.Message
	for rows.Next() {
		var msg model.Message

		err := rows.Scan(
			&msg.ID,
			&msg.ConversationID,
			&msg.SenderID,
			&msg.SenderUsername,
			&msg.Content,
			&msg.Type,
			&msg.SentAt,
		)
		if err != nil {
			return nil, err
		}

		messages = append(messages, msg)
	}

	// Check if there are more messages
	hasMore := len(messages) > limit
	if hasMore {
		messages = messages[:limit]
	}

	// Set next cursor
	var nextCursor *string
	if hasMore && len(messages) > 0 {
		cursorStr := messages[len(messages)-1].SentAt.Format(time.RFC3339Nano)
		nextCursor = &cursorStr
	}

	return &model.MessagesPage{
		Messages:   messages,
		HasMore:    hasMore,
		NextCursor: nextCursor,
	}, nil
}
