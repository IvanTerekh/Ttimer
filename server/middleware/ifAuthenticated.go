package middleware

import (
	"ttimer/app"
	"net/http"
	"log"
)

func IfAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, err := IsAuthenticated(r)
		if err != nil {
			log.Println(err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		if auth {
			next(w, r)
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	}
}

func IsAuthenticated(r *http.Request) (bool, error) {
	session, err := app.Store.Get(r, "auth-session")
	if err != nil {
		return false, err
	}

	_, ok := session.Values["profile"]
	return ok, nil
}
