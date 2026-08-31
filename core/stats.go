package core

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/singbox-ui/singbox-ui/internal/database"
	"github.com/singbox-ui/singbox-ui/internal/database/models"
	"gorm.io/gorm"
)

type ClashConnection struct {
	ID          string    `json:"id"`
	Metadata    Metadata  `json:"metadata"`
	Upload      int64     `json:"upload"`
	Download    int64     `json:"download"`
	Start       time.Time `json:"start"`
	Rule        string    `json:"rule"`
	RulePayload string    `json:"rulePayload"`
}

type Metadata struct {
	Network         string `json:"network"`
	Type            string `json:"type"`
	SourceIP        string `json:"sourceIP"`
	DestinationIP   string `json:"destinationIP"`
	SourcePort      string `json:"sourcePort"`
	DestinationPort string `json:"destinationPort"`
	Host            string `json:"host"`
	InboundTag      string `json:"inboundTag"`
	InboundUser     string `json:"inboundUser"`
}

type ClashConnectionsResponse struct {
	DownloadTotal int64             `json:"downloadTotal"`
	UploadTotal   int64             `json:"uploadTotal"`
	Connections   []ClashConnection `json:"connections"`
}

type connSnapshot struct {
	Upload   int64
	Download int64
}

type StatsManager struct {
	clashPort     string
	clashSecret   string
	client        *http.Client
	mu            sync.Mutex
	lastConnStats map[string]connSnapshot
}

var StatsInstance *StatsManager

func InitStatsManager(port, secret string) *StatsManager {
	StatsInstance = &StatsManager{
		clashPort:     port,
		clashSecret:   secret,
		client:        &http.Client{Timeout: 3 * time.Second},
		lastConnStats: make(map[string]connSnapshot),
	}
	go StatsInstance.startPeriodicSync()
	return StatsInstance
}

func (sm *StatsManager) GetActiveConnections() (*ClashConnectionsResponse, error) {
	url := fmt.Sprintf("http://127.0.0.1:%s/connections", sm.clashPort)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if sm.clashSecret != "" {
		req.Header.Set("Authorization", "Bearer "+sm.clashSecret)
	}

	resp, err := sm.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data ClashConnectionsResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (sm *StatsManager) disconnect(connID string) {
	url := fmt.Sprintf("http://127.0.0.1:%s/connections/%s", sm.clashPort, connID)
	req, _ := http.NewRequest("DELETE", url, nil)
	if sm.clashSecret != "" {
		req.Header.Set("Authorization", "Bearer "+sm.clashSecret)
	}
	sm.client.Do(req)
}

func (sm *StatsManager) syncTrafficAndConnections() {
	active, err := sm.GetActiveConnections()
	if err != nil || active == nil {
		return
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	clientDeltas := make(map[string]struct {
		Up   int64
		Down int64
	})

	clientIPs := make(map[string]map[string][]string)
	currentConns := make(map[string]bool)

	for _, conn := range active.Connections {
		currentConns[conn.ID] = true
		key := conn.Metadata.InboundUser
		if key == "" {
			continue // No user identified
		}

		deltaUp := conn.Upload
		deltaDown := conn.Download
		if last, ok := sm.lastConnStats[conn.ID]; ok {
			deltaUp -= last.Upload
			deltaDown -= last.Download
			if deltaUp < 0 {
				deltaUp = 0
			}
			if deltaDown < 0 {
				deltaDown = 0
			}
		}

		sm.lastConnStats[conn.ID] = connSnapshot{
			Upload:   conn.Upload,
			Download: conn.Download,
		}

		if deltaUp > 0 || deltaDown > 0 {
			stat := clientDeltas[key]
			stat.Up += deltaUp
			stat.Down += deltaDown
			clientDeltas[key] = stat
		}

		if clientIPs[key] == nil {
			clientIPs[key] = make(map[string][]string)
		}
		clientIPs[key][conn.Metadata.SourceIP] = append(clientIPs[key][conn.Metadata.SourceIP], conn.ID)
	}

	for connID := range sm.lastConnStats {
		if !currentConns[connID] {
			delete(sm.lastConnStats, connID)
		}
	}

	if database.DB != nil {
		for key, deltas := range clientDeltas {
			database.DB.Model(&models.Client{}).Where("uuid = ? OR email = ? OR password = ?", key, key, key).
				Updates(map[string]interface{}{
					"up":   gorm.Expr("up + ?", deltas.Up),
					"down": gorm.Expr("down + ?", deltas.Down),
				})
		}

		var clients []models.Client
		if err := database.DB.Where("enable = ?", true).Find(&clients).Error; err == nil {
			for _, c := range clients {
				if c.LimitIP <= 0 {
					continue
				}

				var ipMap map[string][]string
				if m, ok := clientIPs[c.UUID]; ok {
					ipMap = m
				} else if m, ok := clientIPs[c.Email]; ok {
					ipMap = m
				} else if m, ok := clientIPs[c.Password]; ok {
					ipMap = m
				}

				if ipMap != nil && len(ipMap) > c.LimitIP {
					ipCount := 0
					for _, conns := range ipMap {
						ipCount++
						if ipCount > c.LimitIP {
							for _, connID := range conns {
								sm.disconnect(connID)
							}
						}
					}
				}
			}
		}
	}
}

func (sm *StatsManager) startPeriodicSync() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		sm.syncTrafficAndConnections()
		sm.checkLimitsAndQuotas()
	}
}

