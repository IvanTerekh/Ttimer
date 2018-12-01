package middleware

import (
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"ttimer/app"
)

// IfAuthenticated passes request to given handler only if user is authenticated.
func IfAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := IsAuthenticated(w, r)

		if auth {
			next(w, r)
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	}
}

// IsAuthenticated checks if user did log in.
func IsAuthenticated(w http.ResponseWriter, r *http.Request) bool {
	session, err := app.Store.Get(r, "auth-session")
	if err != nil {
		options := sessions.Options{MaxAge: -1}
		sessions.NewCookie("auth-session", "_", &options)
		log.Println(err)
		return false
	}

	_, ok := session.Values["profile"]
	return ok
}
