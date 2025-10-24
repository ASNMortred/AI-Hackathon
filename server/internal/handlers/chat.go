package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ASNMortred/AI-Hackathon/server/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ChatHandler struct{}

func NewChatHandler() *ChatHandler {
	return &ChatHandler{}
}

type ChatRequest struct {
	MemoryId string `json:"memoryId" binding:"required"`
	Message  string `json:"message" binding:"required"`
}

func (h *ChatHandler) Chat(c *gin.Context) {
	// 验证认证头
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization header is required",
		})
		return
	}

	// 简单的token验证（实际项目中应该使用JWT或其他安全的认证机制）
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid authorization header format",
		})
		return
	}

	// 检查token前缀（实际项目中应该验证token的有效性）
	if !strings.HasPrefix(tokenParts[1], "token_") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
		})
		return
	}

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
		zap.String("memoryId", req.MemoryId),
	)

	// 模拟流式响应
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// 简单的响应示例
	response := "感谢您的消息: " + req.Message
	c.String(http.StatusOK, response)
}
