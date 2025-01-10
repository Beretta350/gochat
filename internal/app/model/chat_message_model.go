package model

import (
	"fmt"
	"time"
)

//TODO: Put some validator (required fields and rules)

type ChatMessage struct {
	Sender    string    `json:"sender"`
	Content   string    `json:"content"`
	Recipient string    `json:"recipient"`
	Created   time.Time `json:"created"`
}

func (m ChatMessage) Send() error {
	return nil
}

// Implement the String() method
func (m ChatMessage) String() string {
	return fmt.Sprintf(
		"{sender: %s, content: %s, recipient: %s, created: %s}",
		m.Sender, m.Content, m.Recipient, m.Created.String(),
	)
}
