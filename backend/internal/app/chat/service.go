package chat

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"

	"github.com/Beretta350/gochat/internal/app/repository"
	"github.com/Beretta350/gochat/pkg/logger"
	"github.com/Beretta350/gochat/pkg/redisclient"
)

// ConnectedUsers stores WebSocket connections by user ID
type ConnectedUsers struct {
	connections sync.Map
}

// NewConnectedUsers creates a new ConnectedUsers instance
func NewConnectedUsers() *ConnectedUsers {
	return &ConnectedUsers{}
}

// Add adds a connection
func (c *ConnectedUsers) Add(userID string, conn *websocket.Conn) {
	c.connections.Store(userID, conn)
	logger.Infof("User %s connected", userID)
}

// Remove removes a connection
func (c *ConnectedUsers) Remove(userID string) {
	c.connections.Delete(userID)
	logger.Infof("User %s disconnected", userID)
}

// Get gets a connection
func (c *ConnectedUsers) Get(userID string) *websocket.Conn {
	if conn, ok := c.connections.Load(userID); ok {
		return conn.(*websocket.Conn)
	}
	return nil
}

// IsOnline checks if a user is online
func (c *ConnectedUsers) IsOnline(userID string) bool {
	_, ok := c.connections.Load(userID)
	return ok
}

// WebSocketMessage represents a message received via WebSocket
type WebSocketMessage struct {
	ConversationID string `json:"conversation_id"`
	Content        string `json:"content"`
	Type           string `json:"type,omitempty"`
}

// OutgoingMessage represents a message sent to WebSocket clients
type OutgoingMessage struct {
	ID             string `json:"id"`
	ConversationID string `json:"conversation_id"`
	SenderID       string `json:"sender_id"`
	SenderUsername string `json:"sender_username,omitempty"`
	Content        string `json:"content"`
	Type           string `json:"type"`
	SentAt         int64  `json:"sent_at"`
}

// PresenceEvent represents an online/offline status change
type PresenceEvent struct {
	Type     string `json:"type"`     // "presence"
	UserID   string `json:"user_id"`
	Username string `json:"username,omitempty"`
	Status   string `json:"status"`   // "online" or "offline"
}

// Service handles chat operations
type Service struct {
	redis    *redisclient.Client
	convRepo repository.ConversationRepository
	userRepo repository.UserRepository
	users    *ConnectedUsers
}

// NewService creates a new chat service (Fx provider)
func NewService(redis *redisclient.Client, convRepo repository.ConversationRepository, userRepo repository.UserRepository) *Service {
	logger.Info("Chat service initialized")
	return &Service{
		redis:    redis,
		convRepo: convRepo,
		userRepo: userRepo,
		users:    NewConnectedUsers(),
	}
}

// HandleConnection handles a WebSocket connection
func (s *Service) HandleConnection(ctx context.Context, conn *websocket.Conn, userID string) {
	s.users.Add(userID, conn)

	userCtx, cancel := context.WithCancel(ctx)

	// Mark user as online in Redis and broadcast presence
	s.handleUserOnline(userCtx, userID, conn)

	defer func() {
		cancel()
		s.users.Remove(userID)
		// Mark user as offline in Redis and broadcast presence
		s.handleUserOffline(context.Background(), userID)
	}()

	// Deliver pending messages first
	s.deliverPendingMessages(userCtx, conn, userID)

	// Start listening to Redis for messages
	go s.listenForMessages(userCtx, conn, userID)

	// Read messages from WebSocket
	s.readAndPublishMessages(userCtx, conn, userID)
}

func (s *Service) deliverPendingMessages(ctx context.Context, conn *websocket.Conn, userID string) {
	messages, err := s.redis.GetPendingMessages(ctx, userID)
	if err != nil {
		logger.Errorf("Error getting pending messages for %s: %v", userID, err)
		return
	}

	if len(messages) == 0 {
		return
	}

	logger.Infof("Delivering %d pending messages to %s", len(messages), userID)

	for _, msgJSON := range messages {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(msgJSON)); err != nil {
			logger.Errorf("Error sending pending message to %s: %v", userID, err)
			return
		}
	}
}

