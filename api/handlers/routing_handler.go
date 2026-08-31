package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/singbox-ui/singbox-ui/internal/database/models"
	"github.com/singbox-ui/singbox-ui/internal/service"
)

func ListRoutingRules(c *gin.Context) {
	rules, err := service.RoutingSvc.ListRules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rules)
}

func CreateRoutingRule(c *gin.Context) {
	var rule models.RoutingRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := service.RoutingSvc.CreateRule(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rule)
}

func UpdateRoutingRule(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var rule models.RoutingRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rule.ID = uint(id)
	if err := service.RoutingSvc.UpdateRule(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rule)
}

func DeleteRoutingRule(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := service.RoutingSvc.DeleteRule(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Rule deleted successfully"})
}

func GetDNSSettings(c *gin.Context) {
	dns, err := service.RoutingSvc.GetDNSSettings()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"local_dns":     "local",
			"remote_dns":    "https://1.1.1.1/dns-query",
			"china_dns":     "https://223.5.5.5/dns-query",
			"strategy":      "prefer_ipv4",
			"enable_fakeip": false,
		})
		return
	}
	c.JSON(http.StatusOK, dns)
}

func UpdateDNSSettings(c *gin.Context) {
	var dns models.DNSSettings
	if err := c.ShouldBindJSON(&dns); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := service.RoutingSvc.UpdateDNSSettings(&dns); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "DNS settings updated successfully"})
}
