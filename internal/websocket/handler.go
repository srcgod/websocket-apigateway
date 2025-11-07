package websock

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/srcgod/apigateway/internal/websocket/core"
)

type WSocketHandler struct {
	hub       *core.Hub
	processor *MessageProcessor
	logger    *logrus.Logger
}

func New(logger *logrus.Logger, hub *core.Hub) *WSocketHandler {
	router := NewServiceRouter(logger)
	processor := NewMessageProcessor(router, logger)

	return &WSocketHandler{
		hub:       hub,
		processor: processor,
		logger:    logger,
	}
}

func (h *WSocketHandler) SetUpRoutes(r *gin.Engine) {
	wsGroup := r.Group("/gateway")

	//wsGroup.Use(jwtMiddleware.JWTAuthMiddleware(), middleware.UIDMiddleware())
	{
		wsGroup.GET("/ws", h.HandleWebSocket)
	}
}

func (h *WSocketHandler) HandleWebSocket(c *gin.Context) {
	//uid, _ := c.Get(middleware.UidInt64Key)
	//uidInt64, ok := uid.(int64)
	var uidInt64 int64 = 2

	/* if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uid type"})
		return
	}
	*/
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error create websocket": err})
		return
	}

	client := core.NewClient(conn, uidInt64, h.logger)
	h.hub.RegisterClient(client)

	go h.processor.ProcessClientMessages(client)
	go client.ReadPump(h.hub.UnregisterClient)
	go client.WritePump()

}
