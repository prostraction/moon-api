package main

import (
	"log"

	"moon/internal/server"
)

func main() {
	s := server.Server{}
	app := s.NewRouter()
	if err := app.Listen(":9998"); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
