package auth

import (
	"github.com/gorilla/sessions"
	"github.com/ivanterekh/ttimer/app"
	"log"
	"net/http"
)

// IsAuthenticated checks if user did log in.
func IsAuthenticated(r *http.Request) bool {
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
