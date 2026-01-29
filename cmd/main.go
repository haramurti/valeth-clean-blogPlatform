package main

import (
	"log"
	"valeth-clean-blogPlatform/config"

	// 👇 PERBAIKAN 1: Kasih alias "httpDelivery" disini
	httpDelivery "valeth-clean-blogPlatform/internal/delivery/http"

	// 👇 PERBAIKAN 2: Benerin typo "infrasturcture" jadi "infrastructure"
	"valeth-clean-blogPlatform/internal/infrastructure"

	"valeth-clean-blogPlatform/internal/repository"
	"valeth-clean-blogPlatform/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Setup Config & DB Dulu (Sebaiknya load env paling atas)
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
	supabaseStorage.InitializeBucket()

	// 4. Setup Dependency Injection (Repo & Usecase)
	userRepo := repository.NewPostgresUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo)

	postRepo := repository.NewPostgresPostRepository(db)
	postUseCase := usecase.NewPostUseCase(postRepo)

	// 5. Setup Handler (Delivery)
	// 👇 PERBAIKAN 3: Panggil semuanya pake alias "httpDelivery"

	// Auth Handler
	httpDelivery.NewAuthHandler(app, userUseCase)

	// Profile Handler (Fitur Upload Baru)
	httpDelivery.NewProfileHandler(app, userUseCase, supabaseStorage)

	// Post Handler
	httpDelivery.NewPostHandler(app, postUseCase, supabaseStorage)

	// 6. Jalanin Server
	log.Println("Server jalan di port 8080 bos!")
	app.Listen(":8080")
}
