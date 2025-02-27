package handler

import (
	"errors"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/Beretta350/gochat/internal/app/cache"
	"github.com/Beretta350/gochat/internal/app/service"
	"github.com/Beretta350/gochat/pkg/logger"
)

type WebsocketHandler interface {
	HandleConnection(w http.ResponseWriter, r *http.Request)
}

type websocketHandler struct {
	service  service.WebsocketService
	upgrader *websocket.Upgrader
}

func NewWebsocketHandler(s service.WebsocketService, u *websocket.Upgrader) WebsocketHandler {
	return &websocketHandler{service: s, upgrader: u}
}

func (h *websocketHandler) HandleConnection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Upgrade initial GET request to a WebSocket
	logger.Info("Connecting from: ", r.RemoteAddr)
	ws, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		respond(w, err)
		return
	}

	//TODO: Develop a retry method for closing connection
	defer func() { _ = ws.Close() }()

	userToken := r.URL.Query().Get("token")
	if userToken == "" {
		// missing user token error
		respond(w, websocket.HandshakeError{})
		return
	}

	client, err := h.service.SetupSession(ctx, ws, userToken)
	if err != nil {
		respond(w, err)
		return
	}

	defer func() {
		client.CloseConnection()
		cache.GetConnectedUserCache().Remove(userToken)
		logger.Infof("%s fully disconnected", userToken)
	}()

	h.service.HandleSession(ctx, ws, client)
}

// handleResponse handles both logging the error and sending the appropriate HTTP response.
func respond(w http.ResponseWriter, err error) {
	// TODO: Fix this when create custom errors
	if err != nil {
		// Log the error
		logger.Error("Error occurred: ", err)

		// Determine whether the error is internal or a bad request
		var handshakeError websocket.HandshakeError
		switch {
		case errors.As(err, &handshakeError): // Handle internal server error
			http.Error(w, "Bad Request", http.StatusBadRequest)
		default: // Handle bad request for all other errors
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
