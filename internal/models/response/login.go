package response

type LogInResponse struct {
	SessionCookie string `json:"session-cookie"`
}

type IsAuthResponse struct {
	UserID float64 `json:"user_id"`
}
