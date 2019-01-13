package auth

import (
	"context"
	"encoding/json"
	"github.com/ivanterekh/ttimer/app"
	"golang.org/x/oauth2"
	"net/http"
	"os"
)

// CallbackHandler handles callback from auth provider.
func CallbackHandler(w http.ResponseWriter, r *http.Request) {

	authDomain := os.Getenv("AUTH0_DOMAIN")

	conf := &oauth2.Config{
		ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		RedirectURL:  protocol + "://" + os.Getenv("TTIMER_DOMAIN"),
		Scopes:       []string{"openid", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://" + authDomain + "/authorize",
			TokenURL: "https://" + authDomain + "/oauth/token",
		},
	}
	state := r.URL.Query().Get("state")
	session, err := app.Store.Get(r, "state")
	if err != nil {
		handleError(err, w, r)
		return
	}

	if state != session.Values["state"] {
		http.Error(w, "Invalid state parameter", http.StatusInternalServerError)
		return
	}

	code := r.URL.Query().Get("code")

	token, err := conf.Exchange(context.TODO(), code)
	if err != nil {
		handleError(err, w, r)
		return
	}

	// Getting now the userInfo
	client := conf.Client(context.TODO(), token)
	resp, err := client.Get("https://" + authDomain + "/userinfo")
	if err != nil {
		handleError(err, w, r)
		return
	}

	defer resp.Body.Close()

	var profile map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		handleError(err, w, r)
		return
	}

	session, err = app.Store.Get(r, "auth-session")
	if err != nil {
		handleError(err, w, r)
		return
	}

	session.Values["id_token"] = token.Extra("id_token")
	session.Values["access_token"] = token.AccessToken
	session.Values["profile"] = profile
	err = session.Save(r, w)
	if err != nil {
		handleError(err, w, r)
		return
	}

	// Redirect to logged in page
	http.Redirect(w, r, "/", http.StatusSeeOther)

}
