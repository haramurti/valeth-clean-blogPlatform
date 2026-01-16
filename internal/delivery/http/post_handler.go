package http

import (
	"valeth-clean-blogPlatform/internal/domain"

	"github.com/gofiber/fiber/v2"
)

// 1. Struct Handler
// Isinya cuma satu: Dia butuh "Koki" (Usecase) siapa yang bakal dia suruh-suruh.
type PostHandler struct {
	PostUsecase domain.PostUseCase
}

// 2. Constructor
// Ini buat "ngelamar" si pelayan kerja. Pas dia kerja, dia dikasih tau siapa kokinya.
func NewPostHandler(app *fiber.App, usecase domain.PostUseCase) {
	handler := &PostHandler{
		PostUsecase: usecase,
	}

	// Ini Routing-nya (Daftar Menu)
	// Kalau ada yang akses /posts method POST, panggil fungsi Create
	app.Post("/posts", handler.Create)
}

// 3. Fungsi Nganter Pesanan (Handler Function)
func (h *PostHandler) Create(c *fiber.Ctx) error {
	// A. Nerima Pesanan (Parsing Body)
	var post domain.Post
	if err := c.BodyParser(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Pesenannya nggak jelas nih (Bad Request)",
			"error":   err.Error(),
		})
	}

	// B. Teriak ke Dapur (Panggil Usecase)
	// "Eh Usecase, tolong simpenin/masakin data ini dong!"
	err := h.PostUsecase.Store(&post)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Dapur lagi kebakaran (Internal Server Error)",
			"error":   err.Error(),
		})
	}

	// C. Nganter Makanan (Response Sukses)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Siap! Postingan udah tayang ya kakak",
		"data":    post,
	})
}
