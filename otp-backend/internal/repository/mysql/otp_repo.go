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
	return r.db.Save(log).Error
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
