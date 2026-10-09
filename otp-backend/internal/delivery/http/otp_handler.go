package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rafif1220/otp-backend/internal/domain"
)

type OtpHandler struct {
	usecase domain.OtpUseCase
}

func NewOtpHandler(u domain.OtpUseCase) *OtpHandler {
	return &OtpHandler{usecase: u}
}

func (h *OtpHandler) Route(app fiber.Router, authMiddleware fiber.Handler) {
	app.Post("/send-otp", authMiddleware, h.SendOtp)
	app.Get("/history", h.GetHistory)
	app.Get("/metrics", h.GetMetrics)
	app.Post("/send-message", h.SendMessageManual) // Endpoint tambahan untuk dashboard
}

func (h *OtpHandler) SendOtp(c *fiber.Ctx) error {
	type Request struct {
		Channel string `json:"channel"`
		To      string `json:"to"`
		OtpCode string `json:"otp_code"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	// Panggil Usecase
	logID, err := h.usecase.SendOtpAsync(req.Channel, req.To, req.OtpCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memproses OTP"})
	}

	// Langsung return 200 OK ke NUKAR
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Permintaan OTP diterima dan sedang diproses",
		"log_id":  logID,
	})
}

func (h *OtpHandler) GetHistory(c *fiber.Ctx) error {
	logs, err := h.usecase.GetHistory()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data history"})
	}
	return c.JSON(fiber.Map{"data": logs})
}

func (h *OtpHandler) GetMetrics(c *fiber.Ctx) error {
	metrics, err := h.usecase.GetMetrics()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data metrik"})
	}
	return c.JSON(metrics)
}

// Mock handler untuk kirim pesan manual dari SvelteKit Dashboard
func (h *OtpHandler) SendMessageManual(c *fiber.Ctx) error {
	type Request struct {
		Phone   string `json:"phone"`
		Message string `json:"message"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Logic manual Barantum (disimulasikan sukses untuk saat ini)
	return c.JSON(fiber.Map{
		"status":  "success", 
		"message": "Pesan berhasil diantrikan ke " + req.Phone,
	})
}
