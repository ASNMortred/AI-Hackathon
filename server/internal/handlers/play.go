package handlers

import (
	"fmt"
	"net/http"

	"github.com/ASNMortred/AI-Hackathon/server/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PlayHandler struct{}

func NewPlayHandler() *PlayHandler {
	return &PlayHandler{}
}

func (h *PlayHandler) Play(c *gin.Context) {
	videoID := c.Param("videoID")

	fmt.Printf("正在播放视频：%s\n", videoID)

	logger.Logger.Info("playing video",
		zap.String("video_id", videoID),
	)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Playing video",
		"video_id": videoID,
	})
}
