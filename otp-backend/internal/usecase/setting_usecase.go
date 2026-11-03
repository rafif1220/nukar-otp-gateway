package usecase

import "github.com/rafif1220/otp-backend/internal/domain"

type settingUseCase struct {
	repo domain.SettingRepository
}

func NewSettingUseCase(repo domain.SettingRepository) domain.SettingUseCase {
	return &settingUseCase{repo}
}

func (u *settingUseCase) GetChannelStatus(channel string) (bool, error) {
	key := "engine_" + channel + "_active"
	return u.repo.GetStatus(key)
}

func (u *settingUseCase) UpdateChannelStatus(channel string, active bool) error {
	key := "engine_" + channel + "_active"
	return u.repo.UpdateStatus(key, active)
}

func (u *settingUseCase) GetAllStatuses() (map[string]bool, error) {
	wa, _ := u.repo.GetStatus("engine_wa_active")
	email, _ := u.repo.GetStatus("engine_email_active")
	return map[string]bool{"wa": wa, "email": email}, nil
}