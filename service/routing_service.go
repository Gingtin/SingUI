package service

import (
	"github.com/singbox-ui/singbox-ui/internal/database"
	"github.com/singbox-ui/singbox-ui/internal/database/models"
)

type RoutingService struct{}

var RoutingSvc = &RoutingService{}

func (s *RoutingService) ListRules() ([]models.RoutingRule, error) {
	var rules []models.RoutingRule
	err := database.DB.Order("`order` asc, id asc").Find(&rules).Error
	return rules, err
}

func (s *RoutingService) CreateRule(rule *models.RoutingRule) error {
	if err := database.DB.Create(rule).Error; err != nil {
		return err
	}
	_ = InboundSvc.SyncCoreConfig()
	return nil
}

func (s *RoutingService) UpdateRule(rule *models.RoutingRule) error {
	if err := database.DB.Save(rule).Error; err != nil {
		return err
	}
	_ = InboundSvc.SyncCoreConfig()
	return nil
}

func (s *RoutingService) DeleteRule(id uint) error {
	if err := database.DB.Delete(&models.RoutingRule{}, id).Error; err != nil {
		return err
	}
	_ = InboundSvc.SyncCoreConfig()
	return nil
}

func (s *RoutingService) GetDNSSettings() (*models.DNSSettings, error) {
	var dns models.DNSSettings
	err := database.DB.First(&dns).Error
	return &dns, err
}

func (s *RoutingService) UpdateDNSSettings(dns *models.DNSSettings) error {
	var existing models.DNSSettings
	if err := database.DB.First(&existing).Error; err == nil {
		dns.ID = existing.ID
		if err := database.DB.Save(dns).Error; err != nil {
			return err
		}
	} else {
		if err := database.DB.Create(dns).Error; err != nil {
			return err
		}
	}
	_ = InboundSvc.SyncCoreConfig()
	return nil
}
