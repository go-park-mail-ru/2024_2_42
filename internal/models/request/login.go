package request

type LoginRequest struct {
	Email string `json:"user_email"`
	Password string `json:"password"`
}
