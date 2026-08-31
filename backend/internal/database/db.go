package database

import (
	"log"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/singbox-ui/singbox-ui/internal/database/models"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(dbPath string) (*gorm.DB, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Inbound{},
		&models.Client{},
		&models.Setting{},
		&models.RoutingRule{},
		&models.DNSSettings{},
	); err != nil {
		return nil, err
	}

	DB = db
	seedDefaultData(db)
	return db, nil
}

func seedDefaultData(db *gorm.DB) {
	// Seed Admin User
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	if userCount == 0 {
		admin := models.User{
			Username: "admin",
			Role:     "admin",
		}
		_ = admin.SetPassword("admin")
		db.Create(&admin)
		log.Println("[DB] Initialized default admin user (admin/admin). Please change password in settings!")
	}

	// Seed Settings
	defaultSettings := map[string]string{
		"web_port":            "2096",
		"web_base_path":       "/",
		"jwt_secret":          uuid.New().String(),
		"sub_domain":          "",
		"sub_path":            "/sub",
		"sub_title":           "SingUI Nodes",
		"singbox_bin_path":    "sing-box",
		"singbox_config_path": "config/singbox_config.json",
		"clash_api_port":      "9090",
		"clash_api_secret":    uuid.New().String()[:16],
		"tg_bot_token":        "",
		"tg_chat_id":          "",
		"tg_notify_on_expire": "true",
		"tg_notify_on_quota":  "true",
		"traffic_log_days":    "30",
		"auto_backup_enabled": "false",
		"acme_email":          "",
	}

	for k, v := range defaultSettings {
		var setting models.Setting
		if err := db.Where("key = ?", k).First(&setting).Error; err != nil {
			db.Create(&models.Setting{
				Key:   k,
				Value: v,
			})
		}
	}

	// Seed Default DNS Settings
	var dnsCount int64
	db.Model(&models.DNSSettings{}).Count(&dnsCount)
	if dnsCount == 0 {
		db.Create(&models.DNSSettings{
			LocalDNS:     "local",
			RemoteDNS:    "https://1.1.1.1/dns-query",
			ChinaDNS:     "https://223.5.5.5/dns-query",
			EnableFakeIP: false,
			Strategy:     "prefer_ipv4",
		})
	}

	// Seed Default Intelligent Routing Rules
	var ruleCount int64
	db.Model(&models.RoutingRule{}).Count(&ruleCount)
	if ruleCount == 0 {
		defaultRules := []models.RoutingRule{
			{Tag: "DNS Out", Protocol: "dns", Outbound: "dns-out", Enable: true, Order: 1, Remark: "拦截并分流 DNS 查询"},
			{Tag: "Block Ads", Domain: `["geosite:category-ads-all"]`, Outbound: "block", Enable: true, Order: 2, Remark: "拦截常见广告与跟踪域名"},
			{Tag: "Direct Private IP", IP: `["geoip:private"]`, Outbound: "direct", Enable: true, Order: 3, Remark: "局域网与私有 IP 直连"},
			{Tag: "Direct China Domains", Domain: `["geosite:cn"]`, Outbound: "direct", Enable: true, Order: 4, Remark: "国内主流域名直连"},
			{Tag: "Direct China IP", IP: `["geoip:cn"]`, Outbound: "direct", Enable: true, Order: 5, Remark: "中国大陆 IP 地址直连"},
		}
		for _, r := range defaultRules {
			db.Create(&r)
		}
	}
}
