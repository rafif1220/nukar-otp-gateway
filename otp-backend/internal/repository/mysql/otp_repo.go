package mysql

import (
	"github.com/rafif1220/otp-backend/internal/domain"
	"gorm.io/gorm"
)

type otpRepository struct {
	db *gorm.DB
}

func NewOtpRepository(db *gorm.DB) domain.OtpRepository {
	return &otpRepository{db}
}

func (r *otpRepository) Create(log *domain.OtpLog) error {
	return r.db.Create(log).Error
}

func (r *otpRepository) Update(log *domain.OtpLog) error {
	// Gunakan map agar GORM HANYA meng-update kolom status dan response,
	// tanpa menyentuh created_at, to, atau otp_code sama sekali.
	return r.db.Model(&domain.OtpLog{}).Where("id = ?", log.ID).Updates(map[string]interface{}{
		"status":   log.Status,
		"response": log.Response,
	}).Error
}

func (r *otpRepository) GetHistory(limit int) ([]domain.OtpLog, error) {
	var logs []domain.OtpLog
	err := r.db.Order("created_at desc").Limit(limit).Find(&logs).Error
	return logs, err
}

func (r *otpRepository) GetMetrics() (domain.Metrics, error) {
	var metrics domain.Metrics
	r.db.Model(&domain.OtpLog{}).Where("status = ?", "SENT").Count(&metrics.Sent)
	r.db.Model(&domain.OtpLog{}).Where("status = ?", "FAILED").Count(&metrics.Failed)
	return metrics, nil
}
