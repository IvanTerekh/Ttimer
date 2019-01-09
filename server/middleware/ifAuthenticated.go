package middleware

import (
	"github.com/ivanterekh/ttimer/server/auth"
	"net/http"
)

// IfAuthenticated passes request to given handler only if user is authenticated.
func IfAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isAuth := auth.IsAuthenticated(r)

		if isAuth {
			next(w, r)
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	}
}
