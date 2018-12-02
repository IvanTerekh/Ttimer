package db

import (
	"database/sql"
	// Importing mysql driver.
	_ "github.com/go-sql-driver/mysql"
	"log"
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
		os.Getenv("DB_USER")+":"+os.Getenv("DB_PASSWORD")+
			"@tcp("+os.Getenv("DB_IP")+":"+os.Getenv("DB_PORT")+")/ttimer")
	if err != nil {
		log.Fatal(err)
	}

	return db
}
