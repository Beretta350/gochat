package chat

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"

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

// IsOnline checks if a user is online
func (c *ConnectedUsers) IsOnline(userToken string) bool {
	_, ok := c.connections.Load(userToken)
	return ok
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

	// Deliver pending messages first
	s.deliverPendingMessages(userCtx, conn, userToken)

	// Start listening to Redis for messages directed to this user
	go s.listenForMessages(userCtx, conn, userToken)

	// Read messages from WebSocket and publish to Redis
	s.readAndPublishMessages(userCtx, conn, userToken)
}

// deliverPendingMessages delivers any messages that were sent while user was offline
func (s *Service) deliverPendingMessages(ctx context.Context, conn *websocket.Conn, userToken string) {
	messages, err := redisclient.GetPendingMessages(ctx, userToken)
	if err != nil {
		logger.Errorf("Error getting pending messages for %s: %v", userToken, err)
		return
	}

	if len(messages) == 0 {
		return
	}

	logger.Infof("Delivering %d pending messages to %s", len(messages), userToken)

	for _, msgJSON := range messages {
		// Parse to add received_at
		var chatMsg model.ChatMessage
		if err := json.Unmarshal([]byte(msgJSON), &chatMsg); err != nil {
			logger.Errorf("Error parsing pending message: %v", err)
			continue
		}

		// Mark as received
		chatMsg.MarkReceived()

		// Re-serialize
		msgWithTimestamp, err := json.Marshal(chatMsg)
		if err != nil {
			logger.Errorf("Error marshaling message: %v", err)
			continue
		}

		// Send to WebSocket
		if err := conn.WriteMessage(websocket.TextMessage, msgWithTimestamp); err != nil {
			logger.Errorf("Error sending pending message to %s: %v", userToken, err)
			return
		}
	}
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

			// Parse message to add received_at
			var chatMsg model.ChatMessage
			if err := json.Unmarshal([]byte(msg.Payload), &chatMsg); err != nil {
				logger.Errorf("Error parsing message: %v", err)
				continue
			}

			// Mark as received
			chatMsg.MarkReceived()

			// Re-serialize with received_at
			msgWithTimestamp, err := json.Marshal(chatMsg)
			if err != nil {
				logger.Errorf("Error marshaling message: %v", err)
				continue
			}

			// Send message to WebSocket
			if err := conn.WriteMessage(websocket.TextMessage, msgWithTimestamp); err != nil {
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

			// Set ID, sender and sent_at timestamp
			chatMsg.ID = uuid.New().String()
			chatMsg.Sender = userToken
			chatMsg.SentAt = time.Now().UnixMilli()
			if chatMsg.Type == "" {
				chatMsg.Type = "text"
			}

			// Process and route the message
			s.processMessage(ctx, &chatMsg)
		}
	}
}

// processMessage handles message routing and persistence
func (s *Service) processMessage(ctx context.Context, msg *model.ChatMessage) {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		logger.Errorf("Error marshaling message: %v", err)
		return
	}

	// 1. Add to Redis Stream (for persistence worker)
	streamID, err := redisclient.AddToStream(ctx, map[string]interface{}{
		"data": string(msgJSON),
	})
	if err != nil {
		logger.Errorf("Error adding to stream: %v", err)
		// Continue anyway - real-time delivery is more important
	} else {
		logger.Infof("Message added to stream with ID: %s", streamID)
	}

	// 2. Check if recipient is online
	if s.users.IsOnline(msg.Recipient) {
		// Online: publish to Pub/Sub for real-time delivery
		recipientChannel := "user:" + msg.Recipient
		if err := redisclient.Publish(ctx, recipientChannel, msgJSON); err != nil {
			logger.Errorf("Error publishing message to %s: %v", msg.Recipient, err)
		} else {
			logger.Infof("Message from %s to %s published (online)", msg.Sender, msg.Recipient)
		}
	} else {
		// Offline: add to pending queue
		if err := redisclient.AddToPending(ctx, msg.Recipient, string(msgJSON)); err != nil {
			logger.Errorf("Error adding to pending for %s: %v", msg.Recipient, err)
		} else {
			logger.Infof("Message from %s to %s queued (offline)", msg.Sender, msg.Recipient)
		}
	}
}
