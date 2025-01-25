package service

import (
	"context"
	"fmt"
	"github.com/Beretta350/gochat/internal/app/cache"
	"github.com/Beretta350/gochat/internal/app/messaging"
	"github.com/Beretta350/gochat/internal/app/model"
	"github.com/Beretta350/gochat/pkg/logger"
	"github.com/gorilla/websocket"
	"sync"
)

// Mutex for thread-safe access to `users`
var mutex = &sync.Mutex{}

type WebsocketService interface {
	HandleWebsocketSession(ctx context.Context, ws *websocket.Conn, userToken string) error
}

type websocketService struct{}

func NewWebsocketService() WebsocketService {
	return &websocketService{}
}

func (s *websocketService) HandleWebsocketSession(ctx context.Context, ws *websocket.Conn, userToken string) error {
	usersCache := cache.GetConnectedUserCache()
	logger.Info("Initial handshake received for: ", userToken)

	kafkaClient, err := s.authenticateAndInitializeSession(ctx, ws, userToken)
	if err != nil {
		return err
	}

	defer func() {
		kafkaClient.CloseConnection()
		usersCache.Remove(userToken)
	}()

	go s.startChatMessagesHandler(ctx, kafkaClient)
	s.startIncomingMessagesHandler(ws, kafkaClient)

	return nil
}

func (s *websocketService) authenticateAndInitializeSession(ctx context.Context, ws *websocket.Conn, userToken string) (messaging.KafkaClient, error) {
	usersCache := cache.GetConnectedUserCache()

	// Check if the userToken is already connected
	if conn := usersCache.Get(userToken); conn != nil {
		return nil, fmt.Errorf("user %s is already connected", userToken)
	}

	kafkaClient, err := messaging.NewKafkaClient(ctx, userToken)
	if err != nil {
		return nil, fmt.Errorf("kafka client creation error: %w", err)
	}

	// Add the user to the cache (thread-safe)
	usersCache.Add(userToken, ws)
	return kafkaClient, nil
}

func (s *websocketService) startChatMessagesHandler(ctx context.Context, client messaging.KafkaClient) {
	err := client.ConsumeMessage(ctx, handleChatMessages)
	if err != nil {
		logger.Error("Error in Kafka consumer:", err)
	}
}

func (s *websocketService) startIncomingMessagesHandler(ws *websocket.Conn, client messaging.KafkaClient) {
	for {
		var msg model.ChatMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			logger.Infof("Error decoding websocket message: %v", err)
			break
		}

		// Forward message to kafka topic
		pErr := client.ProduceMessage(msg)
		if pErr != nil {
			logger.Errorf("Error sending message to %s: %v", msg.Recipient, pErr)
		}
	}
}

func handleChatMessages(message model.ChatMessage) {
	userCache := cache.GetConnectedUserCache()

	logger.Info("broadcast received:", message)

	// Find the recipient's WebSocket connection
	mutex.Lock()
	recipientConn := userCache.Get(message.Recipient)
	mutex.Unlock()

	if recipientConn != nil {
		// Send the message to the recipient websocket connection
		err := recipientConn.WriteJSON(message)
		if err != nil {
			logger.Infof("Error sending message to %s: %v", message.Recipient, err)
			_ = recipientConn.Close()
			mutex.Lock()
			userCache.Remove(message.Recipient)
			mutex.Unlock()
		}
	} else {
		logger.Infof("Recipient %s not found", message.Recipient)
	}
}
