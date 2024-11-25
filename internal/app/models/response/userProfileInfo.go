package response

import "time"

type UserProfileResponse struct {
	UserName               string     `json:"user_name"`
	NickName               string     `json:"nick_name"`
	Description            *string    `json:"description"`
	BirthTime              *time.Time `json:"birth_time"`
	Gender                 *string    `json:"gender"`
	AvatarUrl              *string    `json:"avatar_url"`
	NumOfUserFollowings    uint64     `json:"number_of_user_followings"`
	NumOfUserSubscriptions uint64     `json:"number_of_user_subscriptions"`
}
