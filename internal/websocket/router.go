package websock

import (
	"github.com/sirupsen/logrus"
	"github.com/srcgod/apigateway/internal/websocket/models"
	"github.com/srcgod/apigateway/internal/websocket/services"
	"github.com/srcgod/apigateway/internal/websocket/types"
)

// serviceRouter маршрутизирует запросы к соответствующим сервисам
type ServiceRouter struct {
	services map[string]types.ServiceHandler
	logger   *logrus.Logger
}

func NewServiceRouter(logger *logrus.Logger) *ServiceRouter {
	router := &ServiceRouter{
		services: make(map[string]types.ServiceHandler),
		logger:   logger,
	}

	router.RegisterService("message-service", services.NewMessageServiceHandler(logger))
	router.RegisterService("auth-service", services.NewAuthServicesHandler(logger))

	return router
}

func (r *ServiceRouter) RegisterService(name string, handler types.ServiceHandler) {
	r.services[name] = handler
	r.logger.WithField("service", name).Info("📝 Service registered")
}

func (r *ServiceRouter) RouteToService(req models.GRPCRequest, client types.ClientInterface) models.GRPCResponse {
	handler, exists := r.services[req.Service]
	if !exists {
		r.logger.WithFields(logrus.Fields{
			"service":    req.Service,
			"request_id": req.ID,
		}).Warn("⚠️ Service not found")

		return models.GRPCResponse{
			ID:      req.ID,
			Success: false,
			Error: &models.ErrorResponse{
				Code:    "unknown_service",
				Message: "Service not found: " + req.Service,
			},
		}
	}

	return handler.Handle(req, client)
}
