package model

import (
	"encoding/json"
	"fmt"
	"time"
)

// ChatMessage represents a chat message
type ChatMessage struct {
	ID         string `json:"id,omitempty"`
	Sender     string `json:"sender"`
	Recipient  string `json:"recipient"`
	Content    string `json:"content"`
	Type       string `json:"type,omitempty"` // "text", "image", "file"
	SentAt     int64  `json:"sent_at"`        // Unix timestamp ms
	ReceivedAt *int64 `json:"received_at"`    // Unix timestamp ms (nil if not received yet)
}

// NewChatMessage creates a new chat message with current timestamp
func NewChatMessage(sender, recipient, content string) *ChatMessage {
	return &ChatMessage{
		Sender:    sender,
		Recipient: recipient,
		Content:   content,
		Type:      "text",
		SentAt:    time.Now().UnixMilli(),
	}
}

// MarkReceived sets the received timestamp
func (m *ChatMessage) MarkReceived() {
	now := time.Now().UnixMilli()
	m.ReceivedAt = &now
}

// String returns a string representation of the message
func (m ChatMessage) String() string {
	receivedStr := "nil"
	if m.ReceivedAt != nil {
		receivedStr = fmt.Sprintf("%d", *m.ReceivedAt)
	}
	return fmt.Sprintf(
		"{sender: %s, recipient: %s, content: %s, sent_at: %d, received_at: %s}",
		m.Sender, m.Recipient, m.Content, m.SentAt, receivedStr,
	)
}

// Bytes returns the JSON representation of the message
func (m ChatMessage) Bytes() ([]byte, error) {
	return json.Marshal(m)
}
