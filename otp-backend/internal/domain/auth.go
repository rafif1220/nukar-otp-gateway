package domain

import "time"

type ApiKey struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Secret    string    `json:"secret" gorm:"type:char(36)"` // UUID panjangnya 36 karakter
	ApiKey    string    `json:"api_key"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
