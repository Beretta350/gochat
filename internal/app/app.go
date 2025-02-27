package app

import (
	"net/http"

	"github.com/Beretta350/gochat/internal/app/handler"
	"github.com/Beretta350/gochat/internal/app/service"
	"github.com/Beretta350/gochat/internal/config"
	"github.com/Beretta350/gochat/pkg/logger"
)

func Run() {
	serverConfig := config.GetServerConfig()
	upgrader := config.GetWebsocketUpgrader()

	websocketService := service.NewWebsocketService()
	websocketHandler := handler.NewWebsocketHandler(websocketService, upgrader)

	// Configure WebSocket route
	http.HandleFunc("/ws", websocketHandler.HandleConnection)

	// Start the server
	logger.Infof("Http server started on port %v", serverConfig.Port)
	err := http.ListenAndServe(":"+serverConfig.Port, nil)
	if err != nil {
		logger.Fatal("Listen and serve:", err)
	}
}
