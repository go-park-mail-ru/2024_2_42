package response

type SignUpResponse struct {
	UserID        uint64 `json:"user_id"`
	SessionCookie string `json:"session_cookie"`
}
