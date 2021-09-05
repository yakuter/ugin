package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/didip/tollbooth"
	"github.com/gin-gonic/gin"
	"github.com/yakuter/ugin/pkg/config"
	"github.com/yakuter/ugin/service"
)

var (
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

// MyLimit middleware
func MyLimit() gin.HandlerFunc {
	config := config.GetConfig()
	lmt := tollbooth.NewLimiter(config.Server.LimitCountPerRequest, nil)
	return func(c *gin.Context) {
		httpError := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if httpError != nil {
			c.Data(httpError.StatusCode, lmt.GetMessageContentType(), []byte(httpError.Message))
			c.Abort()
		} else {
			c.Next()
		}
	}
}

// Security middleware
func Security() gin.HandlerFunc {
	return func(c *gin.Context) {
		// X-XSS-Protection
		c.Writer.Header().Add("X-XSS-Protection", "1; mode=block")

		// HTTP Strict Transport Security
		c.Writer.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		// X-Frame-Options
		c.Writer.Header().Add("X-Frame-Options", "SAMEORIGIN")

		// X-Content-Type-Options
		c.Writer.Header().Add("X-Content-Type-Options", "nosniff")

		// Content Security Policy
		c.Writer.Header().Add("Content-Security-Policy", "default-src 'self';")

		// X-Permitted-Cross-Domain-Policies
		c.Writer.Header().Add("X-Permitted-Cross-Domain-Policies", "none")

		// Referrer-Policy
		c.Writer.Header().Add("Referrer-Policy", "no-referrer")

		// Feature-Policy
		c.Writer.Header().Add("Feature-Policy", "microphone 'none'; camera 'none'")

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
