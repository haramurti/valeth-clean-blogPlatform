package middleware

import (
	"valeth-clean-blogPlatform/internal/utils" // Import file jwt.go kamu

	"github.com/gofiber/fiber/v2"
)

// Ini Satpamnya
func AuthProtected(c *fiber.Ctx) error {
	// 1. Ambil Token dari Cookie bernama "jwt_token"
	tokenString := c.Cookies("jwt_token")

	// Kalau gak ada token, tendang ke login
	if tokenString == "" {
		return c.Redirect("/login")
	}

	// 2. Buka Segel Token (Validasi)
	// Kita pake fungsi ParseToken yang udah kamu buat di utils/jwt.go
	userID, err := utils.ParseToken(tokenString)
	if err != nil {
		// Kalau token rusak/palsu, tendang ke login
		return c.Redirect("/login")
	}

	// 3. Simpan User ID asli di saku (Locals) biar bisa dipake di Handler selanjutnya
	c.Locals("user_id", int(userID))

	return c.Next()

}
