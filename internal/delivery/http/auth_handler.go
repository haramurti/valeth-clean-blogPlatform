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
	"valeth-clean-blogPlatform/internal/utils"

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

	// Struct buat nangkep data Google
	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"` // <--- Ini URL Foto Terbaru
	}

	if err := json.Unmarshal(userData, &googleUser); err != nil {
		return c.Status(500).SendString("Gagal baca JSON Google")
	}

	// D. CEK DI DATABASE
	user, err := h.userUseCase.CheckGoogleLogin(googleUser.ID)

	if err != nil {
		// KASUS: USER BELUM ADA -> REDIRECT KE REGISTER
		if err.Error() == "USER_NOT_FOUND" {
			targetURL := fmt.Sprintf("/register?email=%s&name=%s&avatar=%s&google_id=%s",
				url.QueryEscape(googleUser.Email),
				url.QueryEscape(googleUser.Name),
				url.QueryEscape(googleUser.Picture),
				url.QueryEscape(googleUser.ID),
			)
			return c.Redirect(targetURL)
		}
		return c.Status(500).SendString("Error database: " + err.Error())
	}

	// --- [PERBAIKAN DISINI] ---
	// E. SYNC AVATAR: Update foto database kalau beda sama Google
	if user.Avatar != googleUser.Picture {
		// Update data di memory
		user.Avatar = googleUser.Picture

		// Panggil fungsi Update di Usecase (Pastikan fungsi ini ada!)
		// Kalau error update, kita ignore aja (log doang), yang penting user tetep bisa login
		_ = h.userUseCase.UpdateUser(&user)
	}
	// --------------------------

	// ✅ CUKUP TULIS SATU BARIS INI AJA
	// Fungsi ini otomatis bikin token JWT DAN set cookie avatar buat kamu.
	return h.generateTokenAndLogin(c, user.ID, user.Avatar)
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

	// --- LOGIC TAMBAHAN (FIX POINTER) ---
	// Kita cek dulu: Kalau kosong "", kita paksa jadi nil biar Database seneng (NULL).
	var googleIDPtr *string
	if form.GoogleID != "" {
		googleIDPtr = &form.GoogleID
	}
	// ------------------------------------

	// C. Masukin ke Struct Domain
	newUser := domain.User{
		GoogleID: googleIDPtr, // ✅ Sekarang tipe datanya udah pas (*string)
		Email:    form.Email,
		Name:     form.Name,
		Avatar:   form.Avatar,
	}

	// D. Simpan ke Database
	if err := h.userUseCase.RegisterUser(&newUser); err != nil {
		return c.Status(500).SendString("Gagal simpan user: " + err.Error())
	}

	// ✅ Generate Token & Cookie
	return h.generateTokenAndLogin(c, newUser.ID, newUser.Avatar)
}

// --- 4. LOGOUT ---
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Hapus cookie JWT
	c.Cookie(&fiber.Cookie{
		Name:    "jwt_token",
		Expires: time.Now().Add(-1 * time.Hour),
	})
	c.Cookie(&fiber.Cookie{
		Name:    "avatar",
		Expires: time.Now().Add(-1 * time.Hour),
	})
	// Bersih-bersih cookie lama
	c.Cookie(&fiber.Cookie{
		Name:    "user_id",
		Expires: time.Now().Add(-1 * time.Hour),
	})

	return c.Redirect("/")
}

// ... import ...

///baru 23 Jumat
//register Manual baru kamis 22

// --- [BARU] PROSES REGISTER MANUAL ---
func (h *AuthHandler) ProcessRegisterManual(c *fiber.Ctx) error {
	var form struct {
		Name     string `form:"name"`
		Email    string `form:"email"`
		Password string `form:"password"`
	}

	if err := c.BodyParser(&form); err != nil {
		return c.Status(400).SendString("Input error")
	}

	newUser := domain.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}

	if err := h.userUseCase.RegisterManual(&newUser); err != nil {
		return c.Status(500).SendString("Gagal Register: " + err.Error())
	}

	// Sukses -> Lempar ke Login
	return c.Redirect("/login")
}

// --- [BARU] PROSES LOGIN MANUAL ---
func (h *AuthHandler) ProcessLoginManual(c *fiber.Ctx) error {
	var form struct {
		Email    string `form:"email"`
		Password string `form:"password"`
	}

	if err := c.BodyParser(&form); err != nil {
		return c.Status(400).SendString("Input error")
	}

	// Panggil Usecase
	user, err := h.userUseCase.LoginManual(form.Email, form.Password)
	if err != nil {
		return c.Status(401).SendString("Login Gagal: " + err.Error())
	}

	// Set Cookie ID
	// ✅ CUKUP TULIS SATU BARIS INI AJA
	// Fungsi ini otomatis bikin token JWT DAN set cookie avatar buat kamu.
	return h.generateTokenAndLogin(c, user.ID, user.Avatar)
}

// 👇 INI FUNGSI BARU (Belum ada sebelumnya)
func (h *AuthHandler) generateTokenAndLogin(c *fiber.Ctx, userID uint, avatar string) error {
	// 1. PANGGIL ALAT UTILS: Bikin token terenkripsi dari ID User
	token, err := utils.GenerateToken(userID)
	if err != nil {
		return c.Status(500).SendString("Gagal bikin token JWT")
	}

	// 2. SIMPAN DI COOKIE "jwt_token" (Bukan "user_id" lagi)
	c.Cookie(&fiber.Cookie{
		Name:     "jwt_token", // <--- Nama cookie berubah
		Value:    token,       // <--- Isinya kode acak panjang (eyJ...), bukan angka "7"
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true, // <--- Aman dari hack JavaScript
	})

	// 3. Cookie Avatar (Ini tetep sama, cuma buat tampilan)
	c.Cookie(&fiber.Cookie{
		Name:     "avatar",
		Value:    avatar,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
	})

	return c.Redirect("/")
}
