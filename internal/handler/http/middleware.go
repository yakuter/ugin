package http

import (
	"net/http"
	"strings"

	"github.com/didip/tollbooth"
	"github.com/gin-gonic/gin"
	"github.com/yakuter/ugin/internal/service"
)

var (
	invalidToken = "Token is expired or not valid!"
	noToken      = "Token could not be found!"
)

// CORS middleware
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}

// RateLimit middleware
func RateLimit(limitPerRequest float64) gin.HandlerFunc {
	lmt := tollbooth.NewLimiter(limitPerRequest, nil)
	
	return func(c *gin.Context) {
		httpError := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if httpError != nil {
			c.Data(httpError.StatusCode, lmt.GetMessageContentType(), []byte(httpError.Message))
			c.Abort()
			return
		}
		c.Next()
	}
}

// Security middleware adds security headers
func Security() gin.HandlerFunc {
	return func(c *gin.Context) {
		// X-XSS-Protection
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")

		// HTTP Strict Transport Security
		c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		// X-Frame-Options
		c.Writer.Header().Set("X-Frame-Options", "SAMEORIGIN")

		// X-Content-Type-Options
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")

		// Content Security Policy
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self';")

		// X-Permitted-Cross-Domain-Policies
		c.Writer.Header().Set("X-Permitted-Cross-Domain-Policies", "none")

		// Referrer-Policy
		c.Writer.Header().Set("Referrer-Policy", "no-referrer")

		// Feature-Policy
		c.Writer.Header().Set("Feature-Policy", "microphone 'none'; camera 'none'")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// JWTAuth middleware validates JWT tokens
func JWTAuth(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		
		// Get token from Authorization header
		bearerToken := c.GetHeader("Authorization")
		if bearerToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": noToken})
			return
		}

		// Extract token
		parts := strings.Split(bearerToken, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		tokenStr := parts[1]

		// Validate token
		claims, err := authService.ValidateToken(ctx, tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": invalidToken})
			return
		}

		// Store claims in context for handlers to use
		c.Set("email", claims.Email)
		c.Set("user_uuid", claims.UserUUID)
		
		c.Next()
	}
}

