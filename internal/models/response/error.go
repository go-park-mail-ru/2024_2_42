package response

type ErrorResponse struct {
	CodeStatus int   `json:"code_status"`
	Message    string `json:"message"`
}
