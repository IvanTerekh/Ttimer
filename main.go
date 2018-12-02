package main

import (
	_ "ttimer/app"
	_ "ttimer/db"
	"ttimer/server"
)

func main() {
	server.Start()
}
