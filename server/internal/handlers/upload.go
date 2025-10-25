package handlers

import (
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/ASNMortred/AI-Hackathon/server/internal/config"
	"github.com/ASNMortred/AI-Hackathon/server/internal/dao"
	"github.com/ASNMortred/AI-Hackathon/server/internal/logger"
	"github.com/ASNMortred/AI-Hackathon/server/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UploadHandler struct {
	config *config.Config
	minio  *storage.MinioService
}

func NewUploadHandler(cfg *config.Config) *UploadHandler {
	minioSvc, err := storage.NewMinioService(cfg.Minio)
	if err != nil {
		logger.Logger.Fatal("failed to init minio: " + err.Error())
	}
	return &UploadHandler{config: cfg, minio: minioSvc}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	// 验证认证头，提取用户名
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" || !strings.HasPrefix(tokenParts[1], "token_") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
		return
	}
	username := strings.TrimPrefix(tokenParts[1], "token_")

	// 读取文件
	file, err := c.FormFile("file")
	if err != nil {
		logger.Logger.Error("failed to get file from request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	// 校验大小与类型
	if file.Size > h.config.Upload.MaxSize {
		logger.Logger.Warn("file size exceeds limit",
			zap.String("filename", file.Filename),
			zap.Int64("size", file.Size),
			zap.Int64("max_size", h.config.Upload.MaxSize),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("File size exceeds limit of %d bytes", h.config.Upload.MaxSize)})
		return
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !h.isAllowedType(ext) {
		logger.Logger.Warn("file type not allowed",
			zap.String("filename", file.Filename),
			zap.String("extension", ext),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "File type not allowed"})
		return
	}

	// 获取用户UID
	userDAO := dao.NewUserDAO()
	user, err := userDAO.GetUserByUsername(username)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	// 上传到 MinIO
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
		return
	}
	defer f.Close()

	objectName := fmt.Sprintf("uid_%d/%s", user.Uid, file.Filename)
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		if ext != "" {
			if mt := mime.TypeByExtension(ext); mt != "" {
				contentType = mt
			}
		}
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	url, err := h.minio.Upload(c.Request.Context(), objectName, f, file.Size, contentType)
	if err != nil {
		logger.Logger.Error("failed to upload to minio", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload to object storage"})
		return
	}

	// 写入数据库
	mfDAO := dao.NewMinioFileDAO()
	if err := mfDAO.Create(user.Uid, file.Filename, url); err != nil {
		logger.Logger.Error("failed to insert minio file record", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record file metadata"})
		return
	}

	logger.Logger.Info("file uploaded to minio successfully",
		zap.String("filename", file.Filename),
		zap.Int64("size", file.Size),
		zap.String("url", url),
		zap.Int("uid", user.Uid),
	)

	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"filename": file.Filename,
		"size":     file.Size,
		"url":      url,
	})
}

func (h *UploadHandler) isAllowedType(ext string) bool {
	for _, allowedType := range h.config.Upload.AllowedTypes {
		if ext == allowedType {
			return true
		}
	}
	return false
}
