package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/singbox-ui/singbox-ui/internal/api"
	"github.com/singbox-ui/singbox-ui/internal/core"
	"github.com/singbox-ui/singbox-ui/internal/cronjob"
	"github.com/singbox-ui/singbox-ui/internal/database"
	"github.com/singbox-ui/singbox-ui/internal/database/models"
	"github.com/singbox-ui/singbox-ui/internal/service"
)

//go:embed dist/*
var webDistFS embed.FS

var (
	Version   = "1.0.0"
	BuildTime = "2026-08-31"
)

func main() {
	portFlag := flag.String("p", "", "Web panel listening port (default: from DB or 2096)")
	dbPathFlag := flag.String("d", "data/singbox_ui.db", "SQLite database file path")
	resetAdminFlag := flag.Bool("reset-admin", false, "Reset admin password to default 'admin'")
	versionFlag := flag.Bool("v", false, "Show version")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("SingUI version %s (Built at %s)\n", Version, BuildTime)
		return
	}

	log.Printf("[SingUI] Starting SingUI v%s...\n", Version)

	// 1. Initialize SQLite Database (with WAL mode enabled)
	db, err := database.InitDB(*dbPathFlag)
	if err != nil {
		log.Fatalf("[DB] Failed to initialize database: %v\n", err)
	}

	// Handle Admin Reset
	if *resetAdminFlag {
		var admin models.User
		if err := db.Where("username = ?", "admin").First(&admin).Error; err == nil {
			_ = admin.SetPassword("admin")
			db.Save(&admin)
			log.Println("[Auth] Admin password has been reset to: admin")
		} else {
			admin = models.User{Username: "admin", Role: "admin"}
			_ = admin.SetPassword("admin")
			db.Create(&admin)
			log.Println("[Auth] Created admin user with password: admin")
		}
		return
	}

	// 2. Initialize Sing-box Process Supervisor
	var binPathSetting, configPathSetting models.Setting
	db.Where("key = ?", "singbox_bin_path").First(&binPathSetting)
	db.Where("key = ?", "singbox_config_path").First(&configPathSetting)

	binPath := binPathSetting.Value
	if binPath == "" {
		binPath = "sing-box"
	}
	configPath := configPathSetting.Value
	if configPath == "" {
		configPath = "config/singbox_config.json"
	}

	supervisor := core.InitSupervisor(binPath, configPath)
	_ = service.InboundSvc.SyncCoreConfig()
	_ = supervisor.Start()

	// 3. Initialize Stats Engine & Clash API Polling
	var clashPortSetting, clashSecretSetting models.Setting
	db.Where("key = ?", "clash_api_port").First(&clashPortSetting)
	db.Where("key = ?", "clash_api_secret").First(&clashSecretSetting)

	core.InitStatsEngine(clashPortSetting.Value, clashSecretSetting.Value)

	// 4. Start Background Cronjobs
	cronManager := cronjob.StartCronJobs()
	defer cronManager.Stop()

	// 5. Setup Web Router & Embedded Assets
	var distSubFS fs.FS
	if sub, err := fs.Sub(webDistFS, "dist"); err == nil {
		distSubFS = sub
	}

	router := api.SetupRouter(distSubFS)

	// Determine Listening Port
	listenPort := "2096"
	if *portFlag != "" {
		listenPort = *portFlag
	} else {
		var portSetting models.Setting
		if err := db.Where("key = ?", "web_port").First(&portSetting).Error; err == nil && portSetting.Value != "" {
			listenPort = portSetting.Value
		}
	}

	srv := &http.Server{
		Addr:    ":" + listenPort,
		Handler: router,
	}

	// 6. Graceful Shutdown Listener
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("[SingUI] Shutting down SingUI gracefully...")
		if core.Instance != nil {
			_ = core.Instance.Stop()
		}
		_ = srv.Close()
		os.Exit(0)
	}()

	log.Printf("[SingUI] Web Panel listening on http://0.0.0.0:%s\n", listenPort)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("[SingUI] Server failed to start: %v\n", err)
	}
}
