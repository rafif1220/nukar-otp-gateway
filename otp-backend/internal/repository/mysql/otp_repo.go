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

func (r *otpRepository) GetChartData(timeRange string) ([]domain.ChartAggregation, error) {
	var data []domain.ChartAggregation
	var query string

	switch timeRange {
	case "today":
		// Group by Jam (contoh: 08:00, 14:00)
		query = `SELECT DATE_FORMAT(created_at, '%H:00') as label, 
				 SUM(CASE WHEN status = 'SENT' THEN 1 ELSE 0 END) as sent, 
				 SUM(CASE WHEN status = 'FAILED' THEN 1 ELSE 0 END) as failed 
				 FROM otp_logs WHERE DATE(created_at) = CURDATE() GROUP BY label ORDER BY label ASC`
	case "6m":
		// Group by Bulan (contoh: 2026-05, 2026-06)
		query = `SELECT DATE_FORMAT(created_at, '%Y-%m') as label, 
				 SUM(CASE WHEN status = 'SENT' THEN 1 ELSE 0 END) as sent, 
				 SUM(CASE WHEN status = 'FAILED' THEN 1 ELSE 0 END) as failed 
				 FROM otp_logs WHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 6 MONTH) GROUP BY label ORDER BY label ASC`
	case "1y":
		// Group by Bulan untuk 1 Tahun
		query = `SELECT DATE_FORMAT(created_at, '%Y-%m') as label, 
				 SUM(CASE WHEN status = 'SENT' THEN 1 ELSE 0 END) as sent, 
				 SUM(CASE WHEN status = 'FAILED' THEN 1 ELSE 0 END) as failed 
				 FROM otp_logs WHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 1 YEAR) GROUP BY label ORDER BY label ASC`
	case "1m":
		// Group by Tanggal (Hari per hari selama 1 bulan)
		query = `SELECT DATE_FORMAT(created_at, '%Y-%m-%d') as label, 
				 SUM(CASE WHEN status = 'SENT' THEN 1 ELSE 0 END) as sent, 
				 SUM(CASE WHEN status = 'FAILED' THEN 1 ELSE 0 END) as failed 
				 FROM otp_logs WHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 1 MONTH) GROUP BY label ORDER BY label ASC`
	default:
		// Default 7 Hari (7d)
		query = `SELECT DATE_FORMAT(created_at, '%Y-%m-%d') as label, 
				 SUM(CASE WHEN status = 'SENT' THEN 1 ELSE 0 END) as sent, 
				 SUM(CASE WHEN status = 'FAILED' THEN 1 ELSE 0 END) as failed 
				 FROM otp_logs WHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 7 DAY) GROUP BY label ORDER BY label ASC`
	}

	err := r.db.Raw(query).Scan(&data).Error
	return data, err
}