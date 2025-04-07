package wsadapter

import (
	"time"

	"github.com/gorilla/websocket"
)

// Conn interface allows flexibility for testing and real connections.
type Conn interface {
	Close() error
	WriteMessage(messageType int, p []byte) error
	ReadMessage() (messageType int, p []byte, err error)
	ReadJSON(v any) error
	WriteJSON(v any) error
}

// NewConn creates a new WebSocket connection with environment-based configuration.
func NewConn(url string, readBufferSize int, writeBufferSize int) (Conn, error) {
	dialer := websocket.Dialer{
		HandshakeTimeout: 45 * time.Second,
		ReadBufferSize:   readBufferSize,
		WriteBufferSize:  writeBufferSize,
	}

	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
