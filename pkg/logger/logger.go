package logger

import (
	"io"
	"os"
	"pinset/configs"

	"github.com/sirupsen/logrus"
)

func Logger() (*logrus.Logger, error) {
	cfg := configs.NewLoggerParams()

	logFile, err := os.OpenFile(cfg.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	logger := logrus.New()

	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logger.Level = logrus.DebugLevel

	mw := io.MultiWriter(os.Stdout, logFile)
	logrus.SetOutput(mw)

	return logger, err
}
