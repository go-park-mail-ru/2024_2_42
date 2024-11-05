package userRepository

import (
	"database/sql"
	"errors"
	"fmt"
	"pinset/internal/app/models"
	"pinset/internal/app/models/response"
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

func (urc *UserRepositoryController) GetLastUserID() (uint64, error) {
	var userID uint64
	err := urc.db.QueryRow(GetLastUserID).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 1, nil
		}
		return 0, fmt.Errorf("getLastUserID: %w", err)
	}
	return userID, nil
}

func (urc *UserRepositoryController) CreateUser(user *models.User) error {
	var userID uint64
	err := urc.db.QueryRow(CreateUser, user.UserName, user.NickName, user.Email, user.Password).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			userID = 0
		}
		return fmt.Errorf("psql CreateUser: %w", err)
	}

	if userID == 0 {
		return internal_errors.ErrBadUserInputData
	}
	urc.logger.WithField("user was succesful created with userID", userID).Info("createUser func")
	return nil
}

func (urc *UserRepositoryController) CheckUserByEmail(user *models.User) (bool, error) {
	var foundEmail sql.NullInt64

	err := urc.db.QueryRow(CheckUserByEmail, user.Email).Scan(&foundEmail)
	if err != nil && err != sql.ErrNoRows {
		return foundEmail.Valid, internal_errors.ErrUserDoesntExists
	}

	return foundEmail.Valid, nil
}

func (urc *UserRepositoryController) CheckUserCredentials(user *models.User) error {
	var userPassword string
	err := urc.db.QueryRow(CheckUserCredentials, user.Email).Scan(&userPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return internal_errors.ErrUserDoesntExists
		}
		return fmt.Errorf("psql GetUserByID: %w", err)
	}

	if userPassword != user.Password {
		return internal_errors.ErrBadPassword
	}
	return nil
}

func (urc *UserRepositoryController) GetUserInfo(user *models.User) (response.UserProfileResponse, error) {
	var userInfo response.UserProfileResponse

	err := urc.db.QueryRow(GetUserInfoByID, user.UserID).Scan(&userInfo.UserName, &userInfo.NickName, &userInfo.Description, &userInfo.BirthTime, &userInfo.Gender, &userInfo.AvatarUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return response.UserProfileResponse{}, internal_errors.ErrUserDoesntExists
		}
		return response.UserProfileResponse{}, fmt.Errorf("psql GetUserByID: %w", err)
	}
	return userInfo, nil
}

func (urc *UserRepositoryController) UpdateUserInfo(user *models.User) error {
	var userID uint64
	err := urc.db.QueryRow(UpdateUserInfoByID,
		user.UserName,
		user.NickName,
		user.Description,
		user.BirthTime,
		user.Gender,
		user.AvatarUrl,
		user.UserID).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return internal_errors.ErrUserDoesntExists
		}
		return fmt.Errorf("psql UpdateUserInfo: %w", err)
	}

	urc.logger.WithField("userInfoUpdated with userID:", userID).Info()
	return nil
}

func (urc *UserRepositoryController) UpdateUserPassword(user *models.User) error {
	var userID uint64
	err := urc.db.QueryRow(UpdateUserPasswordByID, user.UserID).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return internal_errors.ErrUserDoesntExists
		}
		return fmt.Errorf("psql UpdateUserPassword: %w", err)
	}

	return nil
}

func (urc *UserRepositoryController) DeleteUserByID(userID uint64) error {
	_, err := urc.db.Exec(DeleteUserByID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return internal_errors.ErrUserDoesntExists
		}
		return fmt.Errorf("psql DeleteUserByID: %w", err)
	}
	urc.logger.WithField("user was succesfil deleted with userID", userID).Info()
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
	rows, err := urc.db.Query(GetAllFollowings, ownerID)
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

func (urc *UserRepositoryController) GetFollowingsCount(follower_id uint64) (uint64, error) {
	var followingsCount uint64
	err := urc.db.QueryRow(GetFollowingsCount, follower_id).Scan(&followingsCount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return followingsCount, nil
		}
		return 0, fmt.Errorf("psql GetFollowingsCount: %w", err)
	}

	return followingsCount, nil
}

func (urc *UserRepositoryController) GetSubsriptionsCount(ownder_id uint64) (uint64, error) {
	var subscriptionsCount uint64
	err := urc.db.QueryRow(GetSubsriptionsCount, ownder_id).Scan(&subscriptionsCount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return subscriptionsCount, nil
		}
		return 0, fmt.Errorf("psql GetlSubsriptionsCount: %w", err)
	}

	return subscriptionsCount, nil
}

func (urc *UserRepositoryController) UserHasActiveSession(token string) bool {
	return urc.sm.Exists(token)
}

func (urc *UserRepositoryController) Session() *session.SessionsManager {
	return urc.sm
}
