package main

import (
	"log"
	"os"
	"valeth-clean-blogPlatform/config"

	// Alias buat http delivery biar gak konflik
	httpDelivery "valeth-clean-blogPlatform/internal/delivery/http"
	"valeth-clean-blogPlatform/internal/infrastructure"
	"valeth-clean-blogPlatform/internal/repository"
	"valeth-clean-blogPlatform/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Setup Config & DB
	godotenv.Load()
	db, err := config.NewDatabase()
	if err != nil {
		log.Fatal("Gagal konek database:", err)
	}

	// 2. Setup Engine HTML & Fiber
	engine := html.New("./web/templates", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Use(cors.New())
	app.Static("/static", "./web/static")

	// 3. Setup Supabase Storage
	supabaseStorage := infrastructure.NewSupabaseStorage()
	// supabaseStorage.InitializeBucket() // (Opsional, boleh dinyalain buat ngecek koneksi)

	// 4. Setup Dependency Injection (Repo & Usecase)
	userRepo := repository.NewPostgresUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo)

	postRepo := repository.NewPostgresPostRepository(db)
	postUseCase := usecase.NewPostUseCase(postRepo)

	// 5. Setup Handler (Delivery)

	// Auth Handler
	httpDelivery.NewAuthHandler(app, userUseCase)

	// Profile Handler
	httpDelivery.NewProfileHandler(app, userUseCase, supabaseStorage)

	// Post Handler
	// 👇 PERBAIKAN PENTING: Tambahkan 'userUseCase' di sini!
	// Urutannya harus sama dengan definisi di NewPostHandler
	httpDelivery.NewPostHandler(app, postUseCase, userUseCase, supabaseStorage)

	//handler comment
	commentRepo := repository.NewPostgresCommentRepository(db)
	commentUseCase := usecase.NewCommentUseCase(commentRepo)
	httpDelivery.NewCommentHandler(app, commentUseCase)

	// 6. Jalanin Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app.Listen(":" + port)
	log.Println("Server jalan di port:" + port)
}
