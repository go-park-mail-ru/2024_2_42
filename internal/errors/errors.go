package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"pinset/internal/app/models/response"
)

type ErrorInfo struct {
	General  error
	Internal error
}

// Models validation
var (
	ErrUserDataInvalid    = errors.New("данные пользователя невалидны")
	ErrSectionDataInvalid = errors.New("данные секции невалидны")
	ErrPinDataInvalid     = errors.New("данные пина невалидны")
	ErrCommentDataInvalid = errors.New("данные комментария невалидны")
	ErrBoardDataInvalid   = errors.New("данные доски невалидны")
)

// Handlers
var (
	ErrInternalServerError = errors.New("internal server error")
	ErrBadRequest          = errors.New("bad request")

	ErrInvalidOrMissingRequestBody = errors.New("тело запроса невалидно")
	ErrMethodIsNotAllowed          = errors.New("метод не допустим")
	ErrCantSignSessionToken        = errors.New("ошибка подписи сессионного токена")
	ErrCantProcessFormData         = errors.New("ошибка разбора данных логина")

	ErrUserAlreadyRegistered = errors.New("пользователь уже зарегистрирован")
	ErrUserAlreadyAuthorized = errors.New("пользователь уже авторизован")

	ErrUserIsNotRegistered = errors.New("такой пользователь не зарегистрирован")
	ErrUserIsNotAuthorized = errors.New("пользователь не авторизован")

	ErrDuringLogOutOperation = errors.New("ошибка при выходе из аккаунта")

	// JWT token
	ErrInvalidSessionToken = errors.New("сессионный токен невалиден")

	// Feed
	ErrFeedNotAccessible = errors.New("ошибка при загрузке ленты")

	// Media
	ErrExpectedMultipartContentType = errors.New("запрос имеет Content-Type не multipart")
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
	ErrInternalServerError: {HttpCode: 500, InternalCode: 6},
	ErrBadRequest:          {HttpCode: 400, InternalCode: 7},

	ErrInvalidOrMissingRequestBody: {HttpCode: 400, InternalCode: 8},
	ErrMethodIsNotAllowed:          {HttpCode: 400, InternalCode: 9},
	ErrCantSignSessionToken:        {HttpCode: 500, InternalCode: 10},
	ErrCantProcessFormData:         {HttpCode: 500, InternalCode: 11},

	ErrUserAlreadyRegistered: {HttpCode: 400, InternalCode: 12},
	ErrUserAlreadyAuthorized: {HttpCode: 400, InternalCode: 13},

	ErrUserIsNotRegistered: {HttpCode: 400, InternalCode: 14},
	ErrUserIsNotAuthorized: {HttpCode: 400, InternalCode: 15},

	ErrDuringLogOutOperation: {HttpCode: 500, InternalCode: 16},

	ErrInvalidSessionToken: {HttpCode: 400, InternalCode: 17},

	// Feed
	ErrFeedNotAccessible: {HttpCode: 500, InternalCode: 18},

	// Media
	ErrExpectedMultipartContentType: {HttpCode: 400, InternalCode: 19},
}

func SendErrorResponse(w http.ResponseWriter, ei ErrorInfo) {
	if ei.General != nil {
		fmt.Fprintf(os.Stdout, "error: %s\n", ei.General.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(ErrorMapping[ei.Internal].HttpCode)

	erJSON, err := json.Marshal(response.ErrorResponse{
		CodeStatus: ErrorMapping[ei.Internal].InternalCode, Message: ei.Internal.Error(),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		http.Error(w, "Internal server error", 500)
		return
	}

	_, err = w.Write(erJSON)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
	}
}
