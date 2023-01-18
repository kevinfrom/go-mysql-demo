package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	server := SetupServer()
	log.Print(server.Run(":8080"))
}
