package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yakuter/ugin/internal/domain"
	"github.com/yakuter/ugin/internal/repository"
	"github.com/yakuter/ugin/internal/service"
)

type AuthHandler struct {
	service service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// SignIn handles POST /auth/signin
// @Summary Sign in
// @Description User sign in with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body domain.Credentials true "User credentials"
// @Success 200 {object} domain.TokenDetails
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/auth/signin [post]
func (h *AuthHandler) SignIn(c *gin.Context) {
	ctx := c.Request.Context()

	var creds domain.Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	tokenDetails, err := h.service.SignIn(ctx, &creds)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, tokenDetails)
}

// SignUp handles POST /auth/signup
// @Summary Sign up
// @Description Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body domain.Credentials true "User credentials"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/auth/signup [post]
func (h *AuthHandler) SignUp(c *gin.Context) {
	ctx := c.Request.Context()

	var creds domain.Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	if err := h.service.SignUp(ctx, &creds); err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			return
		}
		if errors.Is(err, repository.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
}

// RefreshToken handles POST /auth/refresh
// @Summary Refresh token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh body map[string]string true "Refresh token" SchemaExample({"refresh_token": "your_refresh_token"})
// @Success 200 {object} domain.TokenDetails
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()

	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token is required"})
		return
	}

	tokenDetails, err := h.service.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		if errors.Is(err, service.ErrInvalidToken) || errors.Is(err, service.ErrExpiredToken) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, tokenDetails)
}

// CheckToken handles POST /auth/check
// @Summary Check token
// @Description Validate JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/auth/check [post]
func (h *AuthHandler) CheckToken(c *gin.Context) {
	ctx := c.Request.Context()

	// Get token from Authorization header
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
		return
	}

	// Remove "Bearer " prefix if present
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	claims, err := h.service.ValidateToken(ctx, token)
	if err != nil {
		if errors.Is(err, service.ErrInvalidToken) || errors.Is(err, service.ErrExpiredToken) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": true, "email": claims.Email})
}

