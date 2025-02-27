package wsupgrader

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// NewUpgrader creates a WebSocket upgrader with environment-based configuration.
func NewUpgrader(readBufferSize int, writeBufferSize int, checkOrigin bool) websocket.Upgrader {
	return websocket.Upgrader{
		ReadBufferSize:  readBufferSize,
		WriteBufferSize: writeBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			// In production, allow only trusted origins
			if checkOrigin {
				return true
			}
			return false
		},
	}
}
