package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/singbox-ui/singbox-ui/internal/core"
	"github.com/singbox-ui/singbox-ui/internal/util"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func GetSystemStatus(c *gin.Context) {
	status, err := util.GetSystemStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, status)
}

func GetCoreStatus(c *gin.Context) {
	if core.Instance == nil {
		c.JSON(http.StatusOK, gin.H{"is_running": false, "version": "Not initialized"})
		return
	}
	c.JSON(http.StatusOK, core.Instance.GetStatus())
}

func StartCore(c *gin.Context) {
	if core.Instance == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Supervisor not initialized"})
		return
	}
	if err := core.Instance.Start(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sing-box started successfully"})
}

func StopCore(c *gin.Context) {
	if core.Instance == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Supervisor not initialized"})
		return
	}
	if err := core.Instance.Stop(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sing-box stopped successfully"})
}

func RestartCore(c *gin.Context) {
	if core.Instance == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Supervisor not initialized"})
		return
	}
	if err := core.Instance.Restart(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sing-box restarted successfully"})
}

func GetLogs(c *gin.Context) {
	if core.Instance == nil {
		c.JSON(http.StatusOK, []string{})
		return
	}
	logs := core.Instance.GetRecentLogs()
	c.JSON(http.StatusOK, logs)
}

func StreamLogsWS(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	if core.Instance == nil {
		_ = ws.WriteMessage(websocket.TextMessage, []byte("[System] Core supervisor not initialized"))
		return
	}

	// Send initial buffer
	for _, line := range core.Instance.GetRecentLogs() {
		if err := ws.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
			return
		}
	}

	logChan := core.Instance.SubscribeLogs()
	defer core.Instance.UnsubscribeLogs(logChan)

	done := make(chan struct{})
	go func() {
		for {
			if _, _, err := ws.ReadMessage(); err != nil {
				close(done)
				return
			}
		}
	}()

	for {
		select {
		case line, ok := <-logChan:
			if !ok {
				return
			}
			if err := ws.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
				return
			}
		case <-done:
			return
		}
	}
}

func GetRawConfig(c *gin.Context) {
	var cfgSetting models.Setting
	database.DB.Where("key = ?", "singbox_config_path").First(&cfgSetting)
	cfgPath := cfgSetting.Value
	if cfgPath == "" {
		cfgPath = "config/singbox_config.json"
	}

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"config": "{}"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"config": string(data)})
}

func GetActiveConnections(c *gin.Context) {
	if core.StatsInstance == nil {
		c.JSON(http.StatusOK, gin.H{"connections": []interface{}{}, "uploadTotal": 0, "downloadTotal": 0})
		return
	}

	resp, err := core.StatsInstance.GetActiveConnections()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"connections": []interface{}{}, "uploadTotal": 0, "downloadTotal": 0, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
