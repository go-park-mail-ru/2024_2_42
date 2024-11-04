package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func InitDB(logger *logrus.Logger) *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatal(err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbSSLMode := os.Getenv("DB_SSLMODE")
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	postgresDB, err := sql.Open("pgx", connStr)
	if err != nil {
		logger.Fatalf("Error starting postgres: %v", err)
	}

	retrying := 10
	i := 1
	logger.Info("try ping", i)
	for err = postgresDB.Ping(); err != nil; err = postgresDB.Ping() {
		if i >= retrying {
			logger.Fatal(err)
		}
		i++
		time.Sleep(1 * time.Second)
		log.Printf("try ping postgresql: %v", i)
	}

	logger.Info("PostgreSQL succesful connected!")
	return postgresDB
}
