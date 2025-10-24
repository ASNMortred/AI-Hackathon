package handlers

import (
	"net/http"

	"github.com/ASNMortred/AI-Hackathon/server/internal/logger"
	"github.com/ASNMortred/AI-Hackathon/server/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService *services.UserService
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	Message  string `json:"message"`
	Username string `json:"username"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Message  string `json:"message"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

// NewUserHandler 创建新的用户处理器
func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: services.NewUserService(),
	}
}

// Register 处理用户注册请求
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Error("failed to parse register request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// 注册用户
	err := h.userService.RegisterUser(req.Username, req.Password)
	if err != nil {
		logger.Logger.Error("failed to register user",
			zap.String("username", req.Username),
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	logger.Logger.Info("user registered successfully", zap.String("username", req.Username))

	c.JSON(http.StatusOK, RegisterResponse{
		Message:  "User registered successfully",
		Username: req.Username,
	})
}

// Login 处理用户登录请求
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Error("failed to parse login request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// 验证用户身份
	user, err := h.userService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		logger.Logger.Warn("failed to authenticate user",
			zap.String("username", req.Username),
			zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	}

	// 生成简单token（实际项目中应使用JWT或其他安全的token机制）
	token := "token_" + user.Username

	logger.Logger.Info("user logged in successfully",
		zap.String("username", req.Username),
		zap.Int("user_id", user.Uid))

	c.JSON(http.StatusOK, LoginResponse{
		Message:  "Login successful",
		Username: user.Username,
		Token:    token,
	})
}
