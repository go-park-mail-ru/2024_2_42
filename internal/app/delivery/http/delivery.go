package delivery

import (
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
)

// Controllers
type (
	UserDeliveryController struct {
		Usecase UserUsecase
	}

	FeedDeliveryController struct {
		Usecase FeedUsecase
	}
)
