package api

import (
	"encoding/json"
	"fmt"
	"github.com/ivanterekh/ttimer/app"
	"github.com/ivanterekh/ttimer/server/auth"
	"net/http"
)

// IsAuthenticatedHandler tells client if user is authenticated.
var IsAuthenticatedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	authenticated := auth.IsAuthenticated(r)
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

func retriveUserID(r *http.Request) (string, error) {
	authSession, err := app.Store.Get(r, "auth-session")
	if err != nil {
		return "", err
	}
	profile := authSession.Values["profile"].(map[string]interface{})
	return profile["sub"].(string), nil
}
