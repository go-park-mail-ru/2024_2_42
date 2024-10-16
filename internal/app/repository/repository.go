package repository

import (
	"pinset/internal/models"
	"sync"
)

// Interfaces
type (
	UserRepository interface {
		Insert(*models.User) error
		UserHasActiveSession(string) bool
		UserAlreadySignedUp(models.User) bool
		GetUserId(models.User) uint64
		Session() *SessionsManager
	}

	FeedRepository interface {
		GetPins() []models.Pin
		InsertPin(models.Pin)
	}
)

// Controllers
type (
	UserRepositoryController struct {
		mu *sync.RWMutex
		db map[string]*models.User
		sm *SessionsManager
	}

	FeedRepositoryController struct {
		mu *sync.RWMutex
	}
)
