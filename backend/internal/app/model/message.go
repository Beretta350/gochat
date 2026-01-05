package model

import "time"

// MessageType represents the type of message
type MessageType string

const (
	MessageTypeText  MessageType = "text"
	MessageTypeImage MessageType = "image"
	MessageTypeFile  MessageType = "file"
	MessageTypeAudio MessageType = "audio"
)

// Message represents a chat message
type Message struct {
	ID             string      `json:"id"`
	ConversationID string      `json:"conversation_id"`
	SenderID       string      `json:"sender_id"`
	SenderUsername string      `json:"sender_username,omitempty"`
	Content        string      `json:"content"`
	Type           MessageType `json:"type"`
	SentAt         time.Time   `json:"sent_at"`
}

// MessageCreate represents data to create/send a message
type MessageCreate struct {
	ConversationID string `json:"conversation_id" validate:"required,uuid"`
	Content        string `json:"content" validate:"required,min=1"`
	Type           string `json:"type,omitempty" validate:"omitempty,oneof=text image file audio"`
}

// WebSocketMessage represents a message received via WebSocket
// (simpler format for real-time messaging)
type WebSocketMessage struct {
	RecipientID    string `json:"recipient_id,omitempty"`    // For direct messages
	ConversationID string `json:"conversation_id,omitempty"` // For group messages
	Content        string `json:"content" validate:"required"`
	Type           string `json:"type,omitempty"`
}

// MessageResponse represents a message with sender info
type MessageResponse struct {
	ID             string        `json:"id"`
	ConversationID string        `json:"conversation_id"`
	Sender         *UserResponse `json:"sender"`
	Content        string        `json:"content"`
	Type           MessageType   `json:"type"`
	SentAt         time.Time     `json:"sent_at"`
}

// MessagesPage represents a paginated list of messages
type MessagesPage struct {
	Messages   []Message `json:"messages"`
	HasMore    bool      `json:"has_more"`
	NextCursor *string   `json:"next_cursor,omitempty"` // sent_at of last message
}
