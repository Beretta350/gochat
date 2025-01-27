package handler

import (
	"github.com/Beretta350/gochat/internal/app/cache"
	"github.com/Beretta350/gochat/internal/app/service"
	"github.com/Beretta350/gochat/pkg/logger"
	"github.com/gorilla/websocket"
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
	ctx := r.Context()

	// Upgrade initial GET request to a WebSocket
	logger.Info("Connecting from: ", r.RemoteAddr)
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(err)
		return
	}
	//TODO: Develop a retry method for closing connection
	defer func() { _ = ws.Close() }()

	userToken := r.URL.Query().Get("token")
	if userToken == "" {
		logger.Error("Missing user token")
		return
	}

	client, err := h.service.SetupSession(ctx, ws, userToken)
	if err != nil {
		logger.Error(err)
		return
	}

	defer func() {
		client.CloseConnection()
		cache.GetConnectedUserCache().Remove(userToken)
		logger.Infof("%s fully disconnected", userToken)
	}()

	h.service.HandleSession(ctx, ws, client)
}
