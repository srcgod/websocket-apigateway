package websock

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/srcgod/apigateway/internal/utils"
	"github.com/srcgod/apigateway/internal/websocket/core"
	"github.com/srcgod/apigateway/internal/websocket/models"
)

// MessageProcessor обрабатывает сообщения от клиентов
type MessageProcessor struct {
	router *ServiceRouter
	logger *logrus.Logger
}

func NewMessageProcessor(router *ServiceRouter, logger *logrus.Logger) *MessageProcessor {
	return &MessageProcessor{
		router: router,
		logger: logger,
	}
}

func (p *MessageProcessor) ProcessClientMessages(client *core.Client) {
	p.logger.WithFields(logrus.Fields{
		"client_id": client.ID(),
		"user_id":   client.UserID(),
	}).Info("🔄 Client message handler started")

	for message := range client.Receive {
		p.logger.WithFields(logrus.Fields{
			"client_id":    client.ID(),
			"user_id":      client.UserID(),
			"message_type": message.Type,
			"data_length":  len(message.Data),
			"timestamp":    message.Timestamp.Format(time.RFC3339),
		}).Debug("📨 Processing message from client")

		p.processClientMessage(message)
	}

	p.logger.WithFields(logrus.Fields{
		"client_id": client.ID(),
		"user_id":   client.UserID(),
	}).Info("🛑 Client message handler stopped")
}

func (p *MessageProcessor) processClientMessage(message core.ClientMessage) {
	switch message.Type {
	case websocket.TextMessage:
		p.handleTextMessage(message)
	default:
		p.logger.WithFields(logrus.Fields{
			"client_id":    message.Client.ID(),
			"message_type": message.Type,
		}).Warn("⚠️ Unsupported message type")
	}
}

// handleTextMessage обрабатывает текстовые сообщения (JSON)
func (p *MessageProcessor) handleTextMessage(message core.ClientMessage) {
	var baseMsg struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(message.Data, &baseMsg); err != nil {
		p.logger.WithFields(logrus.Fields{
			"client_id": message.Client.ID(),
			"error":     err.Error(),
		}).Warn("⚠️ Failed to parse JSON message")
		// TODO: отправить ошибку клиенту
		return
	}

	switch baseMsg.Type {
	case "grpc_request":
		p.handleGRPCRequest(message)
	case "ping_request":
		p.handlePingRequest(message)
	default:
		p.logger.WithFields(logrus.Fields{
			"client_id":    message.Client.ID(),
			"message_type": baseMsg.Type,
		}).Warn("⚠️ Unknown message type")
	}
}

func (p *MessageProcessor) sendError(client *core.Client, code string, message string) {
	// TODO: реализовать отправку ошибки
	p.logger.WithFields(logrus.Fields{
		"client_id":  client.ID(),
		"error_code": code,
		"error_msg":  message,
	}).Warn("⚠️ Error occurred")
}

func (p *MessageProcessor) handlePingRequest(message core.ClientMessage) {
	var pingReq models.PingMessage

	if err := json.Unmarshal(message.Data, &pingReq); err != nil {
		p.sendError(message.Client, "invalid_ping", "Failed to parse ping request")
		return
	}
	pongResponse := models.PongMessage{
		ID:        utils.PingIDGenerator.Generate(),
		Timestamp: time.Now().UnixMilli(),
		Type:      "pong",
	}

	payload, _ := json.Marshal(pongResponse)
	message.Client.SendMessage(payload)
}

// handleGRPCRequest обрабатывает gRPC запросы от клиента
func (p *MessageProcessor) handleGRPCRequest(message core.ClientMessage) {
	var grpcReq models.GRPCRequestMessage
	if err := json.Unmarshal(message.Data, &grpcReq); err != nil {
		p.sendError(message.Client, "invalid_grpc", "Failed to parse gRPC request")
		return
	}
	p.logger.WithFields(logrus.Fields{
		"client_id":  message.Client.ID(),
		"user_id":    message.Client.UserID(),
		"service":    grpcReq.Payload.Service,
		"action":     grpcReq.Payload.Action,
		"request_id": grpcReq.Payload.ID,
	}).Info("🔀 Processing gRPC request")

	response := p.router.RouteToService(grpcReq.Payload, message.Client)

	p.sendResponse(response, message.Client)

	p.logger.WithFields(logrus.Fields{
		"request_id": grpcReq.Payload.ID,
		"success":    response.Success,
	}).Info("✅ gRPC request processed")
}

func (p *MessageProcessor) sendResponse(response models.GRPCResponse, client *core.Client) {
	payload, err := json.Marshal(response)
	if err != nil {
		p.logger.WithError(err).Error("Failed to marshal response")
		return
	}

	client.SendMessage(payload)
	p.logger.WithFields(logrus.Fields{
		"client_id":  client.ID(),
		"request_id": response.ID,
	}).Debug("📨 Response sent to client")
}
