package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/singbox-ui/singbox-ui/internal/database/models"
	"github.com/singbox-ui/singbox-ui/internal/service"
	"github.com/singbox-ui/singbox-ui/internal/util"
)

func ListInbounds(c *gin.Context) {
	inbounds, err := service.InboundSvc.ListInbounds()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, inbounds)
}

func GetInbound(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	inbound, err := service.InboundSvc.GetInbound(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inbound not found"})
		return
	}
	c.JSON(http.StatusOK, inbound)
}

func CreateInbound(c *gin.Context) {
	var in models.Inbound
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := service.InboundSvc.CreateInbound(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, in)
}

func UpdateInbound(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var in models.Inbound
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	in.ID = uint(id)
	if err := service.InboundSvc.UpdateInbound(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, in)
}

func DeleteInbound(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := service.InboundSvc.DeleteInbound(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Inbound deleted successfully"})
}

func AddClient(c *gin.Context) {
	idParam := c.Param("id")
	inboundID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid inbound ID"})
		return
	}

	var client models.Client
	if err := c.ShouldBindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := service.InboundSvc.AddClient(uint(inboundID), &client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, client)
}

func UpdateClient(c *gin.Context) {
	clientIdParam := c.Param("clientId")
	clientId, err := strconv.ParseUint(clientIdParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	var client models.Client
	if err := c.ShouldBindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client.ID = uint(clientId)
	if err := service.InboundSvc.UpdateClient(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, client)
}

func DeleteClient(c *gin.Context) {
	clientIdParam := c.Param("clientId")
	clientId, err := strconv.ParseUint(clientIdParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	if err := service.InboundSvc.DeleteClient(uint(clientId)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Client deleted successfully"})
}

func ResetClientTraffic(c *gin.Context) {
	clientIdParam := c.Param("clientId")
	clientId, err := strconv.ParseUint(clientIdParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	if err := service.InboundSvc.ResetClientTraffic(uint(clientId)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Traffic reset successfully"})
}

func ResetAllTraffic(c *gin.Context) {
	if err := service.InboundSvc.ResetAllTraffic(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "All traffic reset successfully"})
}

func GenerateRealityKeypair(c *gin.Context) {
	keypair, err := util.GenerateRealityKeypair()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"private_key": keypair.PrivateKey,
		"public_key":  keypair.PublicKey,
		"short_id":    util.GenerateShortID(),
	})
}

func GenerateRandomUUID(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"uuid": util.GenerateUUID(),
	})
}
