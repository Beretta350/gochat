package cache

import (
	"testing"

	"github.com/Beretta350/gochat/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetConnectedUserCache(t *testing.T) {
	// Test singleton pattern
	cache1 := GetConnectedUserCache()
	cache2 := GetConnectedUserCache()

	// Both instances should be the same
	assert.Equal(t, cache1, cache2)
}

func TestConnectedUserCache_Add(t *testing.T) {
	// Setup
	cache := GetConnectedUserCache()
	mockConn := mocks.NewConn(t)
	username := "testuser"

	// Test adding a connection
	cache.Add(username, mockConn)

	// Verify the connection was added correctly
	conn := cache.Get(username)
	assert.Equal(t, mockConn, conn)
}

func TestConnectedUserCache_Get(t *testing.T) {
	// Setup
	cache := GetConnectedUserCache()
	mockConn := mocks.NewConn(t)
	username := "testuser"

	// Add a connection
	cache.Add(username, mockConn)

	// Test getting an existing connection
	conn := cache.Get(username)
	assert.Equal(t, mockConn, conn)

	// Test getting a non-existent connection
	nonExistentConn := cache.Get("nonexistentuser")
	assert.Nil(t, nonExistentConn)
}

func TestConnectedUserCache_Remove(t *testing.T) {
	// Setup
	cache := GetConnectedUserCache()
	mockConn := mocks.NewConn(t)
	username := "testuser"

	// Add a connection
	cache.Add(username, mockConn)

	// Verify the connection exists
	conn := cache.Get(username)
	assert.NotNil(t, conn)

	// Test removing the connection
	cache.Remove(username)

	// Verify the connection was removed
	conn = cache.Get(username)
	assert.Nil(t, conn)
}

func TestConnectedUserCache_WebsocketInteraction(t *testing.T) {
	// Setup
	cache := GetConnectedUserCache()
	mockConn := mocks.NewConn(t)
	username := "testuser"

	// Setup message expectations
	testMessage := map[string]interface{}{
		"type": "message",
		"text": "Hello, world!",
	}

	mockConn.On("WriteJSON", testMessage).Return(nil)
	mockConn.On("Close").Return(nil)

	// Add the connection to the cache
	cache.Add(username, mockConn)

	// Get the connection from the cache and send a message
	conn := cache.Get(username)
	err := conn.WriteJSON(testMessage)

	// Verify message was sent through the connection
	assert.NoError(t, err)
	mockConn.AssertCalled(t, "WriteJSON", testMessage)

	// Test connection close
	err = conn.Close()
	assert.NoError(t, err)
	mockConn.AssertCalled(t, "Close")

	// Remove the connection
	cache.Remove(username)
}
