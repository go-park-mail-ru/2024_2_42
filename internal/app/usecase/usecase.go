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
		CreateUser(*models.User) error
		CheckUserByEmail(*models.User) (bool, error)
		GetUserInfoByID(uint64) (models.User, error)
		CheckUserCredentials(*models.User) error
		UpdateUserInfoByID(*models.User) error
		UpdateUserPasswordByID(*models.User) error
		DeleteUserByID(uint64) error

		FollowUser(uint64, uint64) error
		UnfollowUser(uint64, uint64) error
		GetAllFollowings(uint64, uint64) ([]uint64, error)
		GetAllSubscriptions(uint64, uint64) ([]uint64, error)

		UserHasActiveSession(string) bool
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
