package db

import (
	// Importing mysql driver.
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
)

var db *sqlx.DB

var stmts = struct {
	insertResults   *sqlx.Stmt
	selectSessionID *sqlx.Stmt
}{}

func init() {
	var err error
	db, err = sqlx.Open("mysql",
		os.Getenv("DB_USER")+":"+os.Getenv("DB_PASSWORD")+
			"@tcp("+os.Getenv("DB_IP")+":"+os.Getenv("DB_PORT")+")/ttimer")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}
	log.Printf("DB connection: OK. Host: %s:%s", os.Getenv("DB_IP"), os.Getenv("DB_PORT"))

	err = prepareStatements()
	if err != nil {
		log.Panic(err)
	}
}

func prepareStatements() error {
	insertResults, err := db.Preparex(
		"INSERT INTO results(sessionID, centis, scramble, penalty, datetime) VALUES(?, ?, ?, ?, ?)")
	stmts.insertResults = insertResults
	if err != nil {
		return err
	}

	selectSessionID, err := db.Preparex(
		"SELECT sessions.id FROM sessions WHERE sessions.userID = ? AND sessions.name = ? AND sessions.event = ?")
	stmts.selectSessionID = selectSessionID
	return err
}
