// Package auth implements user authentication through auth0.com.
package auth

import (
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"os"
)

var protocol string

func init() {
	if os.Getenv("PRODUCTION") == "TRUE" {
		protocol = "https"
	} else {
		protocol = "http"
	}
}

func handleError(err error, w http.ResponseWriter, r *http.Request) {
	deleteCookie("state", w)
	deleteCookie("auth-session", w)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	log.Println(err)
}

func deleteCookie(name string, w http.ResponseWriter) {
	options := sessions.Options{MaxAge: -1}
	cookie := sessions.NewCookie(name, "_", &options)
	http.SetCookie(w, cookie)
}
