package middleware

import (
	"fmt"
	"net/http"
	"os"
)

func Panic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Fprintln(os.Stdout, "Panic recovered!", r.URL.Path)
				http.Error(w, "Internal server error", 500)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
