package middleware

import (
	"context"
	"errors"
	"net/http"
	"pinset/configs"
	delivery "pinset/internal/app/delivery/http"
	"pinset/internal/app/session"
	internal_errors "pinset/internal/errors"

	"github.com/sirupsen/logrus"
)

func requestWithUserContext(r *http.Request, uc delivery.UserUsecase) (*http.Request, error) {
	c, err := r.Cookie(session.SessionTokenCookieKey)
	if err != nil {
		return nil, err
	}

	sessionToken := c.Value

	userId, err := uc.IsAuthorized(sessionToken)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, configs.UserIdKey, uint64(userId))
	return r.WithContext(ctx), nil
}

func RequiredAuthorization(logger *logrus.Logger, uc delivery.UserUsecase, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := requestWithUserContext(r, uc)
		if err != nil {
			if _, ok := internal_errors.ErrorMapping[err]; ok {
				internal_errors.SendErrorResponse(w, logger, internal_errors.ErrorInfo{
					Internal: err,
				})
				return
			}
			if errors.Is(err, http.ErrNoCookie) {
				internal_errors.SendErrorResponse(w, logger, internal_errors.ErrorInfo{
					Internal: internal_errors.ErrUserIsNotAuthorized,
				})
				return
			}
			internal_errors.SendErrorResponse(w, logger, internal_errors.ErrorInfo{
				General: err, Internal: internal_errors.ErrBadRequest,
			})
			return
		}

		next.ServeHTTP(w, request)
	}
}

func NotRequiredAuthorization(logger *logrus.Logger, uc delivery.UserUsecase, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := requestWithUserContext(r, uc)
		if err != nil && !errors.Is(err, http.ErrNoCookie) &&
			!errors.Is(err, internal_errors.ErrUserIsNotAuthorized) {
			if _, ok := internal_errors.ErrorMapping[err]; ok {
				internal_errors.SendErrorResponse(w, logger, internal_errors.ErrorInfo{
					Internal: err,
				})
			} else {
				internal_errors.SendErrorResponse(w, logger, internal_errors.ErrorInfo{
					General: err, Internal: internal_errors.ErrInternalServerError,
				})
			}
			return
		}
		if request == nil {
			request = r
		}
		next.ServeHTTP(w, request)
	}
}
