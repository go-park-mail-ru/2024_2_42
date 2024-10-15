package repository

import (
	"pinset/internal/models"
	"sync"
)

// Interfaces
type (
	UserRepository interface {
		UserHasActiveSession(string) bool
		UserAlreadySignedUp(models.User) bool
		GetUserId(models.User) uint64
	}

	FeedRepository interface {
		GetPins() []models.Pin
		InsertPin(pin models.Pin)
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
