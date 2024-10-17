package repository

import (
	"pinset/internal/app/models"
	"pinset/internal/app/session"
	"sync"
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
)
