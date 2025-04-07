package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/Beretta350/gochat/internal/app/adapters/wsadapter"
	"github.com/Beretta350/gochat/internal/app/cache"
	"github.com/Beretta350/gochat/internal/app/kafkaclient"
	"github.com/Beretta350/gochat/internal/app/model"
	"github.com/Beretta350/gochat/pkg/logger"
)

// Mutex for thread-safe access to `users`
var mutex = &sync.Mutex{}

type WebsocketService interface {
	HandleSession(ctx context.Context, ws wsadapter.Conn, client kafkaclient.KafkaClient)
	SetupSession(ctx context.Context, ws wsadapter.Conn, userToken string) (kafkaclient.KafkaClient, error)
}

type websocketService struct{}

func NewWebsocketService() WebsocketService {
	return &websocketService{}
}

func (s *websocketService) HandleSession(ctx context.Context, ws wsadapter.Conn, client kafkaclient.KafkaClient) {
	userCtx, cancelUserCtx := context.WithCancel(ctx)
	wg := &sync.WaitGroup{}

	go s.startChatMessageReciever(userCtx, client, wg)
	wg.Add(1)
	go s.startChatMessageSender(cancelUserCtx, ws, client, wg)
	wg.Add(1)

	wg.Wait()
}

func (s *websocketService) SetupSession(ctx context.Context, ws wsadapter.Conn, userToken string) (kafkaclient.KafkaClient, error) {
	usersCache := cache.GetConnectedUserCache()
	logger.Info("Initial handshake received for: ", userToken)

	// Check if the userToken is already connected
	if conn := usersCache.Get(userToken); conn != nil {
		return nil, fmt.Errorf("user %s is already connected", userToken)
	}

	kafkaClient, err := kafkaclient.NewKafkaClient(ctx, userToken)
	if err != nil {
		return nil, fmt.Errorf("kafka client creation error: %w", err)
	}

	// Add the user to the cache (thread-safe)
	usersCache.Add(userToken, ws)
	return kafkaClient, nil
}

func (s *websocketService) startChatMessageReciever(ctx context.Context, client kafkaclient.KafkaClient, wg *sync.WaitGroup) {
	defer wg.Done()

	err := client.ConsumeMessage(ctx, handleChatMessages)
	if err != nil {
		logger.Error("Error in Kafka consumer:", err)
	}
}

func (s *websocketService) startChatMessageSender(cancelCtx context.CancelFunc, ws wsadapter.Conn, client kafkaclient.KafkaClient, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		var msg model.ChatMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				logger.Info("Starting user disconnection...")
			} else {
				logger.Errorf("Error decoding websocket message: %v", err)
			}
			cancelCtx()
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
