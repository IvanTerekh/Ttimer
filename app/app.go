// Package app contains global app variables and their initialization.
package app

import (
	"encoding/gob"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"log"
	"os"
	"ttimer/scrambles"
)

var (
	// Store is server store for auth session.
	Store *sessions.FilesystemStore
	// Scrambler stores and updates scrambles.
	Scrambler *scrambles.ScrambleProvider
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}

	Store = sessions.NewFilesystemStore("", []byte(os.Getenv("SESSION_KEY")))

	Scrambler = scrambles.NewScrambleProvider()
	log.Println("Scramblers provider initialized")

	gob.Register(map[string]interface{}{})
}
