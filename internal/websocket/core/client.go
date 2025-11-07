package core

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/srcgod/apigateway/internal/utils"
)

// Client представляет WebSocket соединение клиента
type Client struct {
	id     string
	userID int64
	logger *logrus.Logger

	conn *websocket.Conn

	Send    chan []byte
	Receive chan ClientMessage

	readTimeout  time.Duration
	writeTimeout time.Duration

	isAlive bool
}

// ClientMessage представляет сообщение от клиента
type ClientMessage struct {
	Type      int
	Data      []byte
	Client    *Client
	Timestamp time.Time
}

// NewClient создает новый клиент
func NewClient(ws *websocket.Conn, userID int64, logger *logrus.Logger) *Client {
	return &Client{
		id:           utils.ConnectionIDGenerator.Generate(),
		userID:       userID,
		conn:         ws,
		Send:         make(chan []byte, 256),
		Receive:      make(chan ClientMessage),
		logger:       logger,
		readTimeout:  time.Second * 60,
		writeTimeout: time.Second * 60,
		isAlive:      true,
	}
}

func (c *Client) ID() string {
	return c.id
}

func (c *Client) UserID() int64 {
	return c.userID
}

func (c *Client) SendMessage(data []byte) {
	select {
	case c.Send <- data:
	default:
		c.logger.WithFields(logrus.Fields{
			"client_id": c.id,
		}).Warn("⚠️ Send channel full, message dropped")
	}
}

func (c *Client) Close() {
	c.conn.Close()
}

// ReadPump читает сообщения из WebSocket соединения
func (c *Client) ReadPump(unregister func(*Client)) {
	defer func() {
		c.logger.WithFields(logrus.Fields{
			"client_id": c.id,
			"user_id":   c.userID,
		}).Info("🛑 ReadPump stopped")

		unregister(c)
		c.conn.Close()
	}()

	c.logger.WithFields(logrus.Fields{
		"client_id": c.id,
		"user_id":   c.userID,
	}).Info("📖 ReadPump started")

	for {
		c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))

		msgtype, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.WithFields(logrus.Fields{
					"client_id": c.id,
					"user_id":   c.userID,
					"error":     err.Error(),
				}).Warn("❌ WebSocket read error")
			} else {
				c.logger.WithFields(logrus.Fields{
					"client_id": c.id,
					"user_id":   c.userID,
				}).Info("🔌 WebSocket connection closed")
			}
			break
		}

		clientMsg := ClientMessage{
			Type:      msgtype,
			Data:      msg,
			Client:    c,
			Timestamp: time.Now(),
		}
		select {
		case c.Receive <- clientMsg:
			c.logger.WithFields(logrus.Fields{
				"client_id":    c.id,
				"user_id":      c.userID,
				"message_type": msgtype,
				"data_length":  len(msg),
			}).Debug("📨 Message sent to Receive channel")

		case <-time.After(100 * time.Millisecond):
			c.logger.WithFields(logrus.Fields{
				"client_id": c.id,
				"user_id":   c.userID,
			}).Warn("⚠️ Receive channel full, message dropped")
		}
	}
}

// WritePump записывает сообщения в WebSocket соединение
func (c *Client) WritePump() {
	ticker := time.NewTicker(time.Second * 50)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.Send:
			c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(msg)
			w.Close()
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
