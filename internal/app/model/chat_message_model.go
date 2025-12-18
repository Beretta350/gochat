package model

import (
	"encoding/json"
	"fmt"
	"time"
)

// ChatMessage represents a chat message
type ChatMessage struct {
	ID        string    `json:"id,omitempty"`
	Sender    string    `json:"sender"`
	Recipient string    `json:"recipient"`
	Content   string    `json:"content"`
	Type      string    `json:"type,omitempty"` // "text", "image", "file"
	CreatedAt time.Time `json:"created_at"`
}

// NewChatMessage creates a new chat message with current timestamp
func NewChatMessage(sender, recipient, content string) *ChatMessage {
	return &ChatMessage{
		Sender:    sender,
		Recipient: recipient,
		Content:   content,
		Type:      "text",
		CreatedAt: time.Now(),
	}
}

// String returns a string representation of the message
func (m ChatMessage) String() string {
	return fmt.Sprintf(
		"{sender: %s, recipient: %s, content: %s, created: %s}",
		m.Sender, m.Recipient, m.Content, m.CreatedAt.Format(time.RFC3339),
	)
}

// Bytes returns the JSON representation of the message
func (m ChatMessage) Bytes() ([]byte, error) {
	return json.Marshal(m)
}
