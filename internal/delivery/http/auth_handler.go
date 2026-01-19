package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url" // Penting buat URL Encode
	"time"
	"valeth-clean-blogPlatform/config"
	"valeth-clean-blogPlatform/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	userUseCase domain.UserUseCase
}

func NewAuthHandler(app *fiber.App, u domain.UserUseCase) {
	handler := &AuthHandler{
		userUseCase: u,
	}

	auth := app.Group("/auth")

	// 1. Login ke Google
	auth.Get("/google", handler.LoginGoogle)

	// 2. Balik dari Google
	auth.Get("/google/callback", handler.GoogleCallback)

	// 3. Proses Akhir Register (Nangkep Form) - INI YANG BARU
	auth.Post("/register-final", handler.RegisterFinal)

	// 4. Logout
	auth.Get("/logout", handler.Logout)
}

// --- 1. LEMPAR USER KE GOOGLE ---
func (h *AuthHandler) LoginGoogle(c *fiber.Ctx) error {
	conf := config.SetupGoogleOAuth()
	// "state-random" nanti bisa diganti token unik biar lebih aman
	url := conf.AuthCodeURL("state-random")
	return c.Redirect(url)
}

// --- 2. SAMBUT USER DARI GOOGLE ---
func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	// A. Ambil Code
	code := c.Query("code")
	if code == "" {
		return c.Status(400).SendString("Gagal login: Tidak ada kode dari Google")
	}

	// B. Tukar Code jadi Token
	conf := config.SetupGoogleOAuth()
	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		return c.Status(500).SendString("Gagal tukar token: " + err.Error())
	}

	// C. Ambil Data User via API Google
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return c.Status(500).SendString("Gagal ambil data user: " + err.Error())
	}
	defer resp.Body.Close()

	userData, _ := io.ReadAll(resp.Body)

	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.Unmarshal(userData, &googleUser); err != nil {
		return c.Status(500).SendString("Gagal baca JSON Google")
	}

	// D. CEK DI DATABASE (Panggil Usecase)
	user, err := h.userUseCase.CheckGoogleLogin(googleUser.ID)

	if err != nil {
		// KASUS: USER BELUM ADA -> REDIRECT KE REGISTER
		if err.Error() == "USER_NOT_FOUND" {
			// Titip data Google di URL
			targetURL := fmt.Sprintf("/register?email=%s&name=%s&avatar=%s&google_id=%s",
				url.QueryEscape(googleUser.Email),
				url.QueryEscape(googleUser.Name),
				url.QueryEscape(googleUser.Picture),
				url.QueryEscape(googleUser.ID),
			)
			return c.Redirect(targetURL)
		}

		// KASUS: ERROR LAIN
		return c.Status(500).SendString("Error database: " + err.Error())
	}

	// E. KASUS: USER ADA -> LANGSUNG LOGIN
	c.Cookie(&fiber.Cookie{
		Name:     "user_id",
		Value:    fmt.Sprintf("%d", user.ID),
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
	})

	return c.Redirect("/")
}

// --- 3. FINALISASI REGISTER (NANGKEP FORM HTML) ---
func (h *AuthHandler) RegisterFinal(c *fiber.Ctx) error {
	// A. Siapin struct buat nangkep inputan form
	var form struct {
		GoogleID string `form:"google_id"`
		Email    string `form:"email"`
		Name     string `form:"name"`
		Avatar   string `form:"avatar"`
	}

	// B. Parsing data form
	if err := c.BodyParser(&form); err != nil {
		return c.Status(400).SendString("Gagal baca form: " + err.Error())
	}

	// C. Masukin ke Struct Domain
	newUser := domain.User{
		GoogleID: form.GoogleID,
		Email:    form.Email,
		Name:     form.Name,
		Avatar:   form.Avatar,
	}

	// D. Simpan ke Database
	if err := h.userUseCase.RegisterUser(&newUser); err != nil {
		return c.Status(500).SendString("Gagal simpan user: " + err.Error())
	}

	// E. Auto Login (Kasih Cookie)
	c.Cookie(&fiber.Cookie{
		Name:     "user_id",
		Value:    fmt.Sprintf("%d", newUser.ID), // ID baru dari DB
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
	})

	// F. Masuk ke Home
	return c.Redirect("/")
}

// --- 4. LOGOUT ---
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:    "user_id",
		Expires: time.Now().Add(-1 * time.Hour), // Expire masa lalu
	})
	return c.Redirect("/")
}
