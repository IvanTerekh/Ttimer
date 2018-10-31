package app

import (
	"encoding/gob"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"log"
	"ttimer/scrambles"
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

	Store = sessions.NewFilesystemStore("", []byte(os.Getenv("SESSION_KEY")))

	Scrambler = scrambles.NewScrambleProvider(nil)
	log.Println("Scramblers provider initialized")

	gob.Register(map[string]interface{}{})
}
