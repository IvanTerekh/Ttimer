package api

import (
	"net/http"
	"ttimer/server/middleware"
	"fmt"
	"encoding/json"
	"ttimer/app"
)

var IsAuthenticatedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	authenticated, err := middleware.IsAuthenticated(r)
	if err != nil {
		handleError(err, w)
		return
	}

	w.Write([]byte(fmt.Sprint(authenticated)))
})

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