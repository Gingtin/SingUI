package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/singbox-ui/singbox-ui/internal/service"
)

type VersionHandler struct{}

var VersionHdl = &VersionHandler{}

// CheckVersions returns version status of Panel, Sing-box core and GeoIP
func (h *VersionHandler) CheckVersions(c *gin.Context) {
	info, err := service.CheckAllVersions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}

// UpdateCore upgrades or downgrades the Sing-box core binary
func (h *VersionHandler) UpdateCore(c *gin.Context) {
	var req struct {
		Version string `json:"version" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid version"})
		return
	}

	if err := service.UpdateSingboxCore(req.Version); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sing-box core updated successfully to " + req.Version})
}

// UpdateGeo updates GeoIP / Geosite databases
func (h *VersionHandler) UpdateGeo(c *gin.Context) {
	if err := service.UpdateGeoDatabases(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "GeoIP/Geosite databases updated successfully"})
}
