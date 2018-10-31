package api

import (
	"net/http"
	"encoding/json"
	"ttimer/scrambles"
)

var EventsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	jsonStr, err := json.Marshal(scrambles.Events)
	if err != nil {
		handleError(err, w)
		return
	}

	w.Write(jsonStr)
})
