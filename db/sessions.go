package db

import (
	"errors"
	"github.com/ivanterekh/ttimer/model"
)

// SelectSessions gets sessions from the database.
func SelectSessions(userID string) ([]model.Session, error) {
	var sessions []model.Session

	err := db.Select(&sessions,
		"SELECT sessions.name, sessions.event FROM sessions WHERE sessions.userID = ?",
		userID)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

// InsertSession adds session to the database.
func InsertSession(session *model.Session) error {
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

func getSessionID(session *model.Session) (int, error) {
	rows, err := stmts.selectSessionID.Query(session.UserID, session.Name, session.Event)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	if !rows.Next() {
		return 0, errors.New("db: Session not found")
	}
	var id int
	rows.Scan(&id)
	return id, nil
}

// SessionExists checks if session is already added to the database.
func SessionExists(session *model.Session) (bool, error) {
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
		sessionID, err := getSessionID(&session)
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
