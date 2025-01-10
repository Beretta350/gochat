package model

import (
	"fmt"
	"time"
)

//TODO: Put a validator to handshake for require username

type HandshakeMessage struct {
	Username string    `json:"username"`
	Created  time.Time `json:"created"`
}

func (m HandshakeMessage) Send() error {
	return nil
}

// Implement the String() method
func (m HandshakeMessage) String() string {
	return fmt.Sprintf(
		"{username: %s, created: %s}",
		m.Username, m.Created.String(),
	)
}
