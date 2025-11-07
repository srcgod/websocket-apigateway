package services

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/srcgod/apigateway/internal/websocket/models"
	"github.com/srcgod/apigateway/internal/websocket/types"
)

type MessageServiceHandler struct {
	logger *logrus.Logger
}

func NewMessageServiceHandler(logger *logrus.Logger) *MessageServiceHandler {
	return &MessageServiceHandler{
		logger: logger,
	}
}

func (h *MessageServiceHandler) Handle(req models.GRPCRequest, client types.ClientInterface) models.GRPCResponse {
	h.logger.WithFields(logrus.Fields{
		"client_id":  client.ID(),
		"user_id":    client.UserID(),
		"request_id": req.ID,
		"action":     req.Action,
	}).Info("💬 Message-service gRPC call")

	fmt.Println("Message data:", req.Data)

	// TODO: grpc call
	return models.GRPCResponse{
		ID:      req.ID,
		Success: true,
		Data:    req.Data,
	}
}
