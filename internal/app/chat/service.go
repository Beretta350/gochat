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

// NewConnectedUsers creates a new ConnectedUsers instance
func NewConnectedUsers() *ConnectedUsers {
	return &ConnectedUsers{}
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
	redis *redisclient.Client
	users *ConnectedUsers
}

// NewService creates a new chat service (Fx provider)
func NewService(redis *redisclient.Client) *Service {
	logger.Info("Chat service initialized")
	return &Service{
		redis: redis,
		users: NewConnectedUsers(),
	}
}

// HandleConnection handles a WebSocket connection
func (s *Service) HandleConnection(ctx context.Context, conn *websocket.Conn, userToken string) {
	s.users.Add(userToken, conn)
	defer s.users.Remove(userToken)

	userCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Deliver pending messages first
	s.deliverPendingMessages(userCtx, conn, userToken)

	// Start listening to Redis for messages
	go s.listenForMessages(userCtx, conn, userToken)

	// Read messages from WebSocket
	s.readAndPublishMessages(userCtx, conn, userToken)
}

func (s *Service) deliverPendingMessages(ctx context.Context, conn *websocket.Conn, userToken string) {
	messages, err := s.redis.GetPendingMessages(ctx, userToken)
	if err != nil {
		logger.Errorf("Error getting pending messages for %s: %v", userToken, err)
		return
	}

	if len(messages) == 0 {
		return
	}

	logger.Infof("Delivering %d pending messages to %s", len(messages), userToken)

	for _, msgJSON := range messages {
		var chatMsg model.ChatMessage
		if err := json.Unmarshal([]byte(msgJSON), &chatMsg); err != nil {
			logger.Errorf("Error parsing pending message: %v", err)
			continue
		}

		chatMsg.MarkReceived()

		msgWithTimestamp, err := json.Marshal(chatMsg)
		if err != nil {
			logger.Errorf("Error marshaling message: %v", err)
			continue
		}

		if err := conn.WriteMessage(websocket.TextMessage, msgWithTimestamp); err != nil {
			logger.Errorf("Error sending pending message to %s: %v", userToken, err)
			return
		}
	}
}

func (s *Service) listenForMessages(ctx context.Context, conn *websocket.Conn, userToken string) {
	channel := "user:" + userToken
	pubsub := s.redis.Subscribe(ctx, channel)
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

			var chatMsg model.ChatMessage
			if err := json.Unmarshal([]byte(msg.Payload), &chatMsg); err != nil {
				logger.Errorf("Error parsing message: %v", err)
				continue
			}

			chatMsg.MarkReceived()

			msgWithTimestamp, err := json.Marshal(chatMsg)
			if err != nil {
				logger.Errorf("Error marshaling message: %v", err)
				continue
			}

			if err := conn.WriteMessage(websocket.TextMessage, msgWithTimestamp); err != nil {
				logger.Errorf("Error writing to WebSocket for %s: %v", userToken, err)
				return
			}
		}
	}
}

func (s *Service) readAndPublishMessages(ctx context.Context, conn *websocket.Conn, userToken string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, msgBytes, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					logger.Infof("User %s WebSocket closed", userToken)
				} else {
					logger.Errorf("Error reading from WebSocket for %s: %v", userToken, err)
				}
				return
			}

			var chatMsg model.ChatMessage
			if err := json.Unmarshal(msgBytes, &chatMsg); err != nil {
				logger.Errorf("Error parsing message from %s: %v", userToken, err)
				continue
			}

			chatMsg.ID = uuid.New().String()
			chatMsg.Sender = userToken
			chatMsg.SentAt = time.Now().UnixMilli()
			if chatMsg.Type == "" {
				chatMsg.Type = "text"
			}

			s.processMessage(ctx, &chatMsg)
		}
	}
}

func (s *Service) processMessage(ctx context.Context, msg *model.ChatMessage) {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		logger.Errorf("Error marshaling message: %v", err)
		return
	}

	// Add to Redis Stream
	streamID, err := s.redis.AddToStream(ctx, map[string]interface{}{
		"data": string(msgJSON),
	})
	if err != nil {
		logger.Errorf("Error adding to stream: %v", err)
	} else {
		logger.Infof("Message added to stream: %s", streamID)
	}

	// Route message
	if s.users.IsOnline(msg.Recipient) {
		recipientChannel := "user:" + msg.Recipient
		if err := s.redis.Publish(ctx, recipientChannel, msgJSON); err != nil {
			logger.Errorf("Error publishing to %s: %v", msg.Recipient, err)
		} else {
			logger.Infof("Message %s -> %s (online)", msg.Sender, msg.Recipient)
		}
	} else {
		if err := s.redis.AddToPending(ctx, msg.Recipient, string(msgJSON)); err != nil {
			logger.Errorf("Error adding to pending for %s: %v", msg.Recipient, err)
		} else {
			logger.Infof("Message %s -> %s (offline)", msg.Sender, msg.Recipient)
		}
	}
}
