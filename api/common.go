// Package api implements api of the web application.
package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"ttimer/model"
)

func handleError(err error, w http.ResponseWriter) {
	//log.Panic(err)
	log.Println(err)
	http.Error(w, "Server Error at "+time.Now().String(), http.StatusInternalServerError)
}

func retriveResultsAndSessionFromBody(r *http.Request) (*model.Session, *[]model.Result, error) {

	type resultsContainer struct {
		Results []model.Result `json:"results"`
		Session model.Session  `json:"session"`
	}

	jsonRes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, nil, err
	}
	var container resultsContainer
	err = json.Unmarshal(jsonRes, &container)

	results := container.Results

	session := container.Session

	return &session, &results, err
}
