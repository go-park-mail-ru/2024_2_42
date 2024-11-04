package usecase

import (
	"fmt"
	"pinset/configs"
	delivery "pinset/internal/app/delivery/http"
	"pinset/internal/app/models"
	"pinset/internal/app/models/request"
	"pinset/internal/app/models/response"
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

	// User is not registered
	isUserExists, err := uuc.repo.CheckUserByEmail(&user)
	if err != nil {
		return "", fmt.Errorf("signUp after UserAlreadySignedUp: %w", err)
	}
	if !isUserExists {
		return "", internal_errors.ErrUserDoesntExists
	}

	err = uuc.repo.CheckUserCredentials(&user)
	if err != nil {
		return "", fmt.Errorf("login checkCredentials: %w", err)
	}

	// Does user already have an active session?
	if uuc.repo.UserHasActiveSession(req.Token) {
		return "", internal_errors.ErrUserAlreadyAuthorized
	}

	userID, err := uuc.repo.GetLastUserID()
	if err != nil {
		return "", fmt.Errorf("logIn getLastUserID: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
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
	isUserExists, err := uuc.repo.CheckUserByEmail(user)
	if err != nil {
		return fmt.Errorf("signUp after UserAlreadySignedUp: %w", err)
	}

	if isUserExists {
		return internal_errors.ErrUserAlreadyExists
	}

	err = uuc.repo.CreateUser(user)
	if err != nil {
		return fmt.Errorf("signUp after Insert: %w", err)
	}

	return nil
}

func (uuc *UserUsecaseController) IsAuthorized(token string) (uint64, error) {
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
		return uint64(claims["user_id"].(float64)), nil
	}

	return 0, internal_errors.ErrBadRequest
}

func (uuc *UserUsecaseController) UpdateUserInfo(token string, user *models.User) error {
	userID, err := uuc.IsAuthorized(token)

	if err != nil {
		return fmt.Errorf("updateUserPasswordByID isAuthorized: %w", err)
	}

	if userID != user.UserID {
		return internal_errors.ErrBadUserID
	}

	err = uuc.repo.UpdateUserInfo(user)
	if err != nil {
		return fmt.Errorf("getUserInfoByID usecase: %w", err)
	}
	return nil
}

func (uuc *UserUsecaseController) UpdateUserPassword(token string, user *models.User) error {
	userID, err := uuc.IsAuthorized(token)

	if err != nil {
		return fmt.Errorf("updateUserPasswordByID isAuthorized: %w", err)
	}

	if userID != user.UserID {
		return internal_errors.ErrBadUserID
	}

	err = uuc.repo.UpdateUserPassword(user)
	if err != nil {
		return fmt.Errorf("getUserInfoByID updateUserPasswordByID: %w", err)
	}
	return nil
}

func (uuc *UserUsecaseController) DeleteProfile(token string, user *models.User) error {
	_, err := uuc.IsAuthorized(token)

	if err != nil {
		return fmt.Errorf("deleteUserByID isAuthorized: %w", err)
	}

	isUserExists, err := uuc.repo.CheckUserByEmail(user)
	if err != nil {
		return fmt.Errorf("signUp after UserAlreadySignedUp: %w", err)
	}

	if !isUserExists {
		return internal_errors.ErrUserDoesntExists
	}

	err = uuc.repo.DeleteUserByID(user.UserID)
	if err != nil {
		return fmt.Errorf("getUserInfoByID usecase: %w", err)
	}
	return nil
}

func (uuc *UserUsecaseController) GetUserInfo(user *models.User) (response.UserProfileResponse, error) {
	var userProfilelInfo response.UserProfileResponse
	userProfilelInfo, err := uuc.repo.GetUserInfo(user)
	if err != nil {
		return response.UserProfileResponse{}, fmt.Errorf("userProfile GetUserInfo usecase: %w", err)
	}

	userProfilelInfo.NumOfUserFollowings, err = uuc.repo.GetFollowingsCount(user.UserID)
	if err != nil {
		return response.UserProfileResponse{}, fmt.Errorf("userProfile GetFollowingsCount usecase: %w", err)
	}

	userProfilelInfo.NumOfUserSubscriptions, err = uuc.repo.GetlSubsriptionsCount(user.UserID)
	if err != nil {
		return response.UserProfileResponse{}, fmt.Errorf("userProfile GetFollowingsCount usecase: %w", err)
	}

	return userProfilelInfo, nil
}
