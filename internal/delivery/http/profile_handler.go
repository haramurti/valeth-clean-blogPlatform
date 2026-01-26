package http

import (
	"fmt"
	"io"
	"path/filepath"
	"time"

	"valeth-clean-blogPlatform/internal/domain"
	"valeth-clean-blogPlatform/internal/infrastructure"
	"valeth-clean-blogPlatform/internal/middleware" // Pastikan path ini benar

	"github.com/gofiber/fiber/v2"
)

type ProfileHandler struct {
	userUseCase domain.UserUseCase
	storage     *infrastructure.SupabaseStorage
}

// Constructor: Disini kita daftarin route-nya
func NewProfileHandler(app *fiber.App, u domain.UserUseCase, s *infrastructure.SupabaseStorage) {
	handler := &ProfileHandler{
		userUseCase: u,
		storage:     s,
	}

	// Bikin Group Route /profile
	profile := app.Group("/profile")

	// PASANG SATPAM (Middleware)
	// Biar cuma user yang udah login yang bisa akses group ini
	profile.Use(middleware.AuthProtected)

	// Route: POST /profile/avatar
	profile.Post("/avatar", handler.UploadAvatar)
}

func (h *ProfileHandler) UploadAvatar(c *fiber.Ctx) error {
	// 1. Ambil File dari Form HTML (name="avatar")
	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Wajib upload gambar"})
	}

	// 2. Baca File-nya
	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membuka file"})
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membaca file"})
	}

	// 3. Ambil User ID dari "Saku" (Locals) - Disimpan oleh Middleware
	// Kita perlu casting ke int karena Locals nyimpennya interface{}
	userID := c.Locals("user_id").(int)

	// 4. Bikin Nama File Unik (biar gak bentrok)
	// Format: avatar-USERID-TIMESTAMP.ext
	// Contoh: avatar-12-1709882211.jpg
	ext := filepath.Ext(fileHeader.Filename)
	filename := fmt.Sprintf("avatar-%d-%d%s", userID, time.Now().Unix(), ext)

	// 5. UPLOAD KE SUPABASE
	publicURL, err := h.storage.UploadFile(fileBytes, filename, fileHeader.Header.Get("Content-Type"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal upload ke Supabase: " + err.Error()})
	}

	// 6. UPDATE DATABASE USER
	// Ambil data user lama dulu
	user, err := h.userUseCase.GetProfile(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	// Ganti avatarnya dengan URL baru dari Supabase
	user.Avatar = publicURL

	// Simpan ke DB
	if err := h.userUseCase.UpdateUser(&user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal update database"})
	}

	// 7. Berhasil! Balikin URL-nya ke Frontend
	return c.JSON(fiber.Map{
		"message": "Avatar berhasil diganti!",
		"url":     publicURL,
	})
}
