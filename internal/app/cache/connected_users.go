package cache

import (
	"sync"

	"github.com/gorilla/websocket"
)

var once sync.Once
var instance *connectedUserCache

type WebsocketConnectionCache interface {
	Get(username string) *websocket.Conn
	Add(username string, conn *websocket.Conn)
	Remove(username string)
}

type connectedUserCache struct {
	connectedUsers sync.Map
}

// GetConnectedUserCache Singleton users cache constructor
func GetConnectedUserCache() WebsocketConnectionCache {
	once.Do(func() {
		instance = &connectedUserCache{}
	})
	return instance
}

func (c *connectedUserCache) Get(username string) *websocket.Conn {
	conn, ok := c.connectedUsers.Load(username)
	if ok {
		return conn.(*websocket.Conn)
	}
	return nil
}

// Add adds a WebSocket connection to the cache.
func (c *connectedUserCache) Add(username string, conn *websocket.Conn) {
	c.connectedUsers.Store(username, conn)
}

// Remove removes a WebSocket connection from the cache.
func (c *connectedUserCache) Remove(username string) {
	c.connectedUsers.Delete(username)
}