func (s *Service) listenForMessages(ctx context.Context, conn *websocket.Conn, userID string) {
	channel := "user:" + userID
	pubsub := s.redis.Subscribe(ctx, channel)
	defer func() {
		_ = pubsub.Close()
	}()

	logger.Infof("User %s subscribed to channel %s", userID, channel)

	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			logger.Infof("Stopping listener for %s", userID)
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}

			// Forward message to WebSocket
			if err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
				logger.Errorf("Error writing to WebSocket for %s: %v", userID, err)
				return
			}
		}
	}
}

func (s *Service) readAndPublishMessages(ctx context.Context, conn *websocket.Conn, userID string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, msgBytes, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					logger.Infof("User %s WebSocket closed", userID)
				} else {
					logger.Errorf("Error reading from WebSocket for %s: %v", userID, err)
				}
				return
			}

			var wsMsg WebSocketMessage
			if err := json.Unmarshal(msgBytes, &wsMsg); err != nil {
				logger.Errorf("Error parsing message from %s: %v", userID, err)
				s.sendError(conn, "Invalid message format")
				continue
			}

			if wsMsg.ConversationID == "" {
				s.sendError(conn, "conversation_id is required")
				continue
			}

			if wsMsg.Content == "" {
				s.sendError(conn, "content is required")
				continue
			}

			s.processMessage(ctx, conn, userID, &wsMsg)
		}
	}
}

func (s *Service) processMessage(ctx context.Context, conn *websocket.Conn, senderID string, wsMsg *WebSocketMessage) {
	// Get conversation participants
	participants, err := s.convRepo.GetParticipants(ctx, wsMsg.ConversationID)
	if err != nil {
		logger.Errorf("Error getting participants for conversation %s: %v", wsMsg.ConversationID, err)
		s.sendError(conn, "Conversation not found")
		return
	}

	// Verify sender is participant
	isParticipant := false
	var senderUsername string
	for _, p := range participants {
		if p.UserID == senderID {
			isParticipant = true
			if p.User != nil {
				senderUsername = p.User.Username
			}
			break
		}
	}

	if !isParticipant {
		s.sendError(conn, "You are not a participant of this conversation")
		return
	}

	// Create outgoing message
	msgType := wsMsg.Type
	if msgType == "" {
		msgType = "text"
	}

	outMsg := &OutgoingMessage{
		ID:             uuid.New().String(),
		ConversationID: wsMsg.ConversationID,
		SenderID:       senderID,
		SenderUsername: senderUsername,
		Content:        wsMsg.Content,
		Type:           msgType,
		SentAt:         time.Now().UnixMilli(),
	}

	msgJSON, err := json.Marshal(outMsg)
	if err != nil {
		logger.Errorf("Error marshaling message: %v", err)
		return
	}

	// Add to Redis Stream for persistence
	streamData := map[string]interface{}{
		"data": string(msgJSON),
	}
	// Add fields for worker to save to PostgreSQL
	streamData["id"] = outMsg.ID
	streamData["conversation_id"] = outMsg.ConversationID
	streamData["sender_id"] = outMsg.SenderID
	streamData["content"] = outMsg.Content
	streamData["type"] = outMsg.Type
	streamData["sent_at"] = outMsg.SentAt

	if _, err := s.redis.AddToStream(ctx, streamData); err != nil {
		logger.Errorf("Error adding to stream: %v", err)
	}

	// Send to all participants
	for _, p := range participants {
		if p.UserID == senderID {
			continue // Don't send back to sender
		}

		if s.users.IsOnline(p.UserID) {
			// Online: publish to Pub/Sub
			recipientChannel := "user:" + p.UserID
			if err := s.redis.Publish(ctx, recipientChannel, msgJSON); err != nil {
				logger.Errorf("Error publishing to %s: %v", p.UserID, err)
			}
		} else {
			// Offline: add to pending queue
			if err := s.redis.AddToPending(ctx, p.UserID, string(msgJSON)); err != nil {
				logger.Errorf("Error adding to pending for %s: %v", p.UserID, err)
			}
		}
	}

	logger.Infof("Message in conversation %s from %s", wsMsg.ConversationID, senderID)
}

func (s *Service) sendError(conn *websocket.Conn, message string) {
	errMsg := map[string]interface{}{
		"error":   true,
		"message": message,
	}
	msgBytes, _ := json.Marshal(errMsg)
	_ = conn.WriteMessage(websocket.TextMessage, msgBytes)
}

