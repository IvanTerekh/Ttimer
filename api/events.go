package api

import (
	"encoding/json"
	"net/http"
	"ttimer/scrambles"
)

// EventsHandler provides a list of all possible events.
var EventsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	jsonStr, err := json.Marshal(scrambles.Events)
	if err != nil {
		handleError(err, w)
		return
	}

	w.Write(jsonStr)
})
