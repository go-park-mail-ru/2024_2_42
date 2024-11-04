package request

import "time"

type UpdateUserInfoRequest struct {
	UserName    string    `json:"user_name"`
	NickName    string    `json:"nick_name"`
	Description string    `json:"description"`
	BirthTime   time.Time `json:"birth_time"`
	Gender      string    `json:"gender"`
	AvatarUrl   string    `json:"avatar_url"`
}

func (uuir UpdateUserInfoRequest) Valid() bool {
	return len(uuir.NickName) > 0 && len(uuir.NickName) < 25 &&
		len(uuir.Description) < 500 && len(uuir.Gender) < 25 &&
		uuir.BirthTime.Before(time.Now())
}
