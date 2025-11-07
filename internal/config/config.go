package config

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSConfig представляет конфигурацию CORS
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           time.Duration
}

// NewCORSConfig создает новую конфигурацию CORS с дефолтными значениями
func NewCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
}

// ToGinHandler конвертирует конфигурацию в gin.HandlerFunc
func (c CORSConfig) ToGinHandler() gin.HandlerFunc {
	cfg := cors.Config{
		AllowOrigins:     c.AllowOrigins,
		AllowMethods:     c.AllowMethods,
		AllowHeaders:     c.AllowHeaders,
		ExposeHeaders:    c.ExposeHeaders,
		AllowCredentials: c.AllowCredentials,
		MaxAge:           c.MaxAge,
	}
	return cors.New(cfg)
}

// CorsConfig устаревшая функция, используйте NewCORSConfig().ToGinHandler()
// Deprecated: используйте NewCORSConfig().ToGinHandler()
func CorsConfig() gin.HandlerFunc {
	return NewCORSConfig().ToGinHandler()
}

