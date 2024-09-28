package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"youpin/internal/models/response"
)

type ErrorInfo struct {
	General  error
	Internal error
}

// Models validation
var (
	ErrorUserDataInvalid    = errors.New("user data is invalid")
	ErrorSectionDataInvalid = errors.New("section data is invalid")
	ErrorPinDataInvalid     = errors.New("pin data is invalid")
	ErrorCommentDataInvalid = errors.New("comment data is invalid")
	ErrorBoardDataInvalid   = errors.New("board data is invalid")
)

// Handlers
var (
	ErrorInvalidOrMissingRequestBody = errors.New("request body is invalid")

	ErrorUserAlreadyRegistered = errors.New("user is already registered")

	ErrFeedNotAccessible = errors.New("can't load the feed")
)

var ErrorMapping = map[error]struct {
	HttpCode     int
	InternalCode int
}{
	// Models validation
	ErrorUserDataInvalid:    {HttpCode: 400, InternalCode: 1},
	ErrorSectionDataInvalid: {HttpCode: 400, InternalCode: 2},
	ErrorPinDataInvalid:     {HttpCode: 400, InternalCode: 3},
	ErrorCommentDataInvalid: {HttpCode: 400, InternalCode: 4},
	ErrorBoardDataInvalid:   {HttpCode: 400, InternalCode: 5},

	// Handlers
	ErrorUserAlreadyRegistered: {HttpCode: 403, InternalCode: 6},

	ErrorInvalidOrMissingRequestBody: {HttpCode: 400, InternalCode: 7},

	ErrFeedNotAccessible: {HttpCode: 500, InternalCode: 15},
}

func SendErrorResponse(w http.ResponseWriter, ei ErrorInfo) {
	fmt.Fprintf(os.Stderr, "error: %s\n", ei.General.Error())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(ErrorMapping[ei.Internal].HttpCode)

	erJSON, err := json.Marshal(response.ErrorResponse{
		CodeStatus: ErrorMapping[ei.Internal].InternalCode, Message: ei.Internal.Error(),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return
	}

	_, err = w.Write(erJSON)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
	}
}
