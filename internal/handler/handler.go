package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/srcgod/apigateway/internal/websocket/core"
	websock "github.com/srcgod/apigateway/internal/websocket"
)

// GatewayHandler управляет всеми HTTP handlers приложения
type GatewayHandler struct {
	logger           *logrus.Logger
	websocketHandler *websock.WSocketHandler
}

// New создает новый GatewayHandler
func New(logger *logrus.Logger, hub *core.Hub) *GatewayHandler {
	return &GatewayHandler{
		logger:           logger,
		websocketHandler: websock.New(logger, hub),
	}
}

// SetupRoutes настраивает все маршруты приложения
func (g *GatewayHandler) SetupRoutes(router *Router) {
	// Настраиваем WebSocket маршруты
	g.websocketHandler.SetUpRoutes(router.Engine())
}

// Router обертка над gin.Engine для лучшей инкапсуляции
type Router struct {
	engine *gin.Engine
}

// NewRouter создает новый роутер
func NewRouter() *Router {
	return &Router{
		engine: gin.Default(),
	}
}

// Engine возвращает gin.Engine
func (r *Router) Engine() *gin.Engine {
	return r.engine
}

// SetUpRoutes устаревший метод, используйте SetupRoutes
// Deprecated: используйте SetupRoutes
func (g *GatewayHandler) SetUpRoutes(r *gin.Engine) {
	router := &Router{engine: r}
	g.SetupRoutes(router)
}
