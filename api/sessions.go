package api

import (
	"net/http"
	"ttimer/db"
	"ttimer/model"
	"encoding/json"
	"io/ioutil"
)

var ProvideSessionsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	userId, err := retriveUserId(r)
	if err != nil {
		handleError(err, w)
		return
	}

	sessions, err := db.SelectSessions(userId)
	if err != nil {
		handleError(err, w)
		return
	}

	jsonString, err := encodeSessions(sessions)
	if err != nil {
		handleError(err, w)
		return
	}

	w.Write(jsonString)
})

func encodeSessions(sessions []model.Session) ([]byte, error) {
	type responseSession struct {
		Name  string `json:"name"`
		Event string `json:"event"`
	}
	responseArray := make([]responseSession, len(sessions), len(sessions))
	for i, session := range sessions {
		responseArray[i] = responseSession{Name: session.Name, Event: session.Event}
	}

	jsonString, err := json.Marshal(responseArray)
	if err != nil {
		return nil, err
	}

	return jsonString, nil
}

var DeleteSessionsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	sessions, err := retriveSessionsFromBody(r)
	if err != nil {
		handleError(err, w)
		return
	}

	err = db.DeleteSessions(sessions)
	if err != nil {
		handleError(err, w)
		return
	}

	w.Write([]byte{})
})

var AddSessionHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	session, err := retriveNewSession(r)
	if err != nil {
		handleError(err, w)
		return
	}

	err = db.InsertSession(session)
	if err != nil {
		handleError(err, w)
		return
	}

	w.Write([]byte(""))
})

func retriveSessionsFromBody(r *http.Request) (*[]model.Session, error) {

	type container struct {
		Sessions []model.Session `json:"sessions"`
	}

	jsonRes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var cont container
	err = json.Unmarshal(jsonRes, &cont)

	sessions := cont.Sessions

	userId, err := retriveUserId(r)
	if err != nil {
		return nil, err
	}

	for i := range sessions {
		sessions[i].UserId = userId
	}

	return &sessions, err
}

func retriveNewSession(r *http.Request) (*model.Session, error) {
	type sessionContainer struct {
		Name  string `json:"name"`
		Event string `json:"event"`
	}

	jsonStr, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var container sessionContainer
	err = json.Unmarshal(jsonStr, &container)

	userId, err := retriveUserId(r)
	if err != nil {
		return nil, err
	}

	session := model.Session{
		UserId: userId,
		Name:   container.Name,
		Event:  container.Event,
	}

	return &session, err
}

func retriveSessionFromUrl(r *http.Request) (*model.Session, error) {

	var session model.Session

	userId, err := retriveUserId(r)
	if err != nil {
		return nil, err
	}

	session.UserId = userId
	session.Name = r.URL.Query().Get("sessionname")
	session.Event = r.URL.Query().Get("event")
	return &session, nil
}
