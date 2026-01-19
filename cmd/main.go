package main

import (
	"log"
	"valeth-clean-blogPlatform/config"
	"valeth-clean-blogPlatform/internal/delivery/http"
	"valeth-clean-blogPlatform/internal/repository"
	"valeth-clean-blogPlatform/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Setup Engine HTML
	engine := html.New("./web/templates", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Use(cors.New())
	app.Static("/static", "./web/static")

	// 2. Load Config & DB
	godotenv.Load()
	db, _ := config.NewDatabase()

	userRepo := repository.NewPostgresUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo)

	// 2. PASANG HANDLER AUTH (BARU)
	http.NewAuthHandler(app, userUseCase)

	// 3. Wiring (Ritual)
	postRepo := repository.NewPostgresPostRepository(db)
	postUseCase := usecase.NewPostUseCase(postRepo)

	// 4. Panggil Handler (SEMUA RUTE DITANGANI DISINI)
	http.NewPostHandler(app, postUseCase)

	// 5. Jalanin
	log.Println("Server jalan di port 8080 bos!")
	app.Listen(":8080")
}
