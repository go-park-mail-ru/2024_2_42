package middleware

import (
	"net/http"

	internal_errors "pinset/internal/errors"

	"github.com/sirupsen/logrus"
)

func Panic(logger *logrus.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				internal_errors.SendErrorResponse(w, logger, internal_errors.ErrorInfo{
					Internal: internal_errors.ErrInternalServerError,
				})
			}
		}()

		next.ServeHTTP(w, r)
	})
}
