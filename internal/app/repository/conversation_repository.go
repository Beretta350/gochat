package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/Beretta350/gochat/internal/app/model"
	"github.com/Beretta350/gochat/pkg/logger"
	"github.com/Beretta350/gochat/pkg/postgres"
)

var (
	ErrConversationNotFound = errors.New("conversation not found")
)

// ConversationRepository defines the interface for conversation persistence
type ConversationRepository interface {
	Create(ctx context.Context, conv *model.Conversation, participantIDs []string) error
	GetByID(ctx context.Context, id string) (*model.Conversation, error)
	GetByUserID(ctx context.Context, userID string) ([]model.Conversation, error)
	FindDirectConversation(ctx context.Context, userID1, userID2 string) (*model.Conversation, error)
	AddParticipant(ctx context.Context, conversationID, userID string, role *model.ParticipantRole) error
	RemoveParticipant(ctx context.Context, conversationID, userID string) error
	GetParticipants(ctx context.Context, conversationID string) ([]model.Participant, error)
}

// PostgresConversationRepository implements ConversationRepository with PostgreSQL
type PostgresConversationRepository struct {
	db *postgres.Client
}

// NewConversationRepository creates a new conversation repository (Fx provider)
func NewConversationRepository(db *postgres.Client) ConversationRepository {
	logger.Info("Conversation repository initialized")
	return &PostgresConversationRepository{db: db}
}

func (r *PostgresConversationRepository) Create(ctx context.Context, conv *model.Conversation, participantIDs []string) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Create conversation
	query := `
		INSERT INTO conversations (type, name, created_by)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	err = tx.QueryRow(ctx, query, conv.Type, conv.Name, conv.CreatedBy).
		Scan(&conv.ID, &conv.CreatedAt, &conv.UpdatedAt)
	if err != nil {
		return err
	}

	// Add participants
	for _, userID := range participantIDs {
		var role *model.ParticipantRole
		if conv.Type == model.ConversationTypeGroup {
			if conv.CreatedBy != nil && *conv.CreatedBy == userID {
				adminRole := model.ParticipantRoleAdmin
				role = &adminRole
			} else {
				memberRole := model.ParticipantRoleMember
				role = &memberRole
			}
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO conversation_participants (conversation_id, user_id, role)
			VALUES ($1, $2, $3)
		`, conv.ID, userID, role)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *PostgresConversationRepository) GetByID(ctx context.Context, id string) (*model.Conversation, error) {
	query := `
		SELECT id, type, name, created_by, created_at, updated_at
		FROM conversations
		WHERE id = $1
	`
	var conv model.Conversation
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&conv.ID,
		&conv.Type,
		&conv.Name,
		&conv.CreatedBy,
		&conv.CreatedAt,
		&conv.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrConversationNotFound
	}
	if err != nil {
		return nil, err
	}
	return &conv, nil
}

func (r *PostgresConversationRepository) GetByUserID(ctx context.Context, userID string) ([]model.Conversation, error) {
	query := `
		SELECT c.id, c.type, c.name, c.created_by, c.created_at, c.updated_at
		FROM conversations c
		JOIN conversation_participants cp ON c.id = cp.conversation_id
		WHERE cp.user_id = $1 AND cp.left_at IS NULL
		ORDER BY c.updated_at DESC
	`
	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []model.Conversation
	for rows.Next() {
		var conv model.Conversation
		err := rows.Scan(
			&conv.ID,
			&conv.Type,
			&conv.Name,
			&conv.CreatedBy,
			&conv.CreatedAt,
			&conv.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, conv)
	}
	return conversations, nil
}

func (r *PostgresConversationRepository) FindDirectConversation(ctx context.Context, userID1, userID2 string) (*model.Conversation, error) {
	query := `
		SELECT c.id, c.type, c.name, c.created_by, c.created_at, c.updated_at
		FROM conversations c
		WHERE c.type = 'direct'
		AND EXISTS (
			SELECT 1 FROM conversation_participants cp1
			WHERE cp1.conversation_id = c.id AND cp1.user_id = $1 AND cp1.left_at IS NULL
		)
		AND EXISTS (
			SELECT 1 FROM conversation_participants cp2
			WHERE cp2.conversation_id = c.id AND cp2.user_id = $2 AND cp2.left_at IS NULL
		)
		LIMIT 1
	`
	var conv model.Conversation
	err := r.db.Pool.QueryRow(ctx, query, userID1, userID2).Scan(
		&conv.ID,
		&conv.Type,
		&conv.Name,
		&conv.CreatedBy,
		&conv.CreatedAt,
		&conv.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrConversationNotFound
	}
	if err != nil {
		return nil, err
	}
	return &conv, nil
}

func (r *PostgresConversationRepository) AddParticipant(ctx context.Context, conversationID, userID string, role *model.ParticipantRole) error {
	query := `
		INSERT INTO conversation_participants (conversation_id, user_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (conversation_id, user_id) DO UPDATE SET left_at = NULL, role = $3
	`
	_, err := r.db.Pool.Exec(ctx, query, conversationID, userID, role)
	return err
}

func (r *PostgresConversationRepository) RemoveParticipant(ctx context.Context, conversationID, userID string) error {
	query := `
		UPDATE conversation_participants
		SET left_at = NOW()
		WHERE conversation_id = $1 AND user_id = $2
	`
	_, err := r.db.Pool.Exec(ctx, query, conversationID, userID)
	return err
}

func (r *PostgresConversationRepository) GetParticipants(ctx context.Context, conversationID string) ([]model.Participant, error) {
	query := `
		SELECT cp.conversation_id, cp.user_id, cp.role, cp.joined_at, cp.left_at,
		       u.id, u.email, u.username, u.is_active, u.created_at
		FROM conversation_participants cp
		JOIN users u ON cp.user_id = u.id
		WHERE cp.conversation_id = $1 AND cp.left_at IS NULL
	`
	rows, err := r.db.Pool.Query(ctx, query, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []model.Participant
	for rows.Next() {
		var p model.Participant
		var u model.UserResponse
		err := rows.Scan(
			&p.ConversationID,
			&p.UserID,
			&p.Role,
			&p.JoinedAt,
			&p.LeftAt,
			&u.ID,
			&u.Email,
			&u.Username,
			&u.IsActive,
			&u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		p.User = &u
		participants = append(participants, p)
	}
	return participants, nil
}
