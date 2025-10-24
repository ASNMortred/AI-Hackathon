package handlers

import (
	"fmt"
	"net/http"

	"github.com/ASNMortred/AI-Hackathon/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ChatHandler struct{}

func NewChatHandler() *ChatHandler {
	return &ChatHandler{}
}

type ChatRequest struct {
	Message string `json:"message" binding:"required"`
}

func (h *ChatHandler) Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Error("failed to parse chat request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	fmt.Printf("接收到用户消息：%s\n", req.Message)

	logger.Logger.Info("received user message",
		zap.String("message", req.Message),
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Message received",
		"echo":    req.Message,
	})
}
