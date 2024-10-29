package errors

import (
	"encoding/json"
	"errors"
	"net/http"

	"pinset/internal/app/models/response"

	"github.com/sirupsen/logrus"
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
	ErrWrongMediaContentType = errors.New("загружаемое медиа имеет некорректный Content-Type")
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
	ErrWrongMediaContentType: {HttpCode: 400, InternalCode: 20},
}

func IsInternal(err error) bool {
	_, ok := ErrorMapping[err]
	return ok 
}

func SendErrorResponse(w http.ResponseWriter, logger *logrus.Logger, ei ErrorInfo) {
	var generalErrorText, localErrorText string

	if ei.General != nil {
		generalErrorText = ei.General.Error()
	}
	if ei.Internal != nil {
		localErrorText = ei.Internal.Error()
	}

	logger.WithFields(logrus.Fields{
		"general_error":    generalErrorText,
		"local_error":      localErrorText,
		"local_error_code": ErrorMapping[ei.Internal].InternalCode,
	}).Info("Error response")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(ErrorMapping[ei.Internal].HttpCode)

	erJSON, err := json.Marshal(response.ErrorResponse{
		CodeStatus: ErrorMapping[ei.Internal].InternalCode, Message: ei.Internal.Error(),
	})
	if err != nil {
		logger.Error("Unpredicted error during sending error response")
		http.Error(w, "Internal server error", 500)
		return
	}

	_, err = w.Write(erJSON)
	if err != nil {
		logger.Error("Unpredicted error during sending error response")
		http.Error(w, "Internal server error", 500)
	}
}
