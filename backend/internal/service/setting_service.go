package service

import (
	"errors"

	"github.com/singbox-ui/singbox-ui/internal/database"
	"github.com/singbox-ui/singbox-ui/internal/database/models"
)

type SettingService struct{}

var SettingSvc = &SettingService{}

func (s *SettingService) GetAllSettings() (map[string]string, error) {
	var settings []models.Setting
	if err := database.DB.Find(&settings).Error; err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, item := range settings {
		result[item.Key] = item.Value
	}
	return result, nil
}

func (s *SettingService) UpdateSettings(kvs map[string]string) error {
	for k, v := range kvs {
		var setting models.Setting
		if err := database.DB.Where("key = ?", k).First(&setting).Error; err == nil {
			database.DB.Model(&setting).Update("value", v)
		} else {
			database.DB.Create(&models.Setting{Key: k, Value: v})
		}
	}
	return nil
}

func (s *SettingService) UpdateAdminPassword(username, oldPassword, newPassword string) error {
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return errors.New("user not found")
	}

	if !user.CheckPassword(oldPassword) {
		return errors.New("incorrect old password")
	}

	if err := user.SetPassword(newPassword); err != nil {
		return err
	}

	return database.DB.Save(&user).Error
}
