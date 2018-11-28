package middleware

import (
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"ttimer/app"
)

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
