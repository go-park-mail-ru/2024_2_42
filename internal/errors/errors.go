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
	ErrUserDataInvalid    = errors.New("user data is invalid")
	ErrSectionDataInvalid = errors.New("section data is invalid")
	ErrPinDataInvalid     = errors.New("pin data is invalid")
	ErrCommentDataInvalid = errors.New("comment data is invalid")
	ErrBoardDataInvalid   = errors.New("board data is invalid")
)

// Handlers
var (
	ErrBadRequest = errors.New("bad request")

	ErrInvalidOrMissingRequestBody = errors.New("request body is invalid")
	ErrMethodIsNotAllowed          = errors.New("method is not allowed")
	ErrCantSignSessionToken        = errors.New("token signing failed")

	ErrUserAlreadyRegistered = errors.New("user is already registered")
	ErrUserAlreadyAuthorized = errors.New("user is already authorized")

	ErrUserIsNotRegistered = errors.New("user is not registered")
	ErrUserIsNotAuthorized = errors.New("user is not authorized")

	// JWT token
	ErrInvalidSessionToken = errors.New("session token is invalid")

	// Feed
	ErrFeedNotAccessible = errors.New("can't load the feed")
)

var ErrorMapping = map[error]struct {
	HttpCode     int
	InternalCode int
}{
	// Models validation
	ErrUserDataInvalid:    {HttpCode: 400, InternalCode: 1},
	ErrSectionDataInvalid: {HttpCode: 400, InternalCode: 2},
	ErrPinDataInvalid:     {HttpCode: 400, InternalCode: 3},
	ErrCommentDataInvalid: {HttpCode: 400, InternalCode: 4},
	ErrBoardDataInvalid:   {HttpCode: 400, InternalCode: 5},

	// Handlers
	ErrBadRequest: {HttpCode: 400, InternalCode: 6},

	ErrInvalidOrMissingRequestBody: {HttpCode: 400, InternalCode: 7},
	ErrMethodIsNotAllowed:          {HttpCode: 405, InternalCode: 8},
	ErrCantSignSessionToken:        {HttpCode: 500, InternalCode: 9},

	ErrUserAlreadyRegistered: {HttpCode: 403, InternalCode: 10},
	ErrUserAlreadyAuthorized: {HttpCode: 403, InternalCode: 11},

	ErrUserIsNotRegistered: {HttpCode: 403, InternalCode: 12},
	ErrUserIsNotAuthorized: {HttpCode: 401, InternalCode: 13},

	ErrInvalidSessionToken: {HttpCode: 403, InternalCode: 14},

	// Feed
	ErrFeedNotAccessible: {HttpCode: 500, InternalCode: 15},
}

func SendErrorResponse(w http.ResponseWriter, ei ErrorInfo) {
	if ei.General != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", ei.General.Error())
	}

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
