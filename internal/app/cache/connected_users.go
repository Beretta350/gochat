package cache

import (
	"sync"

	"github.com/Beretta350/gochat/pkg/logger"

	"github.com/Beretta350/gochat/internal/app/adapters/wsadapter"
)

var once sync.Once
var instance *connectedUserCache

type WebsocketConnectionCache interface {
	Get(username string) wsadapter.Conn
	Add(username string, conn wsadapter.Conn)
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

func (c *connectedUserCache) Get(username string) wsadapter.Conn {
	conn, ok := c.connectedUsers.Load(username)
	if ok {
		return conn.(wsadapter.Conn)
	}
	return nil
}

// Add adds a WebSocket connection to the cache.
func (c *connectedUserCache) Add(username string, conn wsadapter.Conn) {
	c.connectedUsers.Store(username, conn)
	logger.Infof("%s added to cache", username)
}

// Remove removes a WebSocket connection from the cache.
func (c *connectedUserCache) Remove(username string) {
	c.connectedUsers.Delete(username)
	logger.Infof("%s removed from cache", username)
}
