package main

import (
	"golang_restapi/internal/server"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	appPort := os.Getenv("APP_PORT")
	handler := server.InitializeServer()

	err := handler.Start(":" + appPort)
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
