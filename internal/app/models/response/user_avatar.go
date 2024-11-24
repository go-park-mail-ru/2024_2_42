package response

type UserAvatar struct {
	Message   string `json:"message"`
	UserID    uint64 `json:"user_id"`
	AvatarUrl string `json:"avatar_url"`
}
