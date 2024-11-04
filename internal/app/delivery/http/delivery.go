package delivery

import (
	"mime/multipart"
	"pinset/internal/app/models"
	"pinset/internal/app/models/request"
	"pinset/internal/app/models/response"

	"github.com/sirupsen/logrus"
)

// Usecase interfaces
type (
	UserUsecase interface {
		LogIn(request.LoginRequest) (string, error)
		LogOut(string) error
		SignUp(user *models.User) error
		IsAuthorized(string) (uint64, error)
		GetUserInfo(*models.User) (response.UserProfileResponse, error)
		UpdateUserInfo(string, *models.User) error
	}

	FeedUsecase interface {
		Feed() models.Feed
	}

	MediaUsecase interface {
		GetMedia() error
		UploadMedia(files []*multipart.FileHeader) ([]string, error)
	}
)

// Controllers
type (
	UserDeliveryController struct {
		Usecase UserUsecase
		Logger  *logrus.Logger
	}

	FeedDeliveryController struct {
		Usecase FeedUsecase
		Logger  *logrus.Logger
	}

	MediaDeliveryController struct {
		Usecase MediaUsecase
		Logger  *logrus.Logger
	}
)
