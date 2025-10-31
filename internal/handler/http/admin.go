package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Dashboard handles GET /admin/dashboard
func Dashboard(c *gin.Context) {
	user := c.MustGet(gin.AuthUserKey).(string)
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to admin dashboard",
		"user":    user,
	})
}

