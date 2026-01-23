package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// KUNCI RAHASIA (Jangan sampe orang tau)
// Di dunia nyata, ini ditaruh di .env
var SecretKey = []byte("rahasia-negara-kuma-blog-2026")

// 1. FUNGSI BIKIN TOKEN (Dipake pas Login)
func GenerateToken(userID uint) (string, error) {
	// Isi data token (Claims)
	claims := jwt.MapClaims{
		"user_id": userID,                                // Simpan ID User
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // Kadaluarsa 24 jam
	}

	// Bikin Token pakai metode HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Tanda tangani token pakai Secret Key
	return token.SignedString(SecretKey)
}

// 2. FUNGSI BACA TOKEN (Dipake pas User buka halaman)
func ParseToken(tokenString string) (uint, error) {
	// Coba buka segelnya
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Pastikan metodenya bener HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metode sign ngaco")
		}
		return SecretKey, nil
	})

	if err != nil {
		return 0, err
	}

	// Kalau segel aman, ambil datanya
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Ambil user_id (Data JSON itu float64, harus dicasting)
		userID := uint(claims["user_id"].(float64))
		return userID, nil
	}

	return 0, fmt.Errorf("token tidak valid")
}