// handleUserOnline marks user as online and broadcasts to their contacts
func (s *Service) handleUserOnline(ctx context.Context, userID string, conn *websocket.Conn) {
	// Mark online in Redis
	if err := s.redis.SetUserOnline(ctx, userID); err != nil {
		logger.Errorf("Error setting user %s online in Redis: %v", userID, err)
	}

	// Get user info for username
	user, err := s.userRepo.GetByID(ctx, userID)
	var username string
	if err == nil && user != nil {
		username = user.Username
	}

	// Get user's conversations to find contacts
	conversations, err := s.convRepo.GetByUserID(ctx, userID)
	if err != nil {
		logger.Errorf("Error getting conversations for user %s: %v", userID, err)
		return
	}

	// Collect all unique contact IDs and send them initial online status
	contactsSet := make(map[string]bool)
	for _, conv := range conversations {
		participants, err := s.convRepo.GetParticipants(ctx, conv.ID)
		if err != nil {
			continue
		}
		for _, p := range participants {
			if p.UserID != userID {
				contactsSet[p.UserID] = true
			}
		}
	}

	// Get list of online contacts to send to the connecting user
	var contactIDs []string
	for contactID := range contactsSet {
		contactIDs = append(contactIDs, contactID)
	}

	// Send initial online users list to the connecting user
	onlineContacts, err := s.redis.GetOnlineUsersFromList(ctx, contactIDs)
	if err == nil && len(onlineContacts) > 0 {
		initialStatus := map[string]interface{}{
			"type":         "presence_list",
			"online_users": onlineContacts,
		}
		if msgBytes, err := json.Marshal(initialStatus); err == nil {
			_ = conn.WriteMessage(websocket.TextMessage, msgBytes)
		}
	}

	// Broadcast online status to all contacts
	presenceEvent := &PresenceEvent{
		Type:     "presence",
		UserID:   userID,
		Username: username,
		Status:   "online",
	}
	s.broadcastPresence(ctx, presenceEvent, contactIDs)
}

// handleUserOffline marks user as offline and broadcasts to their contacts
func (s *Service) handleUserOffline(ctx context.Context, userID string) {
	// Mark offline in Redis
	if err := s.redis.SetUserOffline(ctx, userID); err != nil {
		logger.Errorf("Error setting user %s offline in Redis: %v", userID, err)
	}

	// Get user info for username
	user, err := s.userRepo.GetByID(ctx, userID)
	var username string
	if err == nil && user != nil {
		username = user.Username
	}

	// Get user's conversations to find contacts
	conversations, err := s.convRepo.GetByUserID(ctx, userID)
	if err != nil {
		logger.Errorf("Error getting conversations for user %s: %v", userID, err)
		return
	}

	// Collect all unique contact IDs
	contactsSet := make(map[string]bool)
	for _, conv := range conversations {
		participants, err := s.convRepo.GetParticipants(ctx, conv.ID)
		if err != nil {
			continue
		}
		for _, p := range participants {
			if p.UserID != userID {
				contactsSet[p.UserID] = true
			}
		}
	}

	var contactIDs []string
	for contactID := range contactsSet {
		contactIDs = append(contactIDs, contactID)
	}

	// Broadcast offline status to all contacts
	presenceEvent := &PresenceEvent{
		Type:     "presence",
		UserID:   userID,
		Username: username,
		Status:   "offline",
	}
	s.broadcastPresence(ctx, presenceEvent, contactIDs)
}

// broadcastPresence sends presence event to specified users
func (s *Service) broadcastPresence(ctx context.Context, event *PresenceEvent, userIDs []string) {
	msgBytes, err := json.Marshal(event)
	if err != nil {
		logger.Errorf("Error marshaling presence event: %v", err)
		return
	}

	for _, userID := range userIDs {
		if s.users.IsOnline(userID) {
			// User is online, publish to their channel
			channel := "user:" + userID
			if err := s.redis.Publish(ctx, channel, msgBytes); err != nil {
				logger.Errorf("Error publishing presence to %s: %v", userID, err)
			}
		}
		// Note: We don't queue presence events for offline users
	}
}

// GetOnlineUsers returns the list of online user IDs (for API endpoint)
func (s *Service) GetOnlineUsers(ctx context.Context) ([]string, error) {
	return s.redis.GetOnlineUsers(ctx)
}

// GetOnlineUsersFromList checks which users from the list are online
func (s *Service) GetOnlineUsersFromList(ctx context.Context, userIDs []string) ([]string, error) {
	return s.redis.GetOnlineUsersFromList(ctx, userIDs)
}
