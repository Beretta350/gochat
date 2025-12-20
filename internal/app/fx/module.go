package fx

import (
	"go.uber.org/fx"

	"github.com/Beretta350/gochat/internal/app/chat"
	"github.com/Beretta350/gochat/internal/app/handler"
	"github.com/Beretta350/gochat/internal/app/repository"
	"github.com/Beretta350/gochat/internal/app/worker"
	"github.com/Beretta350/gochat/internal/config"
	"github.com/Beretta350/gochat/pkg/redisclient"
)

// Module provides all application dependencies
var Module = fx.Options(
	// Core
	fx.Provide(config.NewConfig),
	fx.Provide(redisclient.NewClient),

	// Repository
	fx.Provide(repository.NewMessageRepository),

	// Services
	fx.Provide(chat.NewService),

	// Workers
	fx.Provide(worker.NewMessageWorker),

	// Handlers
	fx.Provide(handler.NewHealthHandler),
	fx.Provide(handler.NewWebSocketHandler),
)
