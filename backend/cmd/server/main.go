package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron/v3"
	"github.com/singbox-ui/singbox-ui/internal/api"
	"github.com/singbox-ui/singbox-ui/internal/core"
	"github.com/singbox-ui/singbox-ui/internal/database"
	"github.com/singbox-ui/singbox-ui/internal/database/models"
	"github.com/singbox-ui/singbox-ui/internal/service"
	"github.com/singbox-ui/singbox-ui/internal/telegram"
)

//go:embed dist/*
var embedWebFS embed.FS

const Version = "v1.0.0"

func main() {
	portFlag := flag.String("p", "", "Web panel listen port (default from DB or 2096)")
	dbFlag := flag.String("d", "data/singbox_ui.db", "SQLite database file path")
	versionFlag := flag.Bool("v", false, "Show version")
	resetAdminFlag := flag.Bool("reset-admin", false, "Reset admin password to 'admin'")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("Singbox-UI Panel %s\n", Version)
		return
	}

	log.Printf("[Main] Starting Singbox-UI Panel %s...\n", Version)

	// 1. Initialize SQLite Database
	db, err := database.InitDB(*dbFlag)
	if err != nil {
		log.Fatalf("[DB] Failed to initialize database: %v\n", err)
	}

	if *resetAdminFlag {
		var admin models.User
		if err := db.Where("username = ?", "admin").First(&admin).Error; err == nil {
			_ = admin.SetPassword("admin")
			db.Save(&admin)
			log.Println("[Main] Admin password has been reset to: admin")
			return
		}
	}

	// 2. Fetch Settings
	settings, _ := service.SettingSvc.GetAllSettings()
	port := settings["web_port"]
	if *portFlag != "" {
		port = *portFlag
	}
	if port == "" {
		port = "2096"
	}

	binPath := settings["singbox_bin_path"]
	if binPath == "" {
		binPath = "sing-box"
	}

	configPath := settings["singbox_config_path"]
	if configPath == "" {
		configPath = "config/singbox_config.json"
	}

	clashPort := settings["clash_api_port"]
	if clashPort == "" {
		clashPort = "9090"
	}
	clashSecret := settings["clash_api_secret"]

	// 3. Initialize Sing-box Core Supervisor & Config
	_ = service.InboundSvc.SyncCoreConfig()
	supervisor := core.InitSupervisor(binPath, configPath)
	if err := supervisor.Start(); err != nil {
		log.Printf("[Supervisor] Note: Failed to start sing-box binary (%s): %v. You can set correct binary path in Settings.\n", binPath, err)
	}

	// 4. Initialize Clash API Stats Manager
	core.InitStatsManager(clashPort, clashSecret)

	// 5. Initialize Telegram Bot
	if tgToken := settings["tg_bot_token"]; tgToken != "" {
		tgChatID := settings["tg_chat_id"]
		bot := telegram.InitTelegramBot(tgToken, tgChatID)
		if bot != nil {
			go bot.SendMessage(fmt.Sprintf("🚀 Sing-box UI %s has started on port %s", Version, port))
		}
	}

	// 6. Setup Cron Tasks (Daily traffic reset, periodic report)
	c := cron.New()
	_, _ = c.AddFunc("0 0 1 * *", func() {
		log.Println("[Cron] Running monthly traffic reset check...")
	})
	c.Start()
	defer c.Stop()

	// 7. Setup Static Web Assets
	var webFS fs.FS
	subFS, err := fs.Sub(embedWebFS, "dist")
	if err == nil {
		webFS = subFS
	}

	// 8. Start HTTP API & Web Server
	router := api.SetupRouter(webFS)

	// Graceful shutdown handling
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("[Main] Shutting down Singbox-UI...")
		if supervisor != nil {
			_ = supervisor.Stop()
		}
		os.Exit(0)
	}()

	addr := fmt.Sprintf("0.0.0.0:%s", port)
	log.Printf("[Main] Singbox-UI listening on http://%s\n", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("[Main] Server listen error: %v\n", err)
	}
}
