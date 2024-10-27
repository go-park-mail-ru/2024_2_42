package middleware

import (
	"net/http"
	"slices"
	"strings"
)

var (
	originsAllowedList = []string{"http://localhost:3000", "http://37.139.41.77:8079"}

	methodAllowedList = []string{"GET", "POST", "DELETE", "OPTIONS"}
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 'Preflight' режим. В запросе используется метод OPTIONS
		if isPreflight(r) {
			origin := r.Header.Get("Origin")
			method := r.Header.Get("Access-Control-Request-Method")
			if slices.Contains(originsAllowedList, origin) && slices.Contains(methodAllowedList, method) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(methodAllowedList, ", "))
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
				w.WriteHeader(http.StatusOK)
			}
		} else {
			// Обычный запрос
			origin := r.Header.Get("Origin")
			if slices.Contains(originsAllowedList, origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
			}
		}

		next.ServeHTTP(w, r)
	})
}

func isPreflight(r *http.Request) bool {
	return r.Method == "OPTIONS" &&
		r.Header.Get("Origin") != "" &&
		r.Header.Get("Access-Control-Request-Method") != ""
}
