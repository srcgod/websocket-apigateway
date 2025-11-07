package services

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/srcgod/apigateway/internal/grpcclient"
	"github.com/srcgod/apigateway/internal/websocket/models"
	"github.com/srcgod/apigateway/internal/websocket/types"
	authv1 "github.com/srcgod/authproto/gen/user"
)

type AuthServiceHandler struct {
	logger       *logrus.Logger
	authGrpcInst authv1.AuthClient
}

func NewAuthServicesHandler(logger *logrus.Logger) *AuthServiceHandler {
	authGrpcInst, err := grpcclient.NewAuthClient("localhost:8081")
	if err != nil {
		logger.Error("creating auth grpc instance error", err)
	}
	return &AuthServiceHandler{logger: logger, authGrpcInst: authGrpcInst}
}

func (a *AuthServiceHandler) Handle(req models.GRPCRequest, client types.ClientInterface) models.GRPCResponse {
	switch req.Action {
	case "register":
		return a.HandleRegister(req)
	default:
		return models.GRPCResponse{}
	}

}

func (a *AuthServiceHandler) HandleRegister(req models.GRPCRequest) models.GRPCResponse {
	var regReq models.AuthRegisterRequest

	req.ParseData(&regReq)

	request_to_grpc := &authv1.RegisterRequest{
		Email:    regReq.Email,
		Password: regReq.Password,
		Username: regReq.Username,
	}
	// TODO: upgrade handle error
	response, err := a.authGrpcInst.Register(context.Background(), request_to_grpc)
	if err != nil {
		return a.errorResponse(req.ID, "registration_failed", "Service unavailable")

	}
	return models.GRPCResponse{
		ID:      req.ID,
		Success: true,
		Data:    response,
	}
}

func (a *AuthServiceHandler) errorResponse(id, code, message string) models.GRPCResponse {
	return models.GRPCResponse{
		ID:      id,
		Success: false,
		Error:   &models.ErrorResponse{Code: code, Message: message},
	}
}
