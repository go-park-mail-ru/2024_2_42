package usecase

import (
	"io"
	"pinset/configs"
	"pinset/internal/app/models"
	"pinset/internal/app/models/response"
	"pinset/internal/app/session"
)

// Repository interfaces
type (
	UserRepository interface {
		GetLastUserID() (uint64, error)
		CreateUser(*models.User) error
		CheckUserByEmail(*models.User) (bool, error)
		GetUserInfo(*models.User) (response.UserProfileResponse, error)
		CheckUserCredentials(*models.User) error
		UpdateUserInfo(*models.User) error
		UpdateUserPassword(*models.User) error
		DeleteUserByID(uint64) error

		FollowUser(uint64, uint64) error
		UnfollowUser(uint64, uint64) error
		GetAllFollowings(uint64, uint64) ([]uint64, error)
		GetAllSubscriptions(uint64, uint64) ([]uint64, error)
		GetFollowingsCount(uint64) (uint64, error)
		GetlSubsriptionsCount(uint64) (uint64, error)

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
