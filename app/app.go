package app

import (
	"encoding/gob"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"log"
	"ttimer/scrambles"
	"ttimer/db"
	"os"
)

var (
	Store     *sessions.FilesystemStore
	Scrambler *scrambles.ScrambleProvider
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}

	err = db.Test()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(os.Getenv("TTIMER_HOST") + "DB connection: OK")

	Store = sessions.NewFilesystemStore("", []byte("something-very-secret"))

	Scrambler = scrambles.NewScrambleProvider(nil)
	log.Println("Scramblers provider initialized")

	gob.Register(map[string]interface{}{})
}
