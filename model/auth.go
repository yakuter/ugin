package model

import (
	"time"
)

// AuthLoginDTO model
type AuthLoginDTO struct {
	Email          string `validate:"required" json:"email"`
	MasterPassword string `validate:"required" json:"master_password"`
}

// AuthLoginResponse ...
type AuthLoginResponse struct {
	AccessToken     string `json:"access_token"`
	RefreshToken    string `json:"refresh_token"`
	TransmissionKey string `json:"transmission_key"`
}

// TokenDetailsDTO model
type TokenDetailsDTO struct {
	AccessToken     string `json:"access_token"`
	RefreshToken    string `json:"refresh_token"`
	AtExpiresTime   time.Time
	RtExpiresTime   time.Time
	TransmissionKey string `json:"transmission_key"`
}
