package errors

import (
	"fmt"
	"net/http"
)

// Models validation
var (
	ErrorUserDataInvalid    = fmt.Errorf("user data is invalid")
	ErrorSectionDataInvalid = fmt.Errorf("section data is invalid")
	ErrorPinDataInvalid     = fmt.Errorf("pin data is invalid")
	ErrorCommentDataInvalid = fmt.Errorf("comment data is invalid")
	ErrorBoardDataInvalid   = fmt.Errorf("board data is invalid")
)

// Handlers
var (
	ErrorUserAlreadyRegistered = fmt.Errorf("user already registered")
)

// Errors to HTTP Status Code mapping
var ErrorToHttpStatusCode = map[error]int{
	// Models
	ErrorUserDataInvalid:    http.StatusUnprocessableEntity,
	ErrorSectionDataInvalid: http.StatusUnprocessableEntity,
	ErrorPinDataInvalid:     http.StatusUnprocessableEntity,
	ErrorCommentDataInvalid: http.StatusUnprocessableEntity,
	ErrorBoardDataInvalid:   http.StatusUnprocessableEntity,

	// Handlers
	ErrorUserAlreadyRegistered: http.StatusBadRequest,
}
