package logger

import (
	"fmt"
	"os"
	"pinset/configs"

	"github.com/sirupsen/logrus"
)

func Logger() error {
	cfg := configs.NewLoggerParams()
	f, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		return fmt.Errorf("Logger: %w", err)
	}
	logrus.SetOutput(f)
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	return err
}
