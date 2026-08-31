package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/singbox-ui/singbox-ui/internal/database"
	"github.com/singbox-ui/singbox-ui/internal/database/models"
	"github.com/singbox-ui/singbox-ui/internal/util"
)

type TelegramBot struct {
	token  string
	chatID string
	client *http.Client
}

var Bot *TelegramBot

func InitTelegramBot(token, chatID string) *TelegramBot {
	if token == "" || chatID == "" {
		return nil
	}
	Bot = &TelegramBot{
		token:  token,
		chatID: chatID,
		client: &http.Client{Timeout: 10 * time.Second},
	}
	return Bot
}

func (b *TelegramBot) SendMessage(text string) error {
	if b == nil || b.token == "" || b.chatID == "" {
		return nil
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", b.token)
	payload := map[string]string{
		"chat_id":    b.chatID,
		"text":       text,
		"parse_mode": "HTML",
	}

	body, _ := json.Marshal(payload)
	resp, err := b.client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("[Telegram] Send message error: %v\n", err)
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (b *TelegramBot) SendStatusReport() {
	status, err := util.GetSystemStatus()
	if err != nil {
		return
	}

	var inboundsCount, clientsCount int64
	if database.DB != nil {
		database.DB.Model(&models.Inbound{}).Count(&inboundsCount)
		database.DB.Model(&models.Client{}).Count(&clientsCount)
	}

	msg := fmt.Sprintf(
		"🚀 <b>Sing-box UI 运行报告</b>\n\n"+
			"🖥 <b>系统状态</b>: %s (%s)\n"+
			"⏱ <b>运行时间</b>: %d 小时\n"+
			"📊 <b>CPU 使用率</b>: %.1f%%\n"+
			"💾 <b>内存使用率</b>: %.1f%% (%d MB / %d MB)\n"+
			"📁 <b>磁盘使用率</b>: %.1f%%\n"+
			"🌐 <b>网络速度</b>: ⬇️ %.2f MB/s | ⬆️ %.2f MB/s\n\n"+
			"👥 <b>节点数</b>: %d | <b>用户数</b>: %d\n",
		status.OS, status.Platform,
		status.Uptime/3600,
		status.CPUPercent,
		status.MemPercent, status.MemUsed/(1024*1024), status.MemTotal/(1024*1024),
		status.DiskPercent,
		float64(status.NetDownloadRate)/(1024*1024), float64(status.NetUploadRate)/(1024*1024),
		inboundsCount, clientsCount,
	)

	_ = b.SendMessage(msg)
}
