package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ttimer/app"
	"ttimer/server/middleware"
)

// IsAuthenticatedHandler tells client if user is authenticated.
var IsAuthenticatedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	authenticated := middleware.IsAuthenticated(w, r)
	w.Write([]byte(fmt.Sprint(authenticated)))
})

// UserInfoHandler provides information about user
var UserInfoHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	session, err := app.Store.Get(r, "auth-session")
	if err != nil {
		handleError(err, w)
		return
	}

	profile := session.Values["profile"]
	str, err := json.Marshal(profile)
	if err != nil {
		handleError(err, w)
		return
	}

	w.Write(str)
})

func retriveUserId(r *http.Request) (string, error) {
	authSession, err := app.Store.Get(r, "auth-session")
	if err != nil {
		return "", err
	}
	profile := authSession.Values["profile"].(map[string]interface{})
	return profile["sub"].(string), nil
}
