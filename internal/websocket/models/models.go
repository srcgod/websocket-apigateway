package models

import (
	"encoding/json"
	"fmt"
)

type GRPCRequestMessage struct {
	Type    string `json:"type"`
	Payload GRPCRequest
}

type GRPCRequest struct {
	ID       string          `json:"id"`
	Service  string          `json:"service"`
	Action   string          `json:"action"`
	Metadata map[string]any  `json:"metadata"`
	Data     json.RawMessage `json:"data"`
}

func (r *GRPCRequest) ParseData(target any) error {
	if len(r.Data) == 0 {
		return fmt.Errorf("empty data")
	}
	return json.Unmarshal(r.Data, target)
}

func (r *GRPCRequest) MustParseData(target any) {
	if err := r.ParseData(target); err != nil {
		panic(fmt.Sprintf("Failed to parse request data: %v", err))
	}
}

type GRPCResponse struct {
	ID      string         `json:"id"`
	Success bool           `json:"success"`
	Data    any            `json:"data,omitempty"`
	Error   *ErrorResponse `json:"error,omitempty"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type MessageServiceData struct {
	Text   string `json:"text"`
	RoomID string `json:"room_id,omitempty"`
	Test   string `json:"test,omitempty"`
}

type PingMessage struct {
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`
}

type PongMessage struct {
	Type      string `json:"type"`
	ID        string `json:"ID"`
	Timestamp int64  `json:"timestamp"`
}
