package response

type SignUpResponse struct {
	UserID  uint64 `json:"user_id"`
	Message string `json:"message"`
}
