package response

type SignUpResponse struct {
	UserId  uint64 `json:"user_id"`
	Message string `json:"message"`
}
