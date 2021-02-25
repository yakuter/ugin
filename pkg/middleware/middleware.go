package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/yakuter/ugin/service"
)

var (
	validToken   = "Token is valid"
	invalidToken = "Token is expired or not valid!"
	noToken      = "Token could not found! "
)

// CORS middleware
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// Authorize middleware
func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string
		bearerToken := c.GetHeader("Authorization")
		strArr := strings.Split(bearerToken, " ")
		if len(strArr) == 2 {
			tokenStr = strArr[1]
		}

		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, noToken)
			return
		}

		token, err := service.TokenValid(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, invalidToken)
			return
		}

		if err == nil && token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			fmt.Println(claims)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": invalidToken,
				"token":   nil,
			})
		}

	}
}
