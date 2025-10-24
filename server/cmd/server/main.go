package main

import (
	"fmt"
	"os"

	"github.com/ASNMortred/AI-Hackathon/server/internal/config"
	"github.com/ASNMortred/AI-Hackathon/server/internal/handlers"
	"github.com/ASNMortred/AI-Hackathon/server/internal/logger"
	"github.com/ASNMortred/AI-Hackathon/server/internal/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := logger.InitLogger(); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Logger.Fatal("Failed to load config: " + err.Error())
	}

	router := gin.New()
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.LoggerMiddleware())

	uploadHandler := handlers.NewUploadHandler(cfg)
	playHandler := handlers.NewPlayHandler()
	chatHandler := handlers.NewChatHandler()

	v1 := router.Group("/api/v1")
	{
		v1.POST("/upload", uploadHandler.Upload)
		v1.GET("/play/:videoID", playHandler.Play)
		v1.POST("/chat", chatHandler.Chat)
	}

	addr := ":" + cfg.Server.Port
	logger.Logger.Info("Starting server on " + addr)
	if err := router.Run(addr); err != nil {
		logger.Logger.Fatal("Failed to start server: " + err.Error())
	}
}
