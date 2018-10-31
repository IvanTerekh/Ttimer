package db

import (
	"database/sql"
	"log"
	"ttimer/model"
	_ "github.com/go-sql-driver/mysql"
	"errors"
	"os"
)

func init() {
	db := openDb()
	defer db.Close()
	err := db.Ping()
	if err != nil {
		log.Panic(err)
	}
	log.Printf("DB connection: OK. Host: %s:%s", os.Getenv("DB_IP"), os.Getenv("DB_PORT"))
}

func openDb() *sql.DB {
	db, err := sql.Open("mysql",
		os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD")+
			"@tcp("+ os.Getenv("DB_IP")+ ":"+ os.Getenv("DB_PORT")+ ")/ttimer")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func SelectResults(session *model.Session) ([]model.Result, error) {
	db := openDb()
	defer db.Close()

	rows, err := db.Query(
		"SELECT results.centis, results.scramble, results.penalty, results.datetime FROM results INNER JOIN sessions ON results.sessionId = sessions.id WHERE sessions.userId = ? AND sessions.name = ? AND sessions.event = ?",
		session.UserId, session.Name, session.Event)
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

func SelectSessions(userId string) ([]model.Session, error) {
	db := openDb()
	defer db.Close()

	rows, err := db.Query(
		"SELECT sessions.name, sessions.event FROM sessions WHERE sessions.userId = ?",
		userId)
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

func InsertResults(reses *[]model.Result, session *model.Session) error {
	db := openDb()
	defer db.Close()

	sessionId, err := getSessionId(db, session)
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(
		"INSERT INTO results(sessionId, centis, scramble, penalty, datetime) VALUES(?, ?, ?, ?, ?)")
	defer stmt.Close()
	if err != nil {
		return err
	}
	for _, res := range *reses {
		_, err := stmt.Exec(
			sessionId,
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

func InsertSession(session *model.Session) error {
	db := openDb()
	defer db.Close()

	_, err := db.Exec(
		"INSERT INTO sessions(userId, name, event) VALUES(?, ?, ?)",
		session.UserId,
		session.Name,
		session.Event,
	)
	if err != nil {
		return err
	}
	return nil
}

func getSessionId(db *sql.DB, session *model.Session) (int, error) {
	rows, err := db.Query(
		"SELECT sessions.id FROM sessions WHERE sessions.userId = ? AND sessions.name = ? AND sessions.event = ?",
		session.UserId, session.Name, session.Event)
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

func SessionExists(session *model.Session) (bool, error) {
	db := openDb()
	defer db.Close()

	rows, err := db.Query(
		"SELECT sessions.id FROM sessions WHERE sessions.userId = ? AND sessions.name = ? AND sessions.event = ?",
		session.UserId, session.Name, session.Event)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}

func DeleteSessions(sessions *[]model.Session) error {
	db := openDb()
	defer db.Close()

	stmtRes, err := db.Prepare(
		"DELETE FROM results WHERE results.sessionId = ?")
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
		sessionId, err := getSessionId(db, &session)
		if err != nil {
			return err
		}

		_, err = stmtRes.Exec(sessionId)
		if err != nil {
			return err
		}

		_, err = stmtSes.Exec(sessionId)
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteResults(reses *[]model.Result, session *model.Session) error {
	db := openDb()
	defer db.Close()

	sessionId, err := getSessionId(db, session)
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(
		"DELETE FROM results WHERE results.sessionId = ? AND results.datetime = ?")
	defer stmt.Close()
	if err != nil {
		return err
	}
	for _, res := range *reses {
		_, err := stmt.Exec(
			sessionId,
			res.Datetime,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
