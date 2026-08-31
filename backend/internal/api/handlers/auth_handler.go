package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"github.com/singbox-ui/singbox-ui/internal/api/middleware"
	"github.com/singbox-ui/singbox-ui/internal/database"
	"github.com/singbox-ui/singbox-ui/internal/database/models"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	TwoFACode string `json:"two_fa_code"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Check 2FA if enabled
	if user.TwoFAEnabled && user.TwoFASecret != "" {
		if req.TwoFACode == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Two-factor authentication code required", "need_2fa": true})
			return
		}
		if !totp.Validate(req.TwoFACode, user.TwoFASecret) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid two-factor authentication code"})
			return
		}
	}

	token, err := middleware.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate session token"})
		return
	}

	c.SetCookie("singbox_ui_token", token, 7*24*3600, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":             user.ID,
			"username":       user.Username,
			"role":           user.Role,
			"two_fa_enabled": user.TwoFAEnabled,
		},
	})
}

func GetUserInfo(c *gin.Context) {
	userID, _ := c.Get("userID")
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":             user.ID,
		"username":       user.Username,
		"role":           user.Role,
		"two_fa_enabled": user.TwoFAEnabled,
	})
}

func Logout(c *gin.Context) {
	c.SetCookie("singbox_ui_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
