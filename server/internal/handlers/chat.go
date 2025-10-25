package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ASNMortred/AI-Hackathon/server/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ChatHandler struct {
	mcpServiceURL string
	httpClient    *http.Client
}

func NewChatHandler() *ChatHandler {
	mcpURL := os.Getenv("MCP_SERVICE_URL")
	if mcpURL == "" {
		mcpURL = "http://localhost:8000"
		logger.Logger.Warn("MCP_SERVICE_URL not set, using default", zap.String("url", mcpURL))
	}

	return &ChatHandler{
		mcpServiceURL: mcpURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type ChatRequest struct {
	MemoryId string `json:"memoryId" binding:"required"`
	Message  string `json:"message" binding:"required"`
}

type MCPChatRequest struct {
	Message     string  `json:"message"`
	SessionID   string  `json:"session_id,omitempty"`
	Temperature float64 `json:"temperature"`
}

type MCPChatResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

func (h *ChatHandler) Chat(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization header is required",
		})
		return
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid authorization header format",
		})
		return
	}

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

	logger.Logger.Info("received chat request",
		zap.String("message", req.Message),
		zap.String("memoryId", req.MemoryId),
	)

	mcpReq := MCPChatRequest{
		Message:     req.Message,
		SessionID:   req.MemoryId,
		Temperature: 0.7,
	}

	jsonData, err := json.Marshal(mcpReq)
	if err != nil {
		logger.Logger.Error("failed to marshal MCP request", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	mcpURL := fmt.Sprintf("%s/api/chat", h.mcpServiceURL)
	httpReq, err := http.NewRequestWithContext(c.Request.Context(), "POST", mcpURL, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Logger.Error("failed to create MCP request", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")

	logger.Logger.Info("forwarding request to MCP service", zap.String("url", mcpURL))

	resp, err := h.httpClient.Do(httpReq)
	if err != nil {
		logger.Logger.Error("failed to call MCP service",
			zap.Error(err),
			zap.String("url", mcpURL),
		)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Chat service temporarily unavailable",
		})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Logger.Error("failed to read MCP response", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	if resp.StatusCode != http.StatusOK {
		logger.Logger.Error("MCP service returned error",
			zap.Int("status", resp.StatusCode),
			zap.String("body", string(body)),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate response",
		})
		return
	}

	var mcpResp MCPChatResponse
	if err := json.Unmarshal(body, &mcpResp); err != nil {
		logger.Logger.Error("failed to parse MCP response", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	if !mcpResp.Success {
		logger.Logger.Error("MCP service reported failure", zap.String("error", mcpResp.Error))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": mcpResp.Error,
		})
		return
	}

	logger.Logger.Info("successfully processed chat request",
		zap.String("sessionId", req.MemoryId),
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    mcpResp.Data,
	})
}
