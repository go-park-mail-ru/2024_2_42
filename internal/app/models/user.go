package models

import (
	"html"
	"net/mail"
	"time"

	"pinset/internal/errors"
	"pinset/pkg/utils"
)

const (
	minNickNameLength = 4
	minPasswordLength = 8
	restrictionAge    = 18
)

type User struct {
	UserID             uint64     `json:"user_id"`
	UserName           string     `json:"user_name"`
	NickName           string     `json:"nick_name"`
	Email              string     `json:"email"`
	Password           string     `json:"password"`
	Description        string     `json:"description"`
	BirthTime          *time.Time `json:"birth_date"`
	Gender             string     `json:"gender"`
	AvatarUrl          *string    `json:"avatar_url"`
	FollowingsCount    uint64     `json:"followings_count"`
	SubscriptionsCount uint64     `json:"subscriptions_count"`
	CreationTime       time.Time  `json:"creation_time"`
	UpdateTime         time.Time  `json:"update_time"`
}

type UserPin struct {
	UserID             uint64  `json:"user_id"`
	NickName           string  `json:"nick_name"`
	AvatarUrl          *string `json:"avatar_url"`
	FollowingsCount    uint64  `json:"followings_count"`
	SubscriptionsCount uint64  `json:"subscriptions_count"`
}

type UserProfile struct {
	UserName           *string    `json:"user_name"`
	NickName           string     `json:"nick_name"`
	UserBoards         []*Board   `json:"user_boards"`
	Description        *string    `json:"description"`
	BirthTime          *time.Time `json:"birth_date"`
	Gender             *string    `json:"gender"`
	AvatarUrl          *string    `json:"avatar_url"`
	FollowingsCount    uint64     `json:"followings_count"`
	SubscriptionsCount uint64     `json:"subscriptions_count"`
	CreationTime       time.Time  `json:"creation_time"`
	CurrentUser        bool       `json:"current_user"`
}

func NewUser(userID uint64, userName, email, password string) User {
	return User{
		UserID:   userID,
		UserName: userName,
		Email:    email,
		Password: password,
	}
}

func (u *User) Sanitize() {
	u.UserName = html.EscapeString(u.UserName)
	u.NickName = html.EscapeString(u.NickName)
	u.Email = html.EscapeString(u.Email)
	u.Password = html.EscapeString(u.Password)
	u.Gender = html.EscapeString(u.Email)
}

func (u User) Valid() error {
	if len(u.NickName) >= minNickNameLength &&
		len(u.Password) >= minPasswordLength &&
		u.BirthTime.Before(time.Now()) &&
		u.emailValid() {
		return nil
	}
	return errors.ErrUserDataInvalid
}

func (u User) emailValid() bool {
	_, err := mail.ParseAddress(u.Email)
	return err == nil
}

func (u User) AgeRestricted() bool {
	return utils.YearsBetween(*u.BirthTime, time.Now()) < restrictionAge
}
