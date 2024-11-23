package configs

import (
	"os"
	"strconv"
	"time"
)

func LookUpStringEnvVar(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func LookUpBoolEnvVar(key string, defaultValue bool) bool {
	valStr := LookUpStringEnvVar(key, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultValue
}

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

const (
	loggerfilePath = "./logs/pinset.log"
)

type LoggerParams struct {
	FilePath string
}

func NewLoggerParams() LoggerParams {
	return LoggerParams{
		FilePath: loggerfilePath,
	}
}

type ctxUserIDKeyType string

const UserIdKey ctxUserIDKeyType = "user_id"
