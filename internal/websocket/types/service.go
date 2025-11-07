package types

import "github.com/srcgod/apigateway/internal/websocket/models"

// ServiceHandler интерфейс для обработчиков сервисов
type ServiceHandler interface {
	Handle(req models.GRPCRequest, client ClientInterface) models.GRPCResponse
}
