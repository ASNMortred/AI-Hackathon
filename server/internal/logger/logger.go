package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger() error {
	if err := os.MkdirAll("logs", 0755); err != nil {
		return err
	}

	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.OutputPaths = []string{
		"stdout",
		"logs/app.log",
	}
	config.ErrorOutputPaths = []string{
		"stderr",
		"logs/error.log",
	}

	var err error
	Logger, err = config.Build()
	if err != nil {
		return err
	}

	return nil
}

func Sync() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}
