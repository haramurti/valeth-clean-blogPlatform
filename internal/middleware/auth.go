package middleware

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Ini Satpamnya
func AuthProtected(c *fiber.Ctx) error {

	fmt.Println("👮 [SATPAM] Middleware Jalan di rute:", c.Path())

	// 1. CEK COOKIE
	cookieUserID := c.Cookies("user_id")

	// CCTV (Debugging)
	fmt.Printf("[SATPAM] Cek Rute: %s | Cookie: '%s'\n", c.Path(), cookieUserID)

	if cookieUserID == "" {
		fmt.Println("[SATPAM] 🚫 Gak ada tiket! Redirect ke /login")
		return c.Redirect("/login")
	}

	// 2. PARSING ID
	userID, err := strconv.Atoi(cookieUserID)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"message": "Invalid User ID"})
	}

	// 3. SIMPAN DATA DI SAKU
	c.Locals("user_id", userID)

	return c.Next()
}
