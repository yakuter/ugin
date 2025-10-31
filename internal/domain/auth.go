package domain

import "time"

// TokenDetails contains JWT token information
type TokenDetails struct {
	AccessToken     string    `json:"access_token"`
	RefreshToken    string    `json:"refresh_token"`
	TransmissionKey string    `json:"transmission_key"`
	ATExpiresAt     time.Time `json:"access_token_expires_at"`
	RTExpiresAt     time.Time `json:"refresh_token_expires_at"`
}

// Credentials represents user login credentials
type Credentials struct {
	Email          string `json:"email" binding:"required,email" example:"user@example.com"`
	MasterPassword string `json:"master_password" binding:"required,min=6" example:"password123"`
}

// TokenClaims represents the claims in a JWT token
type TokenClaims struct {
	Email    string `json:"email"`
	UserUUID string `json:"user_uuid"`
	UUID     string `json:"uuid"`
}

