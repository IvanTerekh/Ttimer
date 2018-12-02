package db

import (
	"database/sql"
	"ttimer/model"
)

// SelectResults gets results from the database.
func SelectResults(session *model.Session) ([]model.Result, error) {
	db := openDb()
	defer db.Close()

	rows, err := db.Query(
		"SELECT results.centis, results.scramble, results.penalty, results.datetime FROM results INNER JOIN sessions ON results.sessionID = sessions.id WHERE sessions.userID = ? AND sessions.name = ? AND sessions.event = ?",
		session.UserID, session.Name, session.Event)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	results, err := retriveResults(rows)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// InsertResults adds resuts to the database.
func InsertResults(reses *[]model.Result, session *model.Session) error {
	db := openDb()
	defer db.Close()

	sessionID, err := getSessionID(db, session)
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(
		"INSERT INTO results(sessionID, centis, scramble, penalty, datetime) VALUES(?, ?, ?, ?, ?)")
	defer stmt.Close()
	if err != nil {
		return err
	}
	for _, res := range *reses {
		_, err := stmt.Exec(
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

func retriveResults(rows *sql.Rows) ([]model.Result, error) {
	var (
		reses []model.Result
		r     model.Result
	)
	for rows.Next() {
		err := rows.Scan(
			&r.Centis,
			&r.Scramble,
			&r.Penalty,
			&r.Datetime,
		)
		if err != nil {
			return nil, err
		}
		reses = append(reses, r)
	}
	return reses, nil
}

// DeleteResults deletes results from database.
func DeleteResults(reses *[]model.Result, session *model.Session) error {
	db := openDb()
	defer db.Close()

	sessionID, err := getSessionID(db, session)
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
