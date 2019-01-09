package main

import (
	_ "github.com/ivanterekh/ttimer/app"
	_ "github.com/ivanterekh/ttimer/db"
	"github.com/ivanterekh/ttimer/server"
)

func main() {
	server.Start()
}
