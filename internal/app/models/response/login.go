package response

type LogInResponse struct {
	SessionCookie string `json:"session_cookie"`
}

type IsAuthResponse struct {
	UserID uint64 `json:"user_id"`
}
