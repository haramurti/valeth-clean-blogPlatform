package http

import (
	"valeth-clean-blogPlatform/internal/domain"
	"valeth-clean-blogPlatform/internal/middleware"
	"valeth-clean-blogPlatform/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type CommentHandler struct {
	commentUseCase domain.CommentUseCase
}

// Constructor sekalian daftarin rute
func NewCommentHandler(app *fiber.App, u domain.CommentUseCase) {
	handler := &CommentHandler{
		commentUseCase: u,
	}

	api := app.Group("/api")

	// 1. Ambil Komen (Public)
	api.Get("/comments/:post_id", handler.GetComments)

	// 2. Kirim Komen (Protected / Harus Login)
	protected := api.Group("/", middleware.AuthProtected)
	protected.Post("/comments", handler.Store)
}

// Handler: Kirim Komentar
func (h *CommentHandler) Store(c *fiber.Ctx) error {
	var comment domain.Comment

	// 1. Ambil data JSON
	if err := c.BodyParser(&comment); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Data error"})
	}

	// 2. Cek Login
	userID := h.getMyUserID(c)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"message": "Login dulu bro!"})
	}

	comment.UserID = uint(userID)

	// 3. Simpan
	if err := h.commentUseCase.Create(&comment); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	// Balikin data biar frontend bisa langsung nampilin tanpa refresh
	return c.JSON(fiber.Map{"message": "Success", "data": comment})
}

// Handler: Ambil Daftar Komentar
func (h *CommentHandler) GetComments(c *fiber.Ctx) error {
	postID, err := c.ParamsInt("post_id")

	comments, err := h.commentUseCase.GetByPostID(postID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(comments)
}

// Helper: Ambil ID dari Token
func (h *CommentHandler) getMyUserID(c *fiber.Ctx) int {
	tokenString := c.Cookies("jwt_token")
	if tokenString == "" {
		return 0
	}
	userID, err := utils.ParseToken(tokenString)
	if err != nil {
		return 0
	}
	return int(userID)
}
