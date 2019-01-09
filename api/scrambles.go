package api

import (
	"github.com/ivanterekh/ttimer/app"
	"net/http"
)

// ScrambleHndler provides a scrambler for given event.
var ScrambleHndler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	event := r.URL.Query().Get("event")
	scramble := app.Scrambler.Get(event)
	w.Write([]byte(scramble))
})
