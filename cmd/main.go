package main

import (
	"fmt"
	"log"
	"valeth-clean-blogPlatform/config"
	"valeth-clean-blogPlatform/internal/repository"

	"github.com/joho/godotenv"
)

func main() {
	// 1. Load Env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 2. Connect Database
	db, err := config.NewDatabase()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	fmt.Println("Success connect to Supabase! 🚀")

	postRepo := repository.NewPostgresPostRepository(db)

	// Biar 'postRepo' gak kena error "declared and not used" juga,
	// kita print aja dulu (cuma buat ngecek).
	fmt.Println("Repository initialized successfully:", postRepo)
}
