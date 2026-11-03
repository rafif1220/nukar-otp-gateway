package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rafif1220/otp-backend/internal/domain"
)

type OtpHandler struct {
	usecase domain.OtpUseCase
	settingUseCase domain.SettingUseCase
}

func NewOtpHandler(u domain.OtpUseCase, s domain.SettingUseCase) *OtpHandler {
	return &OtpHandler{
		usecase:        u,
		settingUseCase: s,
	}
}

func (h *OtpHandler) Route(app fiber.Router, authMiddleware fiber.Handler) {
	app.Post("/send-otp", authMiddleware, h.SendOtp)
	app.Get("/history", h.GetHistory)
	app.Get("/metrics", h.GetMetrics)
	app.Post("/send-message", h.SendMessageManual) // Endpoint tambahan untuk dashboard
	app.Get("/chart-data", h.GetChartData)

	// Endpoint tambahan untuk kontrol saklar dari Svelte Dashboard
	app.Get("/settings/engine", h.GetEngineStatus)
	app.Put("/settings/engine", h.UpdateEngineStatus)
}

func (h *OtpHandler) SendOtp(c *fiber.Ctx) error {
	type Request struct {
		Channel string `json:"channel"`
		To      string `json:"to"`
		OtpCode string `json:"otp_code"`
	}

	// 1. Parse body request DULU untuk mengetahui channel-nya
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	targetChannel := req.Channel
	if targetChannel == "" {
		targetChannel = "wa" // Default fallback
	}

	// 2. Cek Saklar berdasarkan channel
	isActive, err := h.settingUseCase.GetChannelStatus(targetChannel)
	if err == nil && !isActive {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error":   "OTP_ENGINE_DISABLED",
			"message": "Sistem OTP untuk jalur " + targetChannel + " sedang dimatikan",
		})
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

func (h *OtpHandler) GetEngineStatus(c *fiber.Ctx) error {
	statuses, err := h.settingUseCase.GetAllStatuses()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membaca status engine"})
	}
	// Mengembalikan object JSON: {"wa_active": true, "email_active": true}
	return c.JSON(statuses) 
}

func (h *OtpHandler) UpdateEngineStatus(c *fiber.Ctx) error {
	type Req struct {
		Channel string `json:"channel"` // isinya "wa" atau "email"
		Active  bool   `json:"active"`
	}
	var req Req
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := h.settingUseCase.UpdateChannelStatus(req.Channel, req.Active)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal update status engine"})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"channel": req.Channel,
		"active":  req.Active,
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

	// Panggil usecase untuk mengirim pesan manual
	err := h.usecase.SendManualMessage(req.Phone, req.Message)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengirim pesan"})
	}

	return c.JSON(fiber.Map{
		"status":  "success", 
		"message": "Pesan berhasil diantrikan ke " + req.Phone,
	})
}

func (h *OtpHandler) GetChartData(c *fiber.Ctx) error {
	// Ambil query parameter ?range=... (default "7d" kalau kosong)
	timeRange := c.Query("range", "7d") 
	
	data, err := h.usecase.GetChartData(timeRange)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal fetch data chart"})
	}
	
	return c.JSON(data)
}