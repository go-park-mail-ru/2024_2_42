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
	"strings"

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

func (urc *UserRepositoryController) GetUserIDWithEmail(email string) (uint64, error) {
	var userID uint64
	err := urc.db.QueryRow(GetUserIDByEmail, email).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 1, nil
		}
		return 0, fmt.Errorf("getLastUserID: %w", err)
	}
	return userID, nil
}

func (urc *UserRepositoryController) CreateUser(user *models.User) (uint64, error) {
	var userID uint64
	var nickName string
	err := urc.db.QueryRow(CreateUser, user.NickName, user.Email, user.Password).Scan(&userID, &nickName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			userID = 0
		}
		return 0, fmt.Errorf("psql error with userID = %d, Nickname = %s. CreateUser: %w", userID, nickName, err)
	}

	if userID == 0 {
		return 0, internal_errors.ErrBadUserInputData
	}
	urc.logger.WithField("user was succesfully created with userID", userID).Info("createUser func")
	return userID, nil
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

func (urc *UserRepositoryController) GetUsersByParams(userParams *models.UserSearchParams) ([]*models.UserInfo, error) {
	fmt.Println("YESYESYES")
	queryString := `SELECT user_id, user_name, nick_name, avatar_url FROM "user" WHERE `
	conditions := make([]string, 0)
	params := make([]interface{}, 0)

	if userParams.NickName != nil {
		conditions = append(conditions, `(LOWER(nick_name) LIKE '%' || $1 || '%')`)
		params = append(params, strings.ToLower(*userParams.NickName))
	}
	if userParams.Email != nil {
		conditions = append(conditions, `(LOWER(email) LIKE '%' || $2 || '%')`)
		params = append(params, strings.ToLower(*userParams.Email))
	}
	if userParams.UserName != nil {
		conditions = append(conditions, `(LOWER(user_name) LIKE '%' || $3 || '%')`)
		params = append(params, strings.ToLower(*userParams.UserName))
	}
	if userParams.Gender != nil {
		conditions = append(conditions, `(LOWER(gender) LIKE $4)`)
		params = append(params, strings.ToLower(*userParams.Gender))
	}
	queryString += strings.Join(conditions, ` AND `)
	fmt.Println(queryString)

	rows, err := urc.db.Query(queryString, params...)
	if err != nil {
		return nil, fmt.Errorf("getUsersByParams: %w", err)
	}
	defer rows.Close()

	res := make([]*models.UserInfo, 0)

	for rows.Next() {
		var foundUser *models.UserInfo = &models.UserInfo{}
		err := rows.Scan(&foundUser.UserID, &foundUser.UserName, &foundUser.NickName, &foundUser.AvatarUrl)
		if err != nil {
			return nil, fmt.Errorf("getUsersByParams: rows.Next %w", err)
		}
		res = append(res, foundUser)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getUsersByParams: rows.Err %w", err)
	}

	return res, nil
}

func (urc *UserRepositoryController) GetUserInfo(user *models.User, currUserID uint64) (*models.UserProfile, error) {
	var userInfo *models.UserProfile = &models.UserProfile{}

	err := urc.db.QueryRow(GetUserInfoByID, user.UserID).Scan(
		&userInfo.UserName,
		&userInfo.NickName,
		&userInfo.Description,
		&userInfo.BirthTime,
		&userInfo.Gender,
		&userInfo.AvatarUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &models.UserProfile{}, internal_errors.ErrUserDoesntExists
		}
		return &models.UserProfile{}, fmt.Errorf("psql GetUserByID: %w", err)
	}
	return userInfo, nil
}

func (urc *UserRepositoryController) GetUserInfoPublic(userID uint64) (*response.UserProfileResponse, error) {
	var userInfo *response.UserProfileResponse = &response.UserProfileResponse{}

	err := urc.db.QueryRow(GetUserInfoByID, userID).Scan(
		&userInfo.UserName,
		&userInfo.NickName,
		&userInfo.Description,
		&userInfo.BirthTime,
		&userInfo.Gender,
		&userInfo.AvatarUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &response.UserProfileResponse{}, internal_errors.ErrUserDoesntExists
		}
		return &response.UserProfileResponse{}, fmt.Errorf("psql GetUserByID: %w", err)
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
	urc.logger.WithField("user was succesfully deleted with userID", userID).Info()
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
	var followingsCount int64
	err := urc.db.QueryRow(GetFollowingsCount, follower_id).Scan(&followingsCount)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uint64(followingsCount), nil
		}
		return 0, fmt.Errorf("psql GetFollowingsCount: %w", err)
	}

	return uint64(followingsCount), nil
}

func (urc *UserRepositoryController) GetSubsriptionsCount(ownder_id uint64) (uint64, error) {
	var subscriptionsCount int64
	err := urc.db.QueryRow(GetSubsriptionsCount, ownder_id).Scan(&subscriptionsCount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uint64(subscriptionsCount), nil
		}
		return 0, fmt.Errorf("psql GetlSubsriptionsCount: %w", err)
	}

	return uint64(subscriptionsCount), nil
}

func (urc *UserRepositoryController) GetUserAvatar(userID uint64) (string, error) {
	var userAvatar *string
	err := urc.db.QueryRow(GetUserAvatar, userID).Scan(&userAvatar)
	if err != nil {
		return "", fmt.Errorf("getUserAvatar repo: %w", err)
	}

	avatar_url := ""
	if userAvatar != nil {
		avatar_url = *userAvatar
	}

	return avatar_url, nil
}

func (urc *UserRepositoryController) UserHasActiveSession(token string) bool {
	return urc.sm.Exists(token)
}

func (urc *UserRepositoryController) Session() *session.SessionsManager {
	return urc.sm
}
