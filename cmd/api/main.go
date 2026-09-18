package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/mdsharifulislam-r/go-backend-template/config"
	"github.com/mdsharifulislam-r/go-backend-template/internal/database"
	"github.com/mdsharifulislam-r/go-backend-template/internal/server"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	app := server.New(cfg, db)
	if err := app.Start(); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
