package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/srcgod/apigateway/internal/config"
)

// Server представляет HTTP сервер
type Server struct {
	router *gin.Engine
	logger *logrus.Logger
	config *Config
}

// Config конфигурация сервера
type Config struct {
	Host         string
	Port         string
	CORSConfig   config.CORSConfig
	ReadTimeout  int
	WriteTimeout int
}

// New создает новый HTTP сервер
func New(logger *logrus.Logger, router *gin.Engine, cfg *Config) *Server {
	// Применяем CORS middleware
	router.Use(cfg.CORSConfig.ToGinHandler())

	return &Server{
		router: router,
		logger: logger,
		config: cfg,
	}
}

// Start запускает HTTP сервер
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
	s.logger.WithFields(logrus.Fields{
		"host": s.config.Host,
		"port": s.config.Port,
	}).Info("🚀 Starting HTTP server")

	if err := s.router.Run(addr); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// GetRouter возвращает роутер для настройки маршрутов
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}

