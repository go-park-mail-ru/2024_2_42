package delivery

import (
	"net/http"
	"pinset/internal/app/usecase"
)

// Interfaces
type (
	UserDelivery interface {
		LogIn(w http.ResponseWriter, r *http.Request)
		LogOut(w http.ResponseWriter, r *http.Request)
		SignUp(w http.ResponseWriter, r *http.Request)
		IsAuthorized(w http.ResponseWriter, r *http.Request)
	}

	FeedDelivery interface {
		Feed(w http.ResponseWriter, r *http.Request)
	}
)

// Controllers
type (
	UserDeliveryController struct {
		usecase usecase.UserUsecase
	}

	FeedDeliveryController struct {
		usecase usecase.FeedUsecase
	}
)
