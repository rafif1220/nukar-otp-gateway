package middleware

import (
	"encoding/base64"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/rafif1220/otp-backend/internal/domain"
)

func RequireBasicAuth(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ambil header Authorization
		authHeader := c.Get("Authorization")

		// 2. Validasi keberadaan header dan prefix "Basic "
		if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header missing or invalid",
			})
		}

		// 3. Ambil string base64-nya saja
		base64Token := strings.TrimPrefix(authHeader, "Basic ")

		// 4. Decode base64
		decodedBytes, err := base64.StdEncoding.DecodeString(base64Token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header missing or invalid",
			})
		}

		// 5. Pisahkan hasil decode yang berformat api_key:secret
		decodedString := string(decodedBytes)
		parts := strings.SplitN(decodedString, ":", 2)
		if len(parts) != 2 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header missing or invalid",
			})
		}

		apiKey := parts[0]
		secret := parts[1]

		// 6. Cek kecocokan di Database
		var keyRecord domain.ApiKey
		err = db.Where("api_key = ? AND secret = ?", apiKey, secret).First(&keyRecord).Error
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized!",
			})
		}

		// 7. Jika cocok, persilakan lanjut ke Handler berikutnya
		return c.Next()
	}
}
