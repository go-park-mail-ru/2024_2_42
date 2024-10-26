package configs

import (
	"os"
	"time"
)

type internalParams struct {
	MainServerPort string
}

func NewInternalParams() internalParams {
	internalParams := internalParams{}

	internalParams.MainServerPort = ":8080"

	return internalParams
}

type AuthParams struct {
	SessionTokenExpirationTime time.Duration
	JwtSecret                  []byte
}

func NewAuthParams() AuthParams {
	return AuthParams{
		SessionTokenExpirationTime: time.Hour * 72,
		JwtSecret:                  []byte(os.Getenv("JWT_SECRET")),
	}
}

const filePath = "./logs/log.log"

type LoggerParams struct {
	FilePath string
}

func NewLoggerParams() LoggerParams {
	return LoggerParams{
		FilePath: filePath,
	}
}
