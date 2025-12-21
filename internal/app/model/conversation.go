package model

import "time"

// ConversationType represents the type of conversation
type ConversationType string

const (
	ConversationTypeDirect ConversationType = "direct"
	ConversationTypeGroup  ConversationType = "group"
)

// Conversation represents a chat conversation
type Conversation struct {
	ID        string           `json:"id"`
	Type      ConversationType `json:"type"`
	Name      *string          `json:"name,omitempty"` // Only for groups
	CreatedBy *string          `json:"created_by,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`

	// Populated fields (not stored directly)
	Participants []Participant `json:"participants,omitempty"`
	LastMessage  *Message      `json:"last_message,omitempty"`
}

// ParticipantRole represents the role in a group conversation
type ParticipantRole string

const (
	ParticipantRoleAdmin  ParticipantRole = "admin"
	ParticipantRoleMember ParticipantRole = "member"
)

// Participant represents a user's participation in a conversation
type Participant struct {
	ConversationID string           `json:"conversation_id"`
	UserID         string           `json:"user_id"`
	Role           *ParticipantRole `json:"role,omitempty"` // NULL for direct, 'admin'/'member' for groups
	JoinedAt       time.Time        `json:"joined_at"`
	LeftAt         *time.Time       `json:"left_at,omitempty"`

	// Populated field
	User *UserResponse `json:"user,omitempty"`
}

// ConversationCreate represents data to create a conversation
type ConversationCreate struct {
	Type           ConversationType `json:"type" validate:"required,oneof=direct group"`
	Name           *string          `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	ParticipantIDs []string         `json:"participant_ids" validate:"required,min=1"`
}

// ConversationResponse represents a conversation in API responses
type ConversationResponse struct {
	ID           string           `json:"id"`
	Type         ConversationType `json:"type"`
	Name         *string          `json:"name,omitempty"`
	Participants []UserResponse   `json:"participants"`
	LastMessage  *Message         `json:"last_message,omitempty"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}
