package request

type SignUPRequest struct {
	UserName string `json:"user_name"`
	NickName string `json:"nick_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
