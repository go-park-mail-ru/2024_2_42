package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"pinset/mailer-service/mailer"
	"pinset/mailer-service/usecase"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
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

	var userName, description, avatarUrl, gender *string
	var createdAt *time.Time
	err := urc.db.QueryRow(`SELECT user_name, nick_name, description, birth_time, gender, avatar_url FROM "user" WHERE user_id = $1 LIMIT 1;`, userID).Scan(
		&userName,
		&userInfo.NickName,
		&description,
		&createdAt,
		&gender,
		&avatarUrl)
	if createdAt != nil {
		userInfo.BirthTime = timestamppb.New(*createdAt)
	}
	if userName != nil {
		userInfo.UserName = &wrapperspb.StringValue{Value: *userName}
	}
	if description != nil {
		userInfo.Description = &wrapperspb.StringValue{Value: *description}
	}
	if gender != nil {
		userInfo.Gender = &wrapperspb.StringValue{Value: *gender}
	}
	if avatarUrl != nil {
		userInfo.AvatarUrl = &wrapperspb.StringValue{Value: *avatarUrl}
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &mailer.UserProfileResponse{}, fmt.Errorf("user does not exist")
		}
		return &mailer.UserProfileResponse{}, fmt.Errorf("psql GetUserByID: %w", err)
	}
	return userInfo, nil
}
