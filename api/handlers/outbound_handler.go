package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/singbox-ui/singbox-ui/internal/database/models"
	"github.com/singbox-ui/singbox-ui/internal/service"
)

func ListOutbounds(c *gin.Context) {
	outbounds, err := service.OutboundSvc.ListOutbounds()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, outbounds)
}

func CreateOutbound(c *gin.Context) {
	var out models.Outbound
	if err := c.ShouldBindJSON(&out); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := service.OutboundSvc.CreateOutbound(&out); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func UpdateOutbound(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var out models.Outbound
	if err := c.ShouldBindJSON(&out); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out.ID = uint(id)
	if err := service.OutboundSvc.UpdateOutbound(&out); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func DeleteOutbound(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := service.OutboundSvc.DeleteOutbound(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Outbound deleted successfully"})
}
