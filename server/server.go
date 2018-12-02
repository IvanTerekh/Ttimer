// Package server implements a web server.
package server

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"ttimer/api"
	"ttimer/server/auth"
	"ttimer/server/middleware"
)

// Start runs a new server
func Start() {
	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("static"))))
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

	log.Println("Listening: " + os.Getenv("TTIMER_PORT"))
	err := http.ListenAndServe(":"+os.Getenv("TTIMER_PORT"), handlers.LoggingHandler(os.Stdout, r))
	if err != nil {
		log.Fatal(err)
	}
}
