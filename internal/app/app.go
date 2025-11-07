package app

import (
	"github.com/sirupsen/logrus"
	"github.com/srcgod/apigateway/internal/config"
	"github.com/srcgod/apigateway/internal/handler"
	"github.com/srcgod/apigateway/internal/websocket/core"
)

type App struct {
	logger  *logrus.Logger
	handler *handler.GatewayHandler
	hub     *core.Hub
}

func New(logger *logrus.Logger) *App {
	hub := core.NewHub(logger)
	go hub.HubRun()

	gatewayHandler := handler.New(logger, hub)

	return &App{
		logger:  logger,
		handler: gatewayHandler,
		hub:     hub,
	}
}

func (a *App) SetupRoutes() *handler.Router {
	router := handler.NewRouter()
	a.handler.SetupRoutes(router)
	return router
}

func (a *App) GetCORSConfig() config.CORSConfig {
	return config.NewCORSConfig()
}

func (a *App) Shutdown() {
	a.logger.Info("Shutting down application...")
}
