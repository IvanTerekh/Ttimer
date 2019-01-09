package db

import (
	"github.com/ivanterekh/ttimer/model"
)

// SelectResults gets results from the database.
func SelectResults(session *model.Session) ([]model.Result, error) {
	var results []model.Result
	err := db.Select(&results, "SELECT results.centis, results.scramble, results.penalty, results.datetime FROM results INNER JOIN sessions ON results.sessionID = sessions.id WHERE sessions.userID = ? AND sessions.name = ? AND sessions.event = ?",
		session.UserID, session.Name, session.Event)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// InsertResults adds resuts to the database.
func InsertResults(reses *[]model.Result, session *model.Session) error {
	sessionID, err := getSessionID(session)
	if err != nil {
		return err
	}

	for _, res := range *reses {
		_, err := stmts.insertResults.Exec(
			sessionID,
			res.Centis,
			res.Scramble,
			res.Penalty,
			res.Datetime,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteResults deletes results from database.
func DeleteResults(reses *[]model.Result, session *model.Session) error {
	sessionID, err := getSessionID(session)
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(
		"DELETE FROM results WHERE results.sessionID = ? AND results.datetime = ?")
	defer stmt.Close()
	if err != nil {
		return err
	}
	for _, res := range *reses {
		_, err := stmt.Exec(
			sessionID,
			res.Datetime,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
