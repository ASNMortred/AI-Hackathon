package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ASNMortred/AI-Hackathon/internal/config"
	"github.com/ASNMortred/AI-Hackathon/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UploadHandler struct {
	config *config.Config
}

func NewUploadHandler(cfg *config.Config) *UploadHandler {
	return &UploadHandler{config: cfg}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		logger.Logger.Error("failed to get file from request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No file provided",
		})
		return
	}

	if file.Size > h.config.Upload.MaxSize {
		logger.Logger.Warn("file size exceeds limit",
			zap.String("filename", file.Filename),
			zap.Int64("size", file.Size),
			zap.Int64("max_size", h.config.Upload.MaxSize),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("File size exceeds limit of %d bytes", h.config.Upload.MaxSize),
		})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !h.isAllowedType(ext) {
		logger.Logger.Warn("file type not allowed",
			zap.String("filename", file.Filename),
			zap.String("extension", ext),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File type not allowed",
		})
		return
	}

	if err := os.MkdirAll(h.config.Upload.UploadDir, 0755); err != nil {
		logger.Logger.Error("failed to create upload directory", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create upload directory",
		})
		return
	}

	filename := filepath.Join(h.config.Upload.UploadDir, file.Filename)
	if err := c.SaveUploadedFile(file, filename); err != nil {
		logger.Logger.Error("failed to save file", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save file",
		})
		return
	}

	logger.Logger.Info("file uploaded successfully",
		zap.String("filename", file.Filename),
		zap.Int64("size", file.Size),
	)

	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"filename": file.Filename,
		"size":     file.Size,
		"path":     filename,
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
