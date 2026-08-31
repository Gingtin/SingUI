package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/singbox-ui/singbox-ui/internal/service"
)

func GetSettings(c *gin.Context) {
	settings, err := service.SettingSvc.GetAllSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, settings)
}

func UpdateSettings(c *gin.Context) {
	var kvs map[string]string
	if err := c.ShouldBindJSON(&kvs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := service.SettingSvc.UpdateSettings(kvs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Settings updated successfully"})
}

type UpdatePasswordReq struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func UpdatePassword(c *gin.Context) {
	username, _ := c.Get("username")
	var req UpdatePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password params"})
		return
	}

	if err := service.SettingSvc.UpdateAdminPassword(username.(string), req.OldPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

func DownloadBackup(c *gin.Context) {
	dbPath := "data/singbox_ui.db"
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Database file not found"})
		return
	}

	filename := fmt.Sprintf("singbox_ui_backup_%s.db", time.Now().Format("20060102_150405"))
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.File(dbPath)
}

func RestoreBackup(c *gin.Context) {
	file, err := c.FormFile("backup_file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read uploaded backup file"})
		return
	}

	dbPath := "data/singbox_ui.db"
	tempPath := filepath.Join("data", "singbox_ui_restore.db")

	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Copy temp to active db
	src, err := os.Open(tempPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()

	dst, err := os.OpenFile(dbPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_ = os.Remove(tempPath)
	c.JSON(http.StatusOK, gin.H{"message": "Backup restored successfully. Please restart the panel."})
}
