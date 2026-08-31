package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/singbox-ui/singbox-ui/internal/database"
	"github.com/singbox-ui/singbox-ui/internal/database/models"
	"github.com/singbox-ui/singbox-ui/internal/subscription"
)

func HandleSubscription(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.String(http.StatusBadRequest, "Invalid subscription token")
		return
	}

	var client models.Client
	if err := database.DB.Where("sub_token = ?", token).First(&client).Error; err != nil {
		c.String(http.StatusNotFound, "Subscription not found or expired")
		return
	}

	if !client.Enable {
		c.String(http.StatusForbidden, "Subscription is disabled")
		return
	}

	// Set standard Subscription-Userinfo header
	expireSec := client.ExpiryTime / 1000
	userinfo := fmt.Sprintf("upload=%d; download=%d; total=%d; expire=%d", client.Up, client.Down, client.Total, expireSec)
	c.Header("Subscription-Userinfo", userinfo)
	c.Header("Profile-Update-Interval", "6")

	// Get server host / domain from setting or request host
	var subDomainSetting models.Setting
	database.DB.Where("key = ?", "sub_domain").First(&subDomainSetting)
	serverHost := subDomainSetting.Value
	if serverHost == "" {
		serverHost = c.Request.Host
		if strings.Contains(serverHost, ":") {
			serverHost = strings.Split(serverHost, ":")[0]
		}
	}

	var inbounds []models.Inbound
	database.DB.Preload("Clients").Where("enable = ?", true).Find(&inbounds)

	builder := subscription.NewSubBuilder(serverHost, &client, inbounds)

	// Determine client format from query param or User-Agent
	flag := strings.ToLower(c.Query("flag"))
	ua := strings.ToLower(c.GetHeader("User-Agent"))

	if flag == "sing-box" || flag == "sbox" || strings.Contains(ua, "sing-box") || strings.Contains(ua, "s-box") {
		jsonStr, err := builder.BuildSingboxJSON()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.String(http.StatusOK, jsonStr)
		return
	}

	if flag == "clash" || flag == "meta" || flag == "mihomo" || strings.Contains(ua, "clash") || strings.Contains(ua, "meta") || strings.Contains(ua, "stash") || strings.Contains(ua, "verge") {
		yamlStr, err := builder.BuildClashMeta()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.Header("Content-Type", "text/yaml; charset=utf-8")
		c.String(http.StatusOK, yamlStr)
		return
	}

	// Default: Standard Base64 URI list
	base64Links := builder.BuildBase64Links()
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusOK, base64Links)
}

func HandleSubView(c *gin.Context) {
	token := c.Param("token")
	var client models.Client
	if err := database.DB.Where("sub_token = ?", token).First(&client).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	var subDomainSetting models.Setting
	database.DB.Where("key = ?", "sub_domain").First(&subDomainSetting)
	serverHost := subDomainSetting.Value
	if serverHost == "" {
		serverHost = c.Request.Host
		if strings.Contains(serverHost, ":") {
			serverHost = strings.Split(serverHost, ":")[0]
		}
	}

	var inbounds []models.Inbound
	database.DB.Preload("Clients").Where("enable = ?", true).Find(&inbounds)

	now := time.Now().UnixMilli()
	isExpired := client.ExpiryTime > 0 && now >= client.ExpiryTime
	isOverQuota := client.Total > 0 && (client.Up+client.Down) >= client.Total

	c.JSON(http.StatusOK, gin.H{
		"email":         client.Email,
		"up":            client.Up,
		"down":          client.Down,
		"total":         client.Total,
		"expiry_time":   client.ExpiryTime,
		"is_expired":    isExpired,
		"is_over_quota": isOverQuota,
		"sub_token":     client.SubToken,
		"server_host":   serverHost,
	})
}
