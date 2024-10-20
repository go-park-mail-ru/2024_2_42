package middleware

import (
	"context"
	"errors"
	"net/http"
)

const UserIdKey = "user_id"

func SessionExists(r *http.Request, a auth.AuthorizationClient) (*http.Request, error) {
	c, err := r.Cookie("session_token")
	if err != nil {
		return nil, err
	}
	sessionToken := c.Value
	res, err := a.CheckSession(r.Context(), &auth.CheckSessionRequest{Session: sessionToken})
	if err != nil {
		return nil, err
	}
	if !res.Valid {
		return nil, errs.GetLocalErrorByCode[res.LocalError]
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, UserIdKey, entity.UserID(res.UserId))
	return r.WithContext(ctx), nil
}

func AuthRequired(a auth.AuthorizationClient, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-CSRF-Token", csrf.Token(r))
		w.Header().Set("Access-Control-Expose-Headers", "X-CSRF-Token")
		requestId := r.Context().Value(RequestIDKey).(string)
		request, err := SessionExists(r, a)
		if err != nil {
			if errs.ErrorCodes[err].HttpCode != 0 {
				handler.WriteErrorResponse(w, l, requestId, errs.ErrorInfo{LocalErr: err})
				return
			}
			if errors.Is(err, http.ErrNoCookie) {
				handler.WriteErrorResponse(w, l, requestId, errs.ErrorInfo{LocalErr: errs.ErrUnauthorized})
				return
			}
			handler.WriteErrorResponse(w, l, requestId, errs.ErrorInfo{GeneralErr: err, LocalErr: errs.ErrReadCookie})
			return
		}
		next.ServeHTTP(w, request)
	}
}

func NoAuthRequired(a auth.AuthorizationClient, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestId := r.Context().Value(RequestIDKey).(string)
		_, err := SessionExists(r, a)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		handler.WriteErrorResponse(w, l, requestId, errs.ErrorInfo{
			LocalErr: errs.ErrAlreadyAuthorized,
		})
	}
}

func CheckAuth(a auth.AuthorizationClient, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := SessionExists(r, a)
		if err != nil {
			request = r
		}
		next.ServeHTTP(w, request)
	}
}
