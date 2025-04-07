package model

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestChatMessage_Send(t *testing.T) {
	message := ChatMessage{
		Sender:    "user1",
		Content:   "Hello",
		Recipient: "user2",
		Created:   time.Now(),
	}

	err := message.Send()
	assert.NoError(t, err, "Send method should not return an error")
}

func TestChatMessage_String(t *testing.T) {
	// Create a message with fixed timestamp to make testing easier
	fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	message := ChatMessage{
		Sender:    "user1",
		Content:   "Hello",
		Recipient: "user2",
		Created:   fixedTime,
	}

	expected := "{sender: user1, content: Hello, recipient: user2, created: 2023-01-01 12:00:00 +0000 UTC}"
	result := message.String()

	assert.Equal(t, expected, result, "String representation should match expected format")
}

func TestChatMessage_Bytes(t *testing.T) {
	// Create a message with fixed timestamp to make testing easier
	fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	message := ChatMessage{
		Sender:    "user1",
		Content:   "Hello",
		Recipient: "user2",
		Created:   fixedTime,
	}

	bytes, err := message.Bytes()
	assert.NoError(t, err, "Bytes method should not return an error")

	// Test that the bytes can be unmarshaled back to the original message
	var unmarshaled ChatMessage
	err = json.Unmarshal(bytes, &unmarshaled)
	assert.NoError(t, err, "Should be able to unmarshal the bytes back to a ChatMessage")

	assert.Equal(t, message.Sender, unmarshaled.Sender, "Sender should match after marshal/unmarshal")
	assert.Equal(t, message.Content, unmarshaled.Content, "Content should match after marshal/unmarshal")
	assert.Equal(t, message.Recipient, unmarshaled.Recipient, "Recipient should match after marshal/unmarshal")
	assert.Equal(t, message.Created.UTC(), unmarshaled.Created.UTC(), "Created time should match after marshal/unmarshal")

	// Additionally, manually check bytes against expected JSON
	expectedJSON, _ := json.Marshal(message)
	assert.Equal(t, expectedJSON, bytes, "Bytes should match result of json.Marshal")
}
