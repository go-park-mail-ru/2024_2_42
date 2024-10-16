package usecase

import (
	"pinset/internal/app/repository"
	"pinset/internal/models"
	"pinset/internal/models/request"
)

// Interfaces
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
	userUsecaseController struct {
		repo repository.UserRepository
	}

	feedUsecaseController struct {
		repo repository.FeedRepository
	}
)
