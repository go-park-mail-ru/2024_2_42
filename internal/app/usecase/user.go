package usecase

import (
	"pinset/configs"
	delivery "pinset/internal/app/delivery/http"
	"pinset/internal/app/models"
	"pinset/internal/app/models/request"
	"time"

	internal_errors "pinset/internal/errors"

	"github.com/golang-jwt/jwt"
)

func NewUserUsecase(repo UserRepository) delivery.UserUsecase {
	return &UserUsecaseController{
		repo:           repo,
		authParameters: configs.NewAuthParams(),
	}
}

func (uuc *UserUsecaseController) LogIn(req request.LoginRequest) (string, error) {
	user := models.User{
		Email:    req.Email,
		Password: req.Password,
	}

	user.UserID = uuc.repo.GetUserId(user)

	// User is not registered
	if !uuc.repo.UserAlreadySignedUp(user) {
		return "", internal_errors.ErrUserIsNotRegistered
	}

	// Does user already have an active session?
	if uuc.repo.UserHasActiveSession(req.Token) {
		return "", internal_errors.ErrUserAlreadyAuthorized
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.UserID,
		"login":   req.Email,
		"exp":     time.Now().Add(uuc.authParameters.SessionTokenExpirationTime).Unix(),
	})

	signedToken, err := token.SignedString(uuc.authParameters.JwtSecret)
	if err != nil {
		return "", internal_errors.ErrCantSignSessionToken
	}

	uuc.repo.Session().Create(signedToken, user.UserID)

	return signedToken, nil
}

func (uuc *UserUsecaseController) LogOut(token string) error {
	// Need to remove user from authorized list
	if !uuc.repo.Session().Exists(token) {
		return internal_errors.ErrUserIsNotAuthorized
	}

	uuc.repo.Session().Remove(token)

	return nil
}

func (uuc *UserUsecaseController) SignUp(user *models.User) error {
	// Incorrect data given
	if err := user.Valid(); err != nil {
		return internal_errors.ErrUserDataInvalid
	}

	// User already registered
	if uuc.repo.UserAlreadySignedUp(*user) {
		return internal_errors.ErrUserAlreadyRegistered
	}

	err := uuc.repo.Insert(user)
	if err != nil {
		return err
	}

	return nil
}

func (uuc *UserUsecaseController) IsAuthorized(token string) (float64, error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return uuc.authParameters.JwtSecret, nil
	})

	if err != nil {
		return 0, err
	}

	if !jwtToken.Valid {
		return 0, internal_errors.ErrInvalidSessionToken
	}

	if !uuc.repo.UserHasActiveSession(token) {
		return 0, internal_errors.ErrUserIsNotAuthorized
	}

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok {
		return claims["user_id"].(float64), nil
	}

	return 0, internal_errors.ErrBadRequest
}
