// Package server implements a web server.
package server

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/ivanterekh/ttimer/api"
	"github.com/ivanterekh/ttimer/server/auth"
	"github.com/ivanterekh/ttimer/server/middleware"
	"log"
	"net/http"
	"os"
)

// Start runs a new server
func Start() {
	go http.ListenAndServe(":80", http.HandlerFunc(redirectHttp))

	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("static"))))
	r.PathPrefix("/.well-known/pki-validation/").Handler(http.StripPrefix("/.well-known/pki-validation/",
		http.FileServer(http.Dir("validation"))))
	r.HandleFunc("/", HomeHandler)

	r.HandleFunc("/login", auth.LoginHandler)
	r.HandleFunc("/callback", auth.CallbackHandler)
	r.HandleFunc("/logout", auth.LogoutHandler)

	r.HandleFunc("/api/sessions", middleware.IfAuthenticated(api.ProvideSessionsHandler)).Methods("GET")
	r.HandleFunc("/api/addsession", middleware.IfAuthenticated(api.AddSessionHandler)).Methods("POST")
	r.HandleFunc("/api/deletesessions", middleware.IfAuthenticated(api.DeleteSessionsHandler)).Methods("POST")

	r.HandleFunc("/api/results", middleware.IfAuthenticated(api.ProvideResultsHandler)).Methods("GET")
	r.HandleFunc("/api/deleteresults", middleware.IfAuthenticated(api.DeleteResultsHandler)).Methods("POST")
	r.HandleFunc("/api/saveresults", middleware.IfAuthenticated(api.SaveResultsHandler)).Methods("POST")

	r.HandleFunc("/api/events", api.EventsHandler).Methods("GET")
	r.HandleFunc("/api/scramble", api.ScrambleHndler).Methods("GET")

	r.HandleFunc("/api/userinfo", api.UserInfoHandler).Methods("GET")
	r.HandleFunc("/api/isauthenticated", api.IsAuthenticatedHandler).Methods("GET")

	var err error

	log.Println("Listening: " + os.Getenv("TTIMER_PORT"))
	if os.Getenv("PRODUCTION") == "TRUE" {
		err = http.ListenAndServeTLS(":"+os.Getenv("TTIMER_PORT"),
			"cert.pem", "key.pem", handlers.LoggingHandler(os.Stdout, r))
	} else {
		err = http.ListenAndServe(":"+os.Getenv("TTIMER_PORT"), handlers.LoggingHandler(os.Stdout, r))
	}
	if err != nil {
		log.Fatal(err)
	}
}

func redirectHttp(w http.ResponseWriter, req *http.Request) {
	// remove/add not default ports from req.Host
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	fmt.Printf("redirect to: %s", target)
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}
