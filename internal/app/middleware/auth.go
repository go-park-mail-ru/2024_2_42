package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	delivery "pinset/internal/app/delivery/http"
	internal_errors "pinset/internal/errors"
)

type ctxUserIDKeyType string

const UserIdKey ctxUserIDKeyType = "user_id"

func requestWithUserContext(r *http.Request, uc delivery.UserUsecase) (*http.Request, error) {
	c, err := r.Cookie("session_token")
	if err != nil {
		return nil, err
	}

	sessionToken := c.Value

	fmt.Println("Checking is authorized")

	userId, err := uc.IsAuthorized(sessionToken)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, UserIdKey, uint64(userId))

	return r.WithContext(ctx), nil
}

func RequiredAuthorization(uc delivery.UserUsecase, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// requestId := r.Context().Value(RequestIDKey).(string)
		request, err := requestWithUserContext(r, uc)
		if err != nil {
			if _, ok := internal_errors.ErrorMapping[err]; ok {
				internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
					Internal: err,
				})
				return
			}
			if errors.Is(err, http.ErrNoCookie) {
				internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
					Internal: internal_errors.ErrUserIsNotAuthorized,
				})
				return
			}
			internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
				General: err, Internal: internal_errors.ErrBadRequest,
			})
			return
		}

		next.ServeHTTP(w, request)
	}
}

func NotRequiredAuthorization(uc delivery.UserUsecase, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// requestId := r.Context().Value(RequestIDKey).(string)
		_, err := requestWithUserContext(r, uc)
		if err != nil && !errors.Is(err, http.ErrNoCookie) {
			fmt.Println("NotRequiredAuthorization")
			if _, ok := internal_errors.ErrorMapping[err]; ok {
				internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
					Internal: err,
				})
			} else {
				internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
					General: err, Internal: internal_errors.ErrInternalServerError,
				})
			}
			return
		}

		next.ServeHTTP(w, r)
	}
}
