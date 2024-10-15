package usecase

import (
	"net/http"
	"pinset/internal/app/repository"
	"pinset/internal/models"
	"pinset/internal/models/request"
	"time"

	internal_errors "pinset/internal/errors"

	"github.com/golang-jwt/jwt"
)

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecaseController{
		repo: repo,
	}
}

func (uuc *userUsecaseController) LogIn(w http.ResponseWriter, req request.LoginRequest) error {
	user := models.User{
		Email:    req.Email,
		Password: req.Password,
	}

	user.UserID = getUserID(user)

	// User is not registered
	if !uuc.repo.UserAlreadySignedUp(user) {
		return internal_errors.ErrUserIsNotRegistered
	}

	// Does user already have an active session?
	if uuc.repo.UserHasActiveSession(req.Token) {
		return internal_errors.ErrUserAlreadyAuthorized
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.UserID,
		"login":   req.Email,
		"exp":     time.Now().Add(repository.SessionTokenExpirationTime).Unix(),
	})

	signedToken, err := token.SignedString(repository.SECRET)
	if err != nil {
		return internal_errors.ErrCantSignSessionToken
	}

	uuc.sm.Create(w, user.UserID, signedToken)

	return nil
}

func (uuc *userUsecaseController) LogOut(token string) error {
	// Need to remove user from authorized list
	if !uuc.sm.Exists(token) {
		return internal_errors.ErrUserIsNotAuthorized
	}

	uuc.sm.Remove(token)

	return nil
}

func (uuc *userUsecaseController) IsAuthorized(cookie *http.Cookie) (float64, error) {
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return repository.SECRET, nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, internal_errors.ErrInvalidSessionToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims["user_id"].(float64), nil
	}

	return 0, internal_errors.ErrBadRequest
}
