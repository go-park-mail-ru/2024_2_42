package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"pinset/mailer-service/mailer"
	"pinset/mailer-service/usecase"

	"github.com/sirupsen/logrus"
)

type UserRepositoryController struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewUserRepositoryController(db *sql.DB, logger *logrus.Logger) usecase.UserRepository {
	return &UserRepositoryController{
		db:     db,
		logger: logger,
	}
}

func (urc *UserRepositoryController) GetUserInfoPublic(userID uint64) (*mailer.UserProfileResponse, error) {
	var userInfo *mailer.UserProfileResponse = &mailer.UserProfileResponse{}

	err := urc.db.QueryRow(`SELECT user_name, nick_name, description, birth_time, gender, avatar_url FROM "user" WHERE user_id = $1 LIMIT 1;`, userID).Scan(
		&userInfo.UserName,
		&userInfo.NickName,
		&userInfo.Description,
		&userInfo.BirthTime,
		&userInfo.Gender,
		&userInfo.AvatarUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &mailer.UserProfileResponse{}, fmt.Errorf("user does not exist")
		}
		return &mailer.UserProfileResponse{}, fmt.Errorf("psql GetUserByID: %w", err)
	}
	return userInfo, nil
}
