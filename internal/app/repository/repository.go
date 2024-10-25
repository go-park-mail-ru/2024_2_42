package repository

import (
	"pinset/internal/app/models"
	"pinset/internal/app/session"
	"sync"

	"github.com/minio/minio-go"
)

// Controllers
type (
	UserRepositoryController struct {
		mu *sync.RWMutex
		db map[string]*models.User
		sm *session.SessionsManager
	}

	FeedRepositoryController struct {
		mu *sync.RWMutex
	}

	MediaRepositoryController struct {
		client          *minio.Client
		ImageBucketName string
		VideoBucketName string
		AudioBucketName string
	}
)
