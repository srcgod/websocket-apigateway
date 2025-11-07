package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/srcgod/apigateway/internal/app"
	"github.com/srcgod/apigateway/internal/server"
)

func main() {
	logger := initLogger()

	if err := godotenv.Load(); err != nil {
		logger.Debug("No .env file found, using environment variables")
	}

	application := app.New(logger)

	router := application.SetupRoutes()

	serverConfig := &server.Config{
		Host:       getEnv("HOST", "localhost"),
		Port:       getEnv("PORT", "8080"),
		CORSConfig: application.GetCORSConfig(),
	}

	srv := server.New(logger, router.Engine(), serverConfig)

	if err := srv.Start(); err != nil {
		logger.WithError(err).Fatal("Failed to start server")
	}
}

func initLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	logLevel := getEnv("LOG_LEVEL", "info")
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logger.Warnf("Invalid log level %s, using info", logLevel)
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	return logger
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
