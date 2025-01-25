package handler

import (
	"context"
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

	userToken := r.URL.Query().Get("token")
	if userToken == "" {
		logger.Error("Missing user token")
		return
	}

	err = h.service.HandleWebsocketSession(ctx, ws, userToken)
	if err != nil {
		logger.Error("WebSocket session error:", err)
	}
}
