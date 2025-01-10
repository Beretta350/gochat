package handler

import (
	"fmt"
	"github.com/Beretta350/gochat/internal/app/cache"
	"github.com/Beretta350/gochat/pkg/logger"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/Beretta350/gochat/internal/app/model"
)

// Mutex for thread-safe access to `users`
var mutex = &sync.Mutex{}

// Channel for incoming messages
var broadcast = make(chan model.ChatMessage)

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebsocketHandler interface {
	HandleConnection(w http.ResponseWriter, r *http.Request)
	HandleChatMessages()
}

type websocketHandler struct {
	//service service.websocketService
}

func NewWebsocketHandler() WebsocketHandler {
	return &websocketHandler{}
}

func (h *websocketHandler) HandleConnection(w http.ResponseWriter, r *http.Request) {
	usersCache := cache.GetConnectedUserCache()

	// Upgrade initial GET request to a WebSocket
	logger.Info("Connecting from:", r.RemoteAddr)
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	//TODO: Develop a retry method for closing connection
	defer func() { _ = ws.Close() }()

	username, err := h.handleHandshakeMessage(ws)
	if err != nil {
		logger.Info(err)
		return
	}

	//The cache use sync.Map that means already has mutex
	defer usersCache.Remove(username)

	listenIncomingMessages(ws)
}

func (h *websocketHandler) HandleChatMessages() {
	userCache := cache.GetConnectedUserCache()

	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		logger.Info("broadcast received:", msg)

		// Find the recipient's WebSocket connection
		mutex.Lock()
		recipientConn, ok := userCache.Get(msg.Recipient)
		mutex.Unlock()

		if ok {
			// Send message to the recipient
			err := recipientConn.WriteJSON(msg)
			if err != nil {
				logger.Infof("Error sending message to %s: %v", msg.Recipient, err)
				_ = recipientConn.Close()
				mutex.Lock()
				userCache.Remove(msg.Recipient)
				mutex.Unlock()
			}
		} else {
			logger.Infof("Recipient %s not found", msg.Recipient)
		}
	}
}

func listenIncomingMessages(ws *websocket.Conn) {
	for {
		var msg model.ChatMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			logger.Infof("Error decoding message: %v", err)
			break
		}

		// Forward message to the channel
		//TODO: Should Forward message to kafka
		broadcast <- msg
	}
}

func (h *websocketHandler) handleHandshakeMessage(ws *websocket.Conn) (string, error) {
	usersCache := cache.GetConnectedUserCache()

	// Register client connection
	var handshakeMsg model.HandshakeMessage
	err := ws.ReadJSON(&handshakeMsg)
	if err != nil {
		return "", fmt.Errorf("error during initial handshake: %v", err)
	}

	handshakeMsg.Created = time.Now()
	logger.Info("initial handshake received:", handshakeMsg)

	//The cache use sync.Map that means already has mutex
	usersCache.Add(handshakeMsg.Username, ws)
	return handshakeMsg.Username, nil
}
