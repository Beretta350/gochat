package chat

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/gofiber/contrib/websocket"

	"github.com/Beretta350/gochat/internal/app/model"
	"github.com/Beretta350/gochat/pkg/logger"
	"github.com/Beretta350/gochat/pkg/redisclient"
)

// ConnectedUsers stores WebSocket connections by user token
type ConnectedUsers struct {
	connections sync.Map
}

var (
	connectedUsers *ConnectedUsers
	usersOnce      sync.Once
)

// GetConnectedUsers returns the singleton instance
func GetConnectedUsers() *ConnectedUsers {
	usersOnce.Do(func() {
		connectedUsers = &ConnectedUsers{}
	})
	return connectedUsers
}

// Add adds a connection
func (c *ConnectedUsers) Add(userToken string, conn *websocket.Conn) {
	c.connections.Store(userToken, conn)
	logger.Infof("User %s connected", userToken)
}

// Remove removes a connection
func (c *ConnectedUsers) Remove(userToken string) {
	c.connections.Delete(userToken)
	logger.Infof("User %s disconnected", userToken)
}

// Get gets a connection
func (c *ConnectedUsers) Get(userToken string) *websocket.Conn {
	if conn, ok := c.connections.Load(userToken); ok {
		return conn.(*websocket.Conn)
	}
	return nil
}

// Service handles chat operations
type Service struct {
	users *ConnectedUsers
}

// NewService creates a new chat service
func NewService() *Service {
	return &Service{
		users: GetConnectedUsers(),
	}
}

// HandleConnection handles a WebSocket connection
func (s *Service) HandleConnection(ctx context.Context, conn *websocket.Conn, userToken string) {
	// Add user to connected users
	s.users.Add(userToken, conn)
	defer s.users.Remove(userToken)

	// Create context for this user session
	userCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Start listening to Redis for messages directed to this user
	go s.listenForMessages(userCtx, conn, userToken)

	// Read messages from WebSocket and publish to Redis
	s.readAndPublishMessages(userCtx, conn, userToken)
}

// listenForMessages subscribes to Redis channel for this user
func (s *Service) listenForMessages(ctx context.Context, conn *websocket.Conn, userToken string) {
	channel := "user:" + userToken
	pubsub := redisclient.Subscribe(ctx, channel)
	defer func() {
		_ = pubsub.Close()
	}()

	logger.Infof("User %s subscribed to channel %s", userToken, channel)

	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			logger.Infof("Stopping listener for %s", userToken)
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}

			// Send message to WebSocket
			if err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
				logger.Errorf("Error writing to WebSocket for %s: %v", userToken, err)
				return
			}
		}
	}
}

// readAndPublishMessages reads from WebSocket and publishes to Redis
func (s *Service) readAndPublishMessages(ctx context.Context, conn *websocket.Conn, userToken string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Read message from WebSocket
			_, msgBytes, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					logger.Infof("User %s WebSocket closed", userToken)
				} else {
					logger.Errorf("Error reading from WebSocket for %s: %v", userToken, err)
				}
				return
			}

			// Parse message
			var chatMsg model.ChatMessage
			if err := json.Unmarshal(msgBytes, &chatMsg); err != nil {
				logger.Errorf("Error parsing message from %s: %v", userToken, err)
				continue
			}

			// Set sender
			chatMsg.Sender = userToken

			// Publish to recipient's channel
			s.publishMessage(ctx, &chatMsg)
		}
	}
}

// publishMessage publishes a message to the recipient's Redis channel
func (s *Service) publishMessage(ctx context.Context, msg *model.ChatMessage) {
	recipientChannel := "user:" + msg.Recipient

	msgJSON, err := json.Marshal(msg)
	if err != nil {
		logger.Errorf("Error marshaling message: %v", err)
		return
	}

	if err := redisclient.Publish(ctx, recipientChannel, msgJSON); err != nil {
		logger.Errorf("Error publishing message to %s: %v", msg.Recipient, err)
		return
	}

	logger.Infof("Message from %s to %s published", msg.Sender, msg.Recipient)
}
