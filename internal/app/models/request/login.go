package request

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"session_token"`
}

func (lr LoginRequest) Valid() bool {
	return len(lr.Email) > 0 && len(lr.Password) > 0
}
