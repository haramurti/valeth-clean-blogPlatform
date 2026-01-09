package main

import (
	"fmt"
	"log"
	"valeth-clean-blogPlatform/config"
	"valeth-clean-blogPlatform/internal/repository"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := config.NewDatabase()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	fmt.Println("Success connect to Supabase! 🚀")

	postRepo := repository.NewPostgresPostRepository(db)

	fmt.Println("Repository initialized successfully:", postRepo)
}