// checkLimitsAndQuotas verifies client quotas and expiration dates
func (sm *StatsManager) checkLimitsAndQuotas() {
	if database.DB == nil {
		return
	}

	var clients []models.Client
	if err := database.DB.Where("enable = ?", true).Find(&clients).Error; err != nil {
		return
	}

	nowMs := time.Now().UnixMilli()
	needReload := false

	for _, c := range clients {
		shouldDisable := false
		reason := ""

		// Check traffic quota
		if c.Total > 0 && (c.Up+c.Down) >= c.Total {
			shouldDisable = true
			reason = fmt.Sprintf("Traffic quota exceeded (%d/%d bytes)", c.Up+c.Down, c.Total)
		}

		// Check expiration time
		if c.ExpiryTime > 0 && nowMs >= c.ExpiryTime {
			shouldDisable = true
			reason = fmt.Sprintf("Account expired at %s", time.UnixMilli(c.ExpiryTime).Format("2006-01-02 15:04:05"))
		}

		if shouldDisable {
			log.Printf("[Stats] Disabling client %s (ID: %d): %s\n", c.Email, c.ID, reason)
			database.DB.Model(&c).Update("enable", false)
			needReload = true
		}
	}

	if needReload {
		var inbounds []models.Inbound
		database.DB.Preload("Clients").Find(&inbounds)
		var portSetting, secretSetting models.Setting
		database.DB.Where("key = ?", "clash_api_port").First(&portSetting)
		database.DB.Where("key = ?", "clash_api_secret").First(&secretSetting)

		var cfgSetting models.Setting
		database.DB.Where("key = ?", "singbox_config_path").First(&cfgSetting)
		configPath := cfgSetting.Value
		if configPath == "" {
			configPath = "config/singbox_config.json"
		}

		var rules []models.RoutingRule
		database.DB.Order("`order` asc, id asc").Find(&rules)
		var outbounds []models.Outbound
		database.DB.Find(&outbounds)
		var dns models.DNSSettings
		database.DB.First(&dns)

		cfg, err := GenerateConfig(inbounds, outbounds, rules, dns, portSetting.Value, secretSetting.Value)
		if err == nil {
			_ = WriteConfigToFile(cfg, configPath)
			if Instance != nil && Instance.GetStatus().IsRunning {
				_ = Instance.Restart()
			}
		}
	}
}
