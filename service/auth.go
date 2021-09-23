package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/yakuter/ugin/model"
	"github.com/yakuter/ugin/pkg/config"
)

// CreateToken method
func CreateToken(email string) (*model.TokenDetailsDTO, error) {

	config := config.GetConfig()

	var err error
	td := &model.TokenDetailsDTO{}

	td.AtExpiresTime = time.Now().Add(time.Hour * time.Duration(config.Server.AccessTokenExpireDuration))
	td.RtExpiresTime = time.Now().Add(time.Hour * time.Duration(config.Server.RefreshTokenExpireDuration))

	//create access token
	atClaims := jwt.MapClaims{}
	atClaims["email"] = email
	atClaims["user_uuid"] = "user_uuid"
	atClaims["exp"] = time.Now().Add(time.Hour * 48).Unix()
	atClaims["uuid"] = ""
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	td.AccessToken, err = token.SignedString([]byte(config.Server.Secret))
	if err != nil {
		return nil, err
	}

	//create refresh  token
	rtClaims := jwt.MapClaims{}
	rtClaims["email"] = email
	rtClaims["user_uuid"] = "user_uuid"
	rtClaims["exp"] = time.Now().Add(time.Hour * 96).Unix() //refresh token expire time config read
	rtClaims["uuid"] = ""
	rtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)

	td.RefreshToken, err = rtoken.SignedString([]byte(config.Server.Secret))
	if err != nil {
		return nil, err
	}

	generatedPass, err := GenerateSecureKey(16)
	if err != nil {
		return nil, err
	}
	td.TransmissionKey = generatedPass

	return td, nil
}

// TokenValid method
func TokenValid(bearerToken string) (*jwt.Token, error) {
	token, err := verifyToken(bearerToken)
	if err != nil {
		if token != nil {
			return token, err
		}
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("Unauthorized")
	}
	return token, nil
}

//verifyToken verify token
func verifyToken(tokenString string) (*jwt.Token, error) {
	config := config.GetConfig()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Server.Secret), nil
	})
	if err != nil {
		return token, fmt.Errorf("Unauthorized")
	}
	return token, nil
}
