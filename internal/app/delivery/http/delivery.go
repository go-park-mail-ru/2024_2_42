package delivery

import (
	"mime/multipart"
	"pinset/internal/app/models"
	"pinset/internal/app/models/request"

	"github.com/sirupsen/logrus"
)

// Usecase interfaces
type (
	UserUsecase interface {
		LogIn(request.LoginRequest) (string, error)
		LogOut(string) error
		SignUp(user *models.User) error
		IsAuthorized(string) (float64, error)
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
		Logger *logrus.Logger
	}

	FeedDeliveryController struct {
		Usecase FeedUsecase
		Logger *logrus.Logger
	}

	MediaDeliveryController struct {
		Usecase MediaUsecase
		Logger *logrus.Logger
	}
)
