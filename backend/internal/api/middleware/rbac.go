package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE" {
			role, exists := c.Get("role")
			if !exists || role != "admin" {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden: Admin access required"})
				return
			}
		}
		c.Next()
	}
}
