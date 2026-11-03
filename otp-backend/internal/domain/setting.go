package domain

type Setting struct {
	Key   string `gorm:"primaryKey;type:varchar(50)"`
	Value bool
}

type SettingRepository interface {
	GetStatus(key string) (bool, error)
	UpdateStatus(key string, value bool) error
}

type SettingUseCase interface {
	GetChannelStatus(channel string) (bool, error)
	UpdateChannelStatus(channel string, active bool) error
	GetAllStatuses() (map[string]bool, error)
}