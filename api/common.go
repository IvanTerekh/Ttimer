package api

import (
	"net/http"
	"log"
	"time"
	"ttimer/model"
	"io/ioutil"
	"encoding/json"
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
