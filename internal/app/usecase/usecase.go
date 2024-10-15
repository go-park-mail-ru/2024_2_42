package usecase

import (
	"net/http"
	"pinset/internal/app/repository"
	"pinset/internal/models"
	"pinset/internal/models/request"
)

// Interfaces
type (
	UserUsecase interface {
		LogIn(http.ResponseWriter, request.LoginRequest) error
		LogOut(string) error
		IsAuthorized(*http.Cookie) (float64, error)
	}

	FeedUsecase interface {
		Feed() models.Feed
	}
)

// Controllers
type (
	userUsecaseController struct {
		repo repository.UserRepository
		sm   *repository.SessionsManager
	}

	feedUsecaseController struct {
		repo repository.FeedRepository
	}
)
