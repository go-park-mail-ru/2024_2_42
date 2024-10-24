package usecase

import (
	"pinset/configs"
	"pinset/internal/app/models"
	"pinset/internal/app/session"
)

// Repository interfaces
type (
	UserRepository interface {
		Insert(*models.User) error
		UserHasActiveSession(string) bool
		UserAlreadySignedUp(models.User) bool
		GetUserId(models.User) uint64
		Session() *session.SessionsManager
	}

	FeedRepository interface {
		GetPins() []models.Pin
		InsertPin(models.Pin)
	}

	MediaRepository interface {
		GetMedia(string, string) error
		UploadMedia(string, string) error
	}
)

// Controllers
type (
	userUsecaseController struct {
		repo UserRepository
		authParameters configs.AuthParams
	}

	feedUsecaseController struct {
		repo FeedRepository
	}
)
