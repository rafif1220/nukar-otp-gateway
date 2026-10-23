package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/rafif1220/otp-backend/internal/delivery/http"
	"github.com/rafif1220/otp-backend/internal/domain"
	"github.com/rafif1220/otp-backend/internal/middleware"
	postgresRepo "github.com/rafif1220/otp-backend/internal/repository/postgres"
	"github.com/rafif1220/otp-backend/internal/usecase"
)

func main() {
	// 1. Load konfigurasi dari file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Tidak menemukan file .env")
	}

	// 2. Inisialisasi Koneksi MySQL
	// Tarik variabel satu per satu dari file .env
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbSslMode := os.Getenv("DB_SSLMODE")

	// Rakit format DSN menggunakan fmt.Sprintf
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta", dbHost, dbUser, dbPass, dbName, dbPort, dbSslMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}
	log.Println("Berhasil terkoneksi ke database")

	// Auto Migrate DB Schema dari entity Domain
	db.AutoMigrate(&domain.OtpLog{}, &domain.ApiKey{})
	
	// 3. Dependency Injection (Jahit semua layer Clean Architecture)
	otpRepo := postgresRepo.NewOtpRepository(db)
    otpUseCase := usecase.NewOtpUseCase(otpRepo)
    otpHandler := http.NewOtpHandler(otpUseCase)

	// Inisialisasi Middleware Auth
	authMiddleware := middleware.RequireBasicAuth(db)

	// 4. Inisialisasi Fiber & Middleware
	app := fiber.New()
	app.Use(cors.New())

	// 5. Daftarkan Routing API
	api := app.Group("/api")

	// Oper authMiddleware ke dalam fungsi Route
        otpHandler.Route(api, authMiddleware)

	// 6. Jalankan Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}
	log.Printf("🔥 NUKAR OTP Gateway berjalan di port %s", port)
	log.Fatal(app.Listen(":" + port))
}
