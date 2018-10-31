package server

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/gorilla/handlers"
	"os"
	"log"
	"ttimer/server/middleware"
	"ttimer/api"
)

func Start() {
	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("static"))))
	r.HandleFunc("/", HomeHandler)

	r.HandleFunc("/login", loginHandler)
	r.HandleFunc("/callback", callbackHandler)
	r.HandleFunc("/logout", logoutHandler)

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
	err := http.ListenAndServe(":" + os.Getenv("TTIMER_PORT"), handlers.LoggingHandler(os.Stdout, r), )
	if err != nil {
		log.Fatal(err)
	}
}