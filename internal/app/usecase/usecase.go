package usecase

import (
	"io"
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
		GetBucketNameForContentType(fileType string) string
		HasCorrectContentType(string) bool
		GetMedia(string, string) ([]byte, error)
		UploadMedia(string, string, io.Reader, int64) (string, error)
	}
)

// Controllers
type (
	UserUsecaseController struct {
		repo           UserRepository
		authParameters configs.AuthParams
	}

	FeedUsecaseController struct {
		repo FeedRepository
	}

	MediaUsecaseController struct {
		repo MediaRepository
	}
)
