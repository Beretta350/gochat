package handler

import (
	"context"
	"github.com/Beretta350/gochat/internal/app/cache"
	"github.com/Beretta350/gochat/internal/app/service"
	"github.com/Beretta350/gochat/pkg/logger"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebsocketHandler interface {
	HandleConnection(w http.ResponseWriter, r *http.Request)
}

type websocketHandler struct {
	service service.WebsocketService
}

func NewWebsocketHandler(s service.WebsocketService) WebsocketHandler {
	return &websocketHandler{service: s}
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

	userToken := r.URL.Query().Get("token")
	if userToken == "" {
		logger.Fatal("Missing user token")
	}

	//TODO: Develop a retry method for closing connection
	defer func() {
		cancel()
		_ = ws.Close()
		usersCache.Remove(userToken)
	}()

	client, err := h.service.SetupSession(ctx, ws, userToken)
	if err != nil {
		logger.Fatal(err)
	}

	h.service.HandleSession(ctx, ws, client)
}
