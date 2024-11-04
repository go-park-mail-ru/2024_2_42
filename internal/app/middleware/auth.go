package middleware

import (
	"context"
	"errors"
	"net/http"
	delivery "pinset/internal/app/delivery/http"
	"pinset/internal/app/session"
	internal_errors "pinset/internal/errors"

	"github.com/sirupsen/logrus"
)

type ctxUserIDKeyType string

const UserIdKey ctxUserIDKeyType = "user_id"

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
	ctx = context.WithValue(ctx, UserIdKey, uint64(userId))

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
		// //_, err := requestWithUserContext(r, uc)
		// fmt.Println(err)
		// if err != nil && !errors.Is(err, http.ErrNoCookie) {
		// 	if _, ok := internal_errors.ErrorMapping[err]; ok {
		// 		internal_errors.SendErrorResponse(w, logger, internal_errors.ErrorInfo{
		// 			Internal: err,
		// 		})
		// 	} else {
		// 		internal_errors.SendErrorResponse(w, logger, internal_errors.ErrorInfo{
		// 			General: err, Internal: internal_errors.ErrInternalServerError,
		// 		})
		// 	}
		// 	return
		// }

		next.ServeHTTP(w, r)
	}
}
