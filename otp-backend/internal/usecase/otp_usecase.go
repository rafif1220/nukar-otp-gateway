package usecase

import (
	"bytes"
	"strings"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"

	"github.com/rafif1220/otp-backend/internal/domain"
)

type otpUseCase struct {
	repo domain.OtpRepository
}

func NewOtpUseCase(repo domain.OtpRepository) domain.OtpUseCase {
	return &otpUseCase{repo}
}

func (u *otpUseCase) SendOtpAsync(channel, to, otpCode string) (uint, error) {
    // Default ke WA jika channel kosong
    if channel == "" {
        channel = "wa"
    }
	// 1. Catat log ke DB dengan status PENDING
	logEntry := &domain.OtpLog{
		To:      to,
		OtpCode: otpCode,
		Status:  "PENDING",
	}
	
	err := u.repo.Create(logEntry)
	if err != nil {
		return 0, err
	}

	// 2. Lempar tugas ke Goroutine di background (Fire and Forget)
	// Routing Goroutine berdasarkan channel
	if channel == "email" {
		go u.processEmailOtp(logEntry.ID, to, otpCode)
	} else {
		go u.processBarantumOtp(logEntry.ID, to, otpCode)
	}

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
				"user_name": "User Tester",
				"number": phone,
				"variabel": map[string]string{
					"{{1}}": "(otp)" + otpCode,
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
		
		// Pastikan HTTP sukses DAN Barantum mengirimkan array error kosong "[]"
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && strings.Contains(responseBody, `"error":"[]"`) {
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

// Goroutine baru khusus untuk kirim Email
func (u *otpUseCase) processEmailOtp(logID uint, to string, otpCode string) {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpEmail := os.Getenv("SMTP_EMAIL")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	// Setup Autentikasi
	auth := smtp.PlainAuth("", smtpEmail, smtpPassword, smtpHost)

	// Merakit Header dan Body Email
	subject := "Kode OTP NUKAR Anda"
	body := "Kode OTP Anda adalah: " + otpCode + ".\n\nJangan berikan kode ini kepada siapapun."
	message := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	// Eksekusi kirim email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpEmail, []string{to}, message)

	var status, responseBody string

	if err != nil {
		status = "FAILED"
		responseBody = err.Error()
	} else {
		status = "SENT"
		responseBody = `{"status":200,"message":"Email successfully sent"}`
	}

	// Update database
	updateLog := &domain.OtpLog{
		ID:       logID,
		Status:   status,
		Response: responseBody,
	}
	
	if err := u.repo.Update(updateLog); err != nil {
		log.Println("Gagal update status OTP Email di DB:", err)
	}
}

// Goroutine baru khusus untuk kirim pesan manual dari dashboard
func (u *otpUseCase) SendManualMessage(phone, message string) error {
	// Jalankan di background agar dashboard SvelteKit tidak loading lama
	go func(p, m string) {
		url := os.Getenv("BARANTUM_URL")
		
		payload := map[string]interface{}{
			"company_uuid":  os.Getenv("BARANTUM_COMPANY_UUID"),
			"template_uuid": os.Getenv("BARANTUM_TEMPLATE_UUID"), // Gunakan template notifikasi umum jika ada
			"chat_bot_uuid": os.Getenv("BARANTUM_CHAT_BOT_UUID"),
			"contacts": []map[string]interface{}{
				{
					"user_name": "User NUKAR",
					"number":    p,
					"variabel": map[string]string{
						"{{1}}": m, // Isi pesan dinamis dari textarea SvelteKit
					},
				},
			},
		}

		jsonPayload, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		
		if err == nil {
			resp.Body.Close()
		}
	}(phone, message)

	return nil
}

func (u *otpUseCase) GetChartData(timeRange string) (domain.ChartResponse, error) {
	aggs, err := u.repo.GetChartData(timeRange)
	
	resp := domain.ChartResponse{
		Labels: []string{},
		Sent:   []int{},
		Failed: []int{},
	}
	
	if err != nil {
		return resp, err
	}

	// Pecah data baris menjadi kolom array
	for _, a := range aggs {
		resp.Labels = append(resp.Labels, a.Label)
		resp.Sent = append(resp.Sent, a.Sent)
		resp.Failed = append(resp.Failed, a.Failed)
	}

	return resp, nil
}