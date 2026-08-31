package service

import (
	"errors"

	"github.com/singbox-ui/singbox-ui/internal/database"
	"github.com/singbox-ui/singbox-ui/internal/database/models"
)

type OutboundService struct{}

var OutboundSvc = &OutboundService{}

func (s *OutboundService) ListOutbounds() ([]models.Outbound, error) {
	var outbounds []models.Outbound
	err := database.DB.Order("id asc").Find(&outbounds).Error
	return outbounds, err
}

func (s *OutboundService) CreateOutbound(out *models.Outbound) error {
	var count int64
	database.DB.Model(&models.Outbound{}).Where("tag = ?", out.Tag).Count(&count)
	if count > 0 {
		return errors.New("outbound tag already exists")
	}

	if err := database.DB.Create(out).Error; err != nil {
		return err
	}
	_ = InboundSvc.SyncCoreConfig()
	return nil
}

func (s *OutboundService) UpdateOutbound(out *models.Outbound) error {
	var existing models.Outbound
	if err := database.DB.First(&existing, out.ID).Error; err != nil {
		return err
	}

	existing.Tag = out.Tag
	existing.Type = out.Type
	existing.Server = out.Server
	existing.Port = out.Port
	existing.Settings = out.Settings
	existing.Enable = out.Enable
	existing.Remark = out.Remark

	if err := database.DB.Save(&existing).Error; err != nil {
		return err
	}
	_ = InboundSvc.SyncCoreConfig()
	return nil
}

func (s *OutboundService) DeleteOutbound(id uint) error {
	var out models.Outbound
	if err := database.DB.First(&out, id).Error; err != nil {
		return err
	}
	if out.Tag == "direct" || out.Tag == "block" || out.Tag == "dns-out" {
		return errors.New("cannot delete built-in system outbounds")
	}

	if err := database.DB.Delete(&models.Outbound{}, id).Error; err != nil {
		return err
	}
	_ = InboundSvc.SyncCoreConfig()
	return nil
}
