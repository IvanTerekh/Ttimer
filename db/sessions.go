package db

import (
	"database/sql"
	"errors"
	"log"
	"ttimer/model"
)

// SelectSessions gets sessions from the database.
func SelectSessions(userID string) ([]model.Session, error) {
	db := openDb()
	defer db.Close()

	rows, err := db.Query(
		"SELECT sessions.name, sessions.event FROM sessions WHERE sessions.userID = ?",
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions, err := retriveSessions(rows)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func retriveSessions(rows *sql.Rows) ([]model.Session, error) {
	var (
		sessions []model.Session
		s        model.Session
	)
	for rows.Next() {
		err := rows.Scan(
			&s.Name,
			&s.Event,
		)
		if err != nil {
			log.Fatal(err)
		}
		sessions = append(sessions, s)
	}
	return sessions, nil
}

// InsertSession adds session to the database.
func InsertSession(session *model.Session) error {
	db := openDb()
	defer db.Close()

	_, err := db.Exec(
		"INSERT INTO sessions(userID, name, event) VALUES(?, ?, ?)",
		session.UserID,
		session.Name,
		session.Event,
	)
	if err != nil {
		return err
	}
	return nil
}

func getSessionID(db *sql.DB, session *model.Session) (int, error) {
	rows, err := db.Query(
		"SELECT sessions.id FROM sessions WHERE sessions.userID = ? AND sessions.name = ? AND sessions.event = ?",
		session.UserID, session.Name, session.Event)
	if err != nil {
		return 0, err
	}
	if !rows.Next() {
		return 0, errors.New("db: Session not found")
	}
	var id int
	rows.Scan(&id)
	return id, nil
}

// SessionExists checks if session is already added to the database.
func SessionExists(session *model.Session) (bool, error) {
	db := openDb()
	defer db.Close()

	rows, err := db.Query(
		"SELECT sessions.id FROM sessions WHERE sessions.userID = ? AND sessions.name = ? AND sessions.event = ?",
		session.UserID, session.Name, session.Event)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}

// DeleteSessions deletes session and all it's results.
func DeleteSessions(sessions *[]model.Session) error {
	db := openDb()
	defer db.Close()

	stmtRes, err := db.Prepare(
		"DELETE FROM results WHERE results.sessionID = ?")
	defer stmtRes.Close()
	if err != nil {
		return err
	}

	stmtSes, err := db.Prepare(
		"DELETE FROM sessions WHERE sessions.id = ?")
	defer stmtRes.Close()
	if err != nil {
		return err
	}

	for _, session := range *sessions {
		sessionID, err := getSessionID(db, &session)
		if err != nil {
			return err
		}

		_, err = stmtRes.Exec(sessionID)
		if err != nil {
			return err
		}

		_, err = stmtSes.Exec(sessionID)
		if err != nil {
			return err
		}
	}

	return nil
}
