package user_repository

import (
	"database/sql"
	"errors"
	"fmt"
	"pinset/internal/app/models"
	"pinset/internal/app/session"
	"pinset/internal/app/usecase"
	internal_errors "pinset/internal/errors"

	"github.com/sirupsen/logrus"
)

type UserRepositoryController struct {
	db     *sql.DB
	sm     *session.SessionsManager
	logger *logrus.Logger
}

func NewUserRepository(db *sql.DB, logger *logrus.Logger) usecase.UserRepository {
	return &UserRepositoryController{
		db:     db,
		logger: logger,
		sm:     session.NewSessionManager(),
	}
}

func (urc *UserRepositoryController) CreateUser(user *models.User) error {
	var userID uint64
	err := urc.db.QueryRow(CreateUser, user.UserName, user.NickName, user.Email, user.Password).Scan(&userID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("psql CreateUser: %w", err)
		}

		if errors.Is(err, sql.ErrNoRows) {
			userID = 0
		}
	}

	if userID == 0 {
		return internal_errors.ErrBadUserInputData
	}
	urc.logger.WithField("user_id", userID).Info("insert func")
	return nil
}

func (urc *UserRepositoryController) CheckUserByEmail(user *models.User) (bool, error) {
	var foundEmail sql.NullInt64

	err := urc.db.QueryRow(CheckUserByEmail, user.Email).Scan(&foundEmail)
	if err != nil && err != sql.ErrNoRows {
		return foundEmail.Valid, internal_errors.ErrUserDoesntExists
	}

	if foundEmail.Valid != true {
		return foundEmail.Valid, internal_errors.ErrUserDoesntExists
	}
	return foundEmail.Valid, nil
}

func (urc *UserRepositoryController) CheckUserCredentials(user *models.User) error {
	var userPassword string
	err := urc.db.QueryRow(CheckUserCredentials, user.UserID).Scan(&userPassword)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return internal_errors.ErrUserDoesntExists
		}
		return fmt.Errorf("psql GetUserByID: %w", err)
	}

	if userPassword != user.Password {
		return internal_errors.ErrBadPassword
	}
	return nil
}

func (urc *UserRepositoryController) GetUserInfoByID(userID uint64) (models.User, error) {
	var userInfo models.User

	err := urc.db.QueryRow(GetUserInfoByID, userID).Scan(&userInfo.UserName, &userInfo.NickName, &userInfo.Description, &userInfo.Gender, &userInfo.BirthTime, &userInfo.AvatarUrl)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return models.User{}, internal_errors.ErrUserDoesntExists
		}
		return models.User{}, fmt.Errorf("psql GetUserByID: %w", err)
	}
	return userInfo, nil
}

func (urc *UserRepositoryController) UpdateUserInfoByID(user *models.User) error {
	var userID uint64
	err := urc.db.QueryRow(UpdateUserInfoByID, user.NickName, user.Description, user.BirthTime, user.Gender, user.UserID).Scan(&userID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return internal_errors.ErrUserDoesntExists
		}
		return fmt.Errorf("psql UpdateUserInfo: %w", err)
	}

	return nil
}

func (urc *UserRepositoryController) UpdateUserPasswordByID(user *models.User) error {
	var userID uint64
	err := urc.db.QueryRow(UpdateUserPasswordByID, user.UserID).Scan(&userID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return internal_errors.ErrUserDoesntExists
		}
		return fmt.Errorf("psql UpdateUserPassword: %w", err)
	}

	return nil
}

func (urc *UserRepositoryController) DeleteUserByID(userID uint64) error {
	_, err := urc.db.Exec(UpdateUserPasswordByID, userID)
	if err != nil {
		return internal_errors.ErrUserDoesntExists
	}
	return nil
}

func (urc *UserRepositoryController) FollowUser(ownerID uint64, followerID uint64) error {
	err := urc.db.QueryRow(FollowUser, ownerID, followerID).Err()
	if err != nil {
		return internal_errors.ErrUserDoesntExists
	}
	return nil
}

func (urc *UserRepositoryController) UnfollowUser(ownerID uint64, followerID uint64) error {
	err := urc.db.QueryRow(UnfollowUser, ownerID, followerID).Err()
	if err != nil {
		return internal_errors.ErrUserDoesntExists
	}
	return nil
}

func (urc *UserRepositoryController) GetAllFollowings(ownerID uint64, followerID uint64) ([]uint64, error) {
	rows, err := urc.db.Query(GetAllFollowings, ownerID, followerID)
	if err != nil {
		return nil, fmt.Errorf("getAllFollowings: %w", err)
	}
	defer rows.Close()

	var followersList []uint64
	for rows.Next() {
		var followerID uint64
		if err := rows.Scan(&followerID); err != nil {
			return nil, fmt.Errorf("getAllFollowings rows.Next: %w", err)
		}
		followersList = append(followersList, followerID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getAllFollowings rows.Err: %w", err)
	}
	return followersList, nil
}

func (urc *UserRepositoryController) GetAllSubscriptions(ownerID uint64, followerID uint64) ([]uint64, error) {
	rows, err := urc.db.Query(GetAllSubscriptions, ownerID, followerID)
	if err != nil {
		return nil, fmt.Errorf("getAllFollowings: %w", err)
	}
	defer rows.Close()

	var followersList []uint64
	for rows.Next() {
		var followerID uint64
		if err := rows.Scan(&followerID); err != nil {
			return nil, fmt.Errorf("getAllFollowings rows.Next: %w", err)
		}
		followersList = append(followersList, followerID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getAllFollowings rows.Err: %w", err)
	}
	return followersList, nil
}

func (urc *UserRepositoryController) UserHasActiveSession(token string) bool {
	return urc.sm.Exists(token)
}

func (urc *UserRepositoryController) Session() *session.SessionsManager {
	return urc.sm
}
