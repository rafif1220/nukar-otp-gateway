package postgres

import (
	"github.com/rafif1220/otp-backend/internal/domain"
	"gorm.io/gorm"
)

type settingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) domain.SettingRepository {
	return &settingRepository{db}
}

func (r *settingRepository) GetStatus(key string) (bool, error) {
	var setting domain.Setting
	err := r.db.First(&setting, "key = ?", key).Error
	return setting.Value, err
}

func (r *settingRepository) UpdateStatus(key string, value bool) error {
	return r.db.Model(&domain.Setting{}).Where("key = ?", key).Update("value", value).Error
}