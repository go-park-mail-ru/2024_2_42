package models

import (
	"html"
	"net/mail"
	"time"

	"pinset/internal/errors"
	"pinset/pkg/utils"
)

const (
	minUserNameLength = 4
	minNickNameLength = 4
	minPasswordLength = 8
	restrictionAge    = 18
)

type User struct {
	UserID       uint64    `json:"user_id"`
	UserName     string    `json:"user_name"`
	NickName     string    `json:"nick_name"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	Description  string    `json:"description"`
	BirthTime    time.Time `json:"birth_date"`
	Gender       string    `json:"gender"`
	AvatarUrl    string    `json:"avatar_url"`
	CreationTime time.Time `json:"creation_time"`
	UpdateTime   time.Time `json:"update_time"`
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
	u.AvatarUrl = html.EscapeString(u.AvatarUrl)
}

func (u User) Valid() error {
	if len(u.UserName) >= minUserNameLength &&
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
	return utils.YearsBetween(u.BirthTime, time.Now()) < restrictionAge
}
