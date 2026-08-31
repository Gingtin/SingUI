package cronjob

import (
	"log"
	"time"

	"github.com/singbox-ui/singbox-ui/internal/database"
	"github.com/singbox-ui/singbox-ui/internal/database/models"
	"github.com/singbox-ui/singbox-ui/internal/service"
)

type CronManager struct {
	stopChan chan struct{}
}

var GlobalCron *CronManager

func StartCronJobs() *CronManager {
	cm := &CronManager{
		stopChan: make(chan struct{}),
	}
	GlobalCron = cm

	go cm.runTrafficResetSchedule()
	go cm.runTelegramDailyReport()

	log.Println("[Cronjob] Background scheduler initialized.")
	return cm
}

func (cm *CronManager) Stop() {
	close(cm.stopChan)
}

// runTrafficResetSchedule checks daily whether a client's monthly reset day has arrived
func (cm *CronManager) runTrafficResetSchedule() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-cm.stopChan:
			return
		case <-ticker.C:
			currentDay := time.Now().Day()
			var clients []models.Client
			database.DB.Where("reset_day = ?", currentDay).Find(&clients)

			for _, c := range clients {
				_ = database.DB.Model(&c).Updates(map[string]interface{}{
					"up":   0,
					"down": 0,
				}).Error
				log.Printf("[Cronjob] Automatically reset traffic for client: %s (Reset Day: %d)\n", c.Email, c.ResetDay)
			}
		}
	}
}

// runTelegramDailyReport sends daily summary to Telegram if configured
func (cm *CronManager) runTelegramDailyReport() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-cm.stopChan:
			return
		case <-ticker.C:
			var botTokenSetting, chatIdSetting models.Setting
			database.DB.Where("key = ?", "tg_bot_token").First(&botTokenSetting)
			database.DB.Where("key = ?", "tg_chat_id").First(&chatIdSetting)

			if botTokenSetting.Value != "" && chatIdSetting.Value != "" {
				var clientCount int64
				database.DB.Model(&models.Client{}).Count(&clientCount)
				var inboundCount int64
				database.DB.Model(&models.Inbound{}).Count(&inboundCount)

				msg := service.FormatDailyReport(inboundCount, clientCount)
				_ = service.SendTelegramMessage(botTokenSetting.Value, chatIdSetting.Value, msg)
			}
		}
	}
}
