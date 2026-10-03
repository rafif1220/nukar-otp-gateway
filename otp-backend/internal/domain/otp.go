package domain

import "time"

// OtpLog merepresentasikan tabel di database
type OtpLog struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	To        string `json:"to"`
	OtpCode   string `json:"otp_code"`
	Status    string `json:"status"`
	Response  string `json:"response"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Metrics struct {
	Sent   int64 `json:"sent"`
	Failed int64 `json:"failed"`
}

// OtpRepository adalah kontrak untuk interaksi ke Database (PostgreSQL)
type OtpRepository interface {
	Create(log *OtpLog) error
	Update(log *OtpLog) error
	GetHistory(limit int) ([]OtpLog, error)
	GetMetrics() (Metrics, error)
}

// OtpUseCase adalah kontrak untuk business logic (termasuk Goroutine)
type OtpUseCase interface {
	SendOtpAsync(phone, otpCode string) (uint, error)
	GetHistory() ([]OtpLog, error)
	GetMetrics() (Metrics, error)
}
