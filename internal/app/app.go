package app

import (
	"github.com/Beretta350/gochat/config"
	"github.com/Beretta350/gochat/pkg/logger"
	"net/http"

	"github.com/Beretta350/gochat/internal/app/handler"
)

func Run() {
	serverConfig := config.GetServerConfig()
	websocketHandler := handler.NewWebsocketHandler()

	// Serve static files
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	// Configure WebSocket route
	http.HandleFunc("/ws", websocketHandler.HandleConnection)

	// Start the server
	logger.Infof("Http server started on port %v", serverConfig.Port)
	err := http.ListenAndServe(":"+serverConfig.Port, nil)
	if err != nil {
		logger.Fatal("Listen and serve:", err)
	}
}
