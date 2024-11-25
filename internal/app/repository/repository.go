package repository

import (
	"database/sql"
	"pinset/internal/app/session"
	"sync"

	"github.com/minio/minio-go"
)

// Controllers
type (
	UserRepositoryController struct {
		mu *sync.RWMutex
		db *sql.DB
		sm *session.SessionsManager
	}

	MediaRepositoryController struct {
		client          *minio.Client
		ImageBucketName string
		VideoBucketName string
		AudioBucketName string
	}
)
