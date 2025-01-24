package handler

import (
	"context"
	"fmt"
	"github.com/Beretta350/gochat/internal/app/cache"
	"github.com/Beretta350/gochat/internal/app/client"
	"github.com/Beretta350/gochat/pkg/logger"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/Beretta350/gochat/internal/app/model"
)

// Mutex for thread-safe access to `users`
var mutex = &sync.Mutex{}

// Channel for incoming messages
//var broadcast = make(chan model.ChatMessage)

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebsocketHandler interface {
	HandleConnection(w http.ResponseWriter, r *http.Request)
}

type websocketHandler struct{}

func NewWebsocketHandler() WebsocketHandler {
	return &websocketHandler{}
}

func (h *websocketHandler) HandleConnection(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	usersCache := cache.GetConnectedUserCache()

	// Upgrade initial GET request to a WebSocket
	logger.Info("Connecting from: ", r.RemoteAddr)
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	//TODO: Develop a retry method for closing connection
	defer func() {
		cancel()
		_ = ws.Close()
	}()

	kafkaClient, userToken, err := h.handleHandshake(ctx, r.URL.Query(), ws)
	if err != nil {
		logger.Error(err)
		return
	}

	//The cache use sync.Map that means already has mutex
	defer usersCache.Remove(userToken)

	go runKafkaConsumer(r.Context(), kafkaClient)
	listenIncomingWSMessages(ws, kafkaClient)
}

func (h *websocketHandler) handleHandshake(ctx context.Context, queryParams url.Values, ws *websocket.Conn) (client.KafkaClient, string, error) {
	usersCache := cache.GetConnectedUserCache()

	// Extract the token from the query parameters
	userToken := queryParams.Get("token")
	if userToken == "" {
		return nil, "", fmt.Errorf("missing user token in WebSocket URL")
	}

	logger.Info("Initial handshake received for: ", userToken)

	// Check if the userToken is already connected
	if conn := usersCache.Get(userToken); conn != nil {
		return nil, "", fmt.Errorf("user %s is already connected", userToken)
	}

	kafkaClient, err := client.NewKafkaClient(ctx, userToken)
	if err != nil {
		return nil, "", fmt.Errorf("kafka client creation error: %w", err)
	}

	// Add the user to the cache (thread-safe)
	usersCache.Add(userToken, ws)
	return kafkaClient, userToken, nil
}

func listenIncomingWSMessages(ws *websocket.Conn, client client.KafkaClient) {
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

func runKafkaConsumer(ctx context.Context, client client.KafkaClient) {
	err := client.ConsumeMessage(ctx, handleChatMessages)
	if err != nil {
		logger.Error("error on kafka consumer: ", err)
	}

	client.CloseConnection()
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
