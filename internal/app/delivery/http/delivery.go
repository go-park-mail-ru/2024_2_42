package delivery

import (
	"mime/multipart"
	"pinset/internal/app/models"
	"pinset/internal/app/models/request"
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
		UploadMedia(files []*multipart.FileHeader) error
	}
)

// Controllers
type (
	UserDeliveryController struct {
		Usecase UserUsecase
	}

	FeedDeliveryController struct {
		Usecase FeedUsecase
	}

	MediaDeliveryController struct {
		Usecase MediaUsecase
	}
)
