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
		"sub_title":           "SingBox UI Nodes",
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
}
