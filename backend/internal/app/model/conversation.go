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

// Participant represents a user's participation in a conversation (internal)
type Participant struct {
	ConversationID string
	UserID         string
	Role           *ParticipantRole
	JoinedAt       time.Time
	LeftAt         *time.Time

	// Populated field
	User *UserResponse
}

// ParticipantResponse is the flattened API response for a participant
type ParticipantResponse struct {
	ID        string           `json:"id"`
	Email     string           `json:"email"`
	Username  string           `json:"username"`
	IsActive  bool             `json:"is_active"`
	Role      *ParticipantRole `json:"role,omitempty"`
	JoinedAt  time.Time        `json:"joined_at"`
	CreatedAt time.Time        `json:"created_at"`
}

// ToResponse converts Participant to flattened response
func (p *Participant) ToResponse() *ParticipantResponse {
	if p.User == nil {
		return nil
	}
	return &ParticipantResponse{
		ID:        p.User.ID,
		Email:     p.User.Email,
		Username:  p.User.Username,
		IsActive:  p.User.IsActive,
		Role:      p.Role,
		JoinedAt:  p.JoinedAt,
		CreatedAt: p.User.CreatedAt,
	}
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
