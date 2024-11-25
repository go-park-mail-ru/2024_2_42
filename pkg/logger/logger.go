package logger

import (
	"fmt"
	"io"
	"os"
	"pinset/configs"

	"github.com/sirupsen/logrus"
)

func NewLogger() (*logrus.Logger, error) {
	cfg := configs.NewLoggerParams()

	logFile, err := os.OpenFile(cfg.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("logrus: %w", err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	logger := &logrus.Logger{
		Out:   mw,
		Level: logrus.DebugLevel,
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		},
	}

	return logger, nil
}
