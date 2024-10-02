package request

type LoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func (lr LoginRequest) Valid() bool {
	return len(lr.Email) > 0 && len(lr.Password) > 0
}