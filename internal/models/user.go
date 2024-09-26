package models

import "time"

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
	BirthTime    time.Time `json:"birth_date"`
	Gender       string    `json:"gender"`
	AvatarUrl    string    `json:"avatar_url"`
	Followers    []User    `json:"followers"`
	Following    []User    `json:"following"`
	Boards       []Board   `json:"boards"`
	SignedUpTime time.Time `json:"created_at"`
}

func (u User) DataValid() bool {
	return len(u.UserName) > minUserNameLength &&
		len(u.NickName) > minNickNameLength &&
		len(u.Password) > minPasswordLength
}

func (u User) AgeRestricted() bool {
	return yearsBetween(u.BirthTime, time.Now()) < restrictionAge
}

func yearsBetween(t1, t2 time.Time) int {
	if t1.Location() != t2.Location() {
		t2 = t2.In(t1.Location())
	}
	if t1.After(t2) {
		t1, t2 = t2, t1
	}
	y1, _, _ := t1.Date()
	y2, _, _ := t2.Date()

	return int(y2 - y1)
}
