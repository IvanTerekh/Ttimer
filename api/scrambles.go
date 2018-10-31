package api

import (
	"net/http"
	"ttimer/app"
)

var ScrambleHndler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	event := r.URL.Query().Get("event")
	scramble := app.Scrambler.Get(event)
	w.Write([]byte(scramble))
})
