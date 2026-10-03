package usecase

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/rafif1220/otp-backend/internal/domain"
)

type otpUseCase struct {
	repo domain.OtpRepository
}

func NewOtpUseCase(repo domain.OtpRepository) domain.OtpUseCase {
	return &otpUseCase{repo}
}

func (u *otpUseCase) SendOtpAsync(phone, otpCode string) (uint, error) {
	// 1. Catat log ke DB dengan status PENDING
	logEntry := &domain.OtpLog{
		To:      phone,
		OtpCode: otpCode,
		Status:  "PENDING",
	}
	
	err := u.repo.Create(logEntry)
	if err != nil {
		return 0, err
	}

	// 2. Lempar tugas ke Goroutine di background (Fire and Forget)
	// Passing ID dari logEntry yang baru saja terbuat
	go u.processBarantumOtp(logEntry.ID, phone, otpCode)

	// 3. Langsung kembalikan ID ke Handler tanpa nunggu Goroutine selesai
	return logEntry.ID, nil
}

func (u *otpUseCase) GetHistory() ([]domain.OtpLog, error) {
	// Ambil 100 data terbaru
	return u.repo.GetHistory(100)
}

func (u *otpUseCase) GetMetrics() (domain.Metrics, error) {
	return u.repo.GetMetrics()
}

// Fungsi internal yang jalan di background (Goroutine)
func (u *otpUseCase) processBarantumOtp(logID uint, phone string, otpCode string) {
	url := os.Getenv("BARANTUM_URL")
	
	// Payload sesuai format Barantum
	payload := map[string]interface{}{
		"company_uuid":  os.Getenv("BARANTUM_COMPANY_UUID"),
		"template_uuid": os.Getenv("BARANTUM_TEMPLATE_UUID"),
		"chat_bot_uuid": os.Getenv("BARANTUM_CHAT_BOT_UUID"),
		"contacts": []map[string]interface{}{
			{
				"number": phone,
				"variabel": map[string]string{
					"{{1}}": otpCode,
				},
			},
		},
	}
	
	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	var status string
	var responseBody string

	// Handle error request atau parsing response
	if err != nil {
		status = "FAILED"
		responseBody = err.Error()
	} else {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)
		responseBody = string(bodyBytes)
		
		if resp.StatusCode == 200 || resp.StatusCode == 201 {
			status = "SENT"
		} else {
			status = "FAILED"
		}
	}

	// 4. Update status di database setelah Barantum merespon
	updateLog := &domain.OtpLog{
		ID:       logID,
		Status:   status,
		Response: responseBody,
	}
	
	if err := u.repo.Update(updateLog); err != nil {
		log.Println("Gagal update status OTP di DB:", err)
	}
}
