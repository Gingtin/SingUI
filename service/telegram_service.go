package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SendTelegramMessage sends markdown notification to telegram bot
func SendTelegramMessage(token, chatID, text string) error {
	if token == "" || chatID == "" || text == "" {
		return nil
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	payload := map[string]string{
		"chat_id":    chatID,
		"text":       text,
		"parse_mode": "Markdown",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// FormatDailyReport generates a clean formatted daily report for Telegram
func FormatDailyReport(inboundCount, clientCount int64) string {
	nowStr := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("📊 *SingUI 每日运维概报*\n\n"+
		"⏰ *报告时间*: `%s`\n"+
		"🌐 *活跃入站节点*: `%d` 个\n"+
		"👥 *总客户端用户*: `%d` 人\n"+
		"🟢 *系统状态*: `Sing-box 运行正常`\n\n"+
		"_由 SingUI 自动推送_", nowStr, inboundCount, clientCount)
}

// FormatQuotaAlert generates quota warning alert
func FormatQuotaAlert(email string, usedGB, totalGB float64) string {
	return fmt.Sprintf("⚠️ *SingUI 流量配额耗尽预警*\n\n"+
		"👤 *用户*: `%s`\n"+
		"📈 *已用流量*: `%.2f GB` / `%.2f GB`\n"+
		"🚫 *状态*: `节点已自动熔断停用`\n", email, usedGB, totalGB)
}
