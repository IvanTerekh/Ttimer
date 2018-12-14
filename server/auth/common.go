// Package auth implements user authentication through auth0.com.
package auth

import (
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

func handleError(err error, w http.ResponseWriter) {
	log.Panic(err)
	http.Error(w, "Server error", http.StatusInternalServerError)
}

func deleteCookie(name string, w http.ResponseWriter) {
	options := sessions.Options{MaxAge: -1}
	cookie := sessions.NewCookie(name, "_", &options)
	http.SetCookie(w, cookie)
}
