package server

import (
	"context"
	"crypto/rand"
	_ "crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"net/url"
	"os"
	"ttimer/app"
)

func callbackHandler(w http.ResponseWriter, r *http.Request) {

	authDomain := os.Getenv("AUTH0_DOMAIN")

	conf := &oauth2.Config{
		ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		RedirectURL:  "http://" + os.Getenv("TTIMER_DOMAIN"),
		Scopes:       []string{"openid", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://" + authDomain + "/authorize",
			TokenURL: "https://" + authDomain + "/oauth/token",
		},
	}
	state := r.URL.Query().Get("state")
	session, err := app.Store.Get(r, "state")
	if err != nil {
		handleError(err, w)
		return
	}

	if state != session.Values["state"] {
		http.Error(w, "Invalid state parameter", http.StatusInternalServerError)
		return
	}

	code := r.URL.Query().Get("code")

	token, err := conf.Exchange(context.TODO(), code)
	if err != nil {
		handleError(err, w)
		return
	}

	// Getting now the userInfo
	client := conf.Client(context.TODO(), token)
	resp, err := client.Get("https://" + authDomain + "/userinfo")
	if err != nil {
		handleError(err, w)
		return
	}

	defer resp.Body.Close()

	var profile map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	session, err = app.Store.Get(r, "auth-session")
	if err != nil {
		log.Println(err)
		handleError(err, w)
		return
	}

	session.Values["id_token"] = token.Extra("id_token")
	session.Values["access_token"] = token.AccessToken
	session.Values["profile"] = profile
	err = session.Save(r, w)
	if err != nil {
		handleError(err, w)
		return
	}

	// Redirect to logged in page
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	domain := os.Getenv("AUTH0_DOMAIN")
	aud := os.Getenv("AUTH0_AUDIENCE")

	conf := &oauth2.Config{
		ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		RedirectURL:  "http://" + os.Getenv("TTIMER_DOMAIN") + "/callback",
		Scopes:       []string{"openid", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://" + domain + "/authorize",
			TokenURL: "https://" + domain + "/oauth/token",
		},
	}

	if aud == "" {
		aud = "https://" + domain + "/userinfo"
	}

	// Generate random state
	b := make([]byte, 32)
	rand.Read(b)
	state := base64.StdEncoding.EncodeToString(b)

	session, err := app.Store.Get(r, "state")
	if err != nil {
		deleteCookie("state", w)
		deleteCookie("auth-session", w)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	session.Values["state"] = state
	err = session.Save(r, w)
	if err != nil {
		handleError(err, w)
		return
	}

	audience := oauth2.SetAuthURLParam("audience", aud)
	url := conf.AuthCodeURL(state, audience)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {

	session, err := app.Store.Get(r, "auth-session")
	if err != nil {
		handleError(err, w)
		return
	}
	session.Options.MaxAge = -1
	app.Store.Save(r, w, session)

	domain := os.Getenv("AUTH0_DOMAIN")

	var Url *url.URL
	Url, err = url.Parse("https://" + domain)
	if err != nil {
		handleError(err, w)
		return
	}

	Url.Path += "/v2/logout"
	parameters := url.Values{}
	parameters.Add("returnTo", "http://"+os.Getenv("TTIMER_DOMAIN"))
	parameters.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))
	Url.RawQuery = parameters.Encode()

	http.Redirect(w, r, Url.String(), http.StatusTemporaryRedirect)
}

func handleError(err error, w http.ResponseWriter) {
	log.Panic(err)
	http.Error(w, "Server error", http.StatusInternalServerError)
}

func deleteCookie(name string, w http.ResponseWriter) {
	options := sessions.Options{MaxAge: -1}
	cookie := sessions.NewCookie(name, "_", &options)
	http.SetCookie(w, cookie)
}
