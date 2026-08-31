package service

import (
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/singbox-ui/singbox-ui/internal/core"
	"github.com/singbox-ui/singbox-ui/internal/database"
	"github.com/singbox-ui/singbox-ui/internal/database/models"
	"github.com/singbox-ui/singbox-ui/internal/util"
)

type InboundService struct{}

var InboundSvc = &InboundService{}

func (s *InboundService) ListInbounds() ([]models.Inbound, error) {
	var inbounds []models.Inbound
	err := database.DB.Preload("Clients").Order("id asc").Find(&inbounds).Error
	return inbounds, err
}

func (s *InboundService) GetInbound(id uint) (*models.Inbound, error) {
	var in models.Inbound
	err := database.DB.Preload("Clients").First(&in, id).Error
	if err != nil {
		return nil, err
	}
	return &in, nil
}

func (s *InboundService) CreateInbound(in *models.Inbound) error {
	var count int64
	database.DB.Model(&models.Inbound{}).Where("port = ?", in.Port).Count(&count)
	if count > 0 {
		return errors.New("port already in use")
	}

	if in.Tag == "" {
		in.Tag = fmt.Sprintf("inbound-%d", in.Port)
	}

	if len(in.Clients) == 0 {
		client := models.Client{
			Email:    "default-user",
			UUID:     util.GenerateUUID(),
			Password: util.GenerateRandomPassword(16),
			SubToken: uuid.New().String(),
			Enable:   true,
		}
		if in.Protocol == "vless" && in.Security == "reality" {
			client.Flow = "xtls-rprx-vision"
		}
		in.Clients = append(in.Clients, client)
	} else {
		for i := range in.Clients {
			if in.Clients[i].UUID == "" {
				in.Clients[i].UUID = util.GenerateUUID()
			}
			if in.Clients[i].Password == "" {
				in.Clients[i].Password = util.GenerateRandomPassword(16)
			}
			if in.Clients[i].SubToken == "" {
				in.Clients[i].SubToken = uuid.New().String()
			}
		}
	}

	if err := database.DB.Create(in).Error; err != nil {
		return err
	}

	_ = s.SyncCoreConfig()
	return nil
}

func (s *InboundService) UpdateInbound(in *models.Inbound) error {
	var existing models.Inbound
	if err := database.DB.First(&existing, in.ID).Error; err != nil {
		return err
	}

	existing.Tag = in.Tag
	existing.Port = in.Port
	existing.Listen = in.Listen
	existing.Network = in.Network
	existing.Security = in.Security
	existing.Settings = in.Settings
	existing.StreamSettings = in.StreamSettings
	existing.Sniffing = in.Sniffing
	existing.Enable = in.Enable
	existing.Remark = in.Remark

	if err := database.DB.Save(&existing).Error; err != nil {
		return err
	}

	_ = s.SyncCoreConfig()
	return nil
}

func (s *InboundService) DeleteInbound(id uint) error {
	if err := database.DB.Delete(&models.Inbound{}, id).Error; err != nil {
		return err
	}
	_ = s.SyncCoreConfig()
	return nil
}

func (s *InboundService) BatchDeleteInbounds(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	if err := database.DB.Where("id IN ?", ids).Delete(&models.Inbound{}).Error; err != nil {
		return err
	}
	_ = s.SyncCoreConfig()
	return nil
}

func (s *InboundService) BatchToggleInbounds(ids []uint, enable bool) error {
	if len(ids) == 0 {
		return nil
	}
	if err := database.DB.Model(&models.Inbound{}).Where("id IN ?", ids).Update("enable", enable).Error; err != nil {
		return err
	}
	_ = s.SyncCoreConfig()
	return nil
}

func (s *InboundService) AddClient(inboundID uint, client *models.Client) error {
	var in models.Inbound
	if err := database.DB.First(&in, inboundID).Error; err != nil {
		return err
	}

	client.InboundID = inboundID
	if client.UUID == "" {
		client.UUID = util.GenerateUUID()
	}
	if client.Password == "" {
		client.Password = util.GenerateRandomPassword(16)
	}
	if client.SubToken == "" {
		client.SubToken = uuid.New().String()
	}
	if in.Protocol == "vless" && in.Security == "reality" && client.Flow == "" {
		client.Flow = "xtls-rprx-vision"
	}

	if err := database.DB.Create(client).Error; err != nil {
		return err
	}

	_ = s.SyncCoreConfig()
	return nil
}

func (s *InboundService) UpdateClient(client *models.Client) error {
	var existing models.Client
	if err := database.DB.First(&existing, client.ID).Error; err != nil {
		return err
	}

	existing.Email = client.Email
	existing.Flow = client.Flow
	existing.Total = client.Total
	existing.ExpiryTime = client.ExpiryTime
	existing.Enable = client.Enable
	existing.LimitIP = client.LimitIP
	existing.ResetDay = client.ResetDay

	if err := database.DB.Save(&existing).Error; err != nil {
		return err
	}

	_ = s.SyncCoreConfig()
	return nil
}

func (s *InboundService) DeleteClient(clientID uint) error {
	if err := database.DB.Delete(&models.Client{}, clientID).Error; err != nil {
		return err
	}
	_ = s.SyncCoreConfig()
	return nil
}

func (s *InboundService) BatchDeleteClients(clientIDs []uint) error {
	if len(clientIDs) == 0 {
		return nil
	}
	if err := database.DB.Where("id IN ?", clientIDs).Delete(&models.Client{}).Error; err != nil {
		return err
	}
	_ = s.SyncCoreConfig()
	return nil
}

func (s *InboundService) ResetClientTraffic(clientID uint) error {
	return database.DB.Model(&models.Client{}).Where("id = ?", clientID).Updates(map[string]interface{}{
		"up":   0,
		"down": 0,
	}).Error
}

func (s *InboundService) ResetAllTraffic() error {
	return database.DB.Model(&models.Client{}).Where("1=1").Updates(map[string]interface{}{
		"up":   0,
		"down": 0,
	}).Error
}

func (s *InboundService) SyncCoreConfig() error {
	inbounds, err := s.ListInbounds()
	if err != nil {
		return err
	}

	var portSetting, secretSetting, cfgSetting, binSetting models.Setting
	database.DB.Where("key = ?", "clash_api_port").First(&portSetting)
	database.DB.Where("key = ?", "clash_api_secret").First(&secretSetting)
	database.DB.Where("key = ?", "singbox_config_path").First(&cfgSetting)
	database.DB.Where("key = ?", "singbox_bin_path").First(&binSetting)

	cfgPath := cfgSetting.Value
	if cfgPath == "" {
		cfgPath = "config/singbox_config.json"
	}

	var rules []models.RoutingRule
	database.DB.Order("`order` asc, id asc").Find(&rules)

	var outbounds []models.Outbound
	database.DB.Find(&outbounds)

	var dns models.DNSSettings
	database.DB.First(&dns)

	cfg, err := core.GenerateConfig(inbounds, outbounds, rules, dns, portSetting.Value, secretSetting.Value)
	if err != nil {
		return err
	}

	// Validate config syntax with sing-box binary if available
	binPath := binSetting.Value
	if binPath != "" {
		if valErr := core.ValidateConfig(binPath, cfg); valErr != nil {
			log.Printf("[ConfigSync] Warning: Sing-box validation returned: %v\n", valErr)
		}
	}

	if err := core.WriteConfigToFile(cfg, cfgPath); err != nil {
		return err
	}

	if core.Instance != nil && core.Instance.GetStatus().IsRunning {
		return core.Instance.Restart()
	}
	return nil
}
