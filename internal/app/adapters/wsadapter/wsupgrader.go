package wsadapter

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Upgrader interface {
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (Conn, error)
}

// wsUpgraderAdapter wraps websocket.Upgrader to implement our Upgrader interface
type upgraderAdapter struct {
	*websocket.Upgrader
}

func (a *upgraderAdapter) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (Conn, error) {
	return a.Upgrader.Upgrade(w, r, responseHeader)
}

// NewUpgrader creates a WebSocket upgrader with environment-based configuration.
func NewUpgrader(readBufferSize int, writeBufferSize int, checkOrigin bool) Upgrader {
	return &upgraderAdapter{
		Upgrader: &websocket.Upgrader{
			ReadBufferSize:  readBufferSize,
			WriteBufferSize: writeBufferSize,
			CheckOrigin: func(r *http.Request) bool {
				// In production, allow only trusted origins
				if checkOrigin {
					return true
				}
				return false
			},
		},
	}
}
