package middleware

import (
	"fmt"
	"net/http"
	"os"

	internal_errors "pinset/internal/errors"
)

func Panic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Fprintln(os.Stdout, "Panic recovered!", r.URL.Path)
				internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
					Internal: internal_errors.ErrInternalServerError,
				})
			}
		}()

		next.ServeHTTP(w, r)
	})
}
