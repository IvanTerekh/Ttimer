package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ttimer/db"
	"ttimer/model"
)

// ProvideResultsHandler provides a list of results for given session.
var ProvideResultsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	session, err := retriveSessionFromUrl(r)
	if err != nil {
		handleError(err, w)
		return
	}

	oldSession, err := db.SessionExists(session)
	if err != nil {
		handleError(err, w)
		return
	}

	if oldSession {
		provideResults(session, w)
	} else {
		err = db.InsertSession(session)
		if err != nil {
			handleError(err, w)
			return
		}
		w.Write([]byte("[]"))
	}
})

// SaveResultsHandler saves given results into db.
var SaveResultsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	session, results, err := retriveResultsAndSessionFromBody(r)
	if err != nil {
		handleError(err, w)
		return
	}

	userId, err := retriveUserId(r)
	if err != nil {
		handleError(err, w)
		return
	}
	session.UserId = userId

	sessionExists, err := db.SessionExists(session)
	if err != nil {
		handleError(err, w)
		return
	}

	if !sessionExists {
		http.Error(w, fmt.Sprintf("Session %s with event %s doesn't exist", session.Name, session.Event), http.StatusBadRequest)
		return
	}

	err = db.InsertResults(results, session)
	if err != nil {
		handleError(err, w)
		return
	}
	w.Write([]byte(""))
})

// DeleteResultsHandler removes results from the db.
var DeleteResultsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	session, results, err := retriveResultsAndSessionFromBody(r)
	if err != nil {
		handleError(err, w)
		return
	}

	userId, err := retriveUserId(r)
	if err != nil {
		handleError(err, w)
		return
	}
	session.UserId = userId

	sessionExists, err := db.SessionExists(session)
	if err != nil {
		handleError(err, w)
		return
	}

	if !sessionExists {
		http.Error(w, fmt.Sprintf("Session %s with event %s doesn't exist", session.Name, session.Event), http.StatusBadRequest)
		return
	}

	err = db.DeleteResults(results, session)
	if err != nil {
		handleError(err, w)
		return
	}
	w.Write([]byte(""))
})

func provideResults(session *model.Session, w http.ResponseWriter) {
	res, err := db.SelectResults(session)
	jsonString, err := json.Marshal(res)
	if err != nil {
		handleError(err, w)
		return
	}

	w.Write([]byte(jsonString))
}
