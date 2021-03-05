package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/yakuter/ugin/model"
	"github.com/yakuter/ugin/service"
)

var (
	userLoginErr   = "User email or master password is wrong."
	userVerifyErr  = "Please verify your email first."
	invalidUser    = "Invalid user"
	validToken     = "Token is valid"
	invalidToken   = "Token is expired or not valid!"
	noToken        = "Token could not found! "
	tokenCreateErr = "Token could not be created"
	signupSuccess  = "User created successfully"
	verifySuccess  = "Email verified successfully"
)

// Signup Controller
func (base *Controller) Signup(c *gin.Context) {

	// Run validator according to model validator tags
	// Check and verify the recaptcha response token. This is needed for web form security.
	// Check if user exist in database
	// Create new user
	// Send confirmation email to new user

}

// Signin godoc
// @Summary Signin
// @Description Signin Process
// @Tags auth
// @Accept  json
// @Produce  json
// @Param Signin body object true "Signin"
// @Success 200 "Success"
// @Router /auth/signin [post]
func (base *Controller) Signin(c *gin.Context) {
	var credential model.AuthLoginDTO
	err := c.ShouldBindJSON(&credential)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": noToken,
		})
		return
	}

	isUserAuthenticated := service.FindByCredentials(credential.Email, credential.MasterPassword)

	if isUserAuthenticated {

		token, err := service.CreateToken(credential.Email)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, tokenCreateErr)
			return
		}

		authLoginResponse := model.AuthLoginResponse{
			AccessToken:     token.AccessToken,
			RefreshToken:    token.RefreshToken,
			TransmissionKey: token.TransmissionKey,
		}

		c.JSON(http.StatusOK, authLoginResponse)
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": userLoginErr,
			"token":   nil,
		})
	}

}

// RefreshToken godoc
// @Summary Refresh Token
// @Description RefreshToken Process
// @Tags auth
// @Accept  json
// @Produce  json
// @Param RefreshToken body object true "RefreshToken"
// @Success 200 "Success"
// @Router /auth/refresh [post]
func (base *Controller) RefreshToken(c *gin.Context) {

	mapToken := map[string]string{}

	decoder := json.NewDecoder(c.Request.Body)
	if err := decoder.Decode(&mapToken); err != nil {
		errs := []string{"REFRESH_TOKEN_ERROR"}
		c.JSON(http.StatusUnprocessableEntity, errs)
		return
	}
	defer c.Request.Body.Close()

	token, err := service.TokenValid(mapToken["refresh_token"])

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, invalidToken)
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	//create token
	newtoken, err := service.CreateToken(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, tokenCreateErr)
		return
	}

	authLoginResponse := model.AuthLoginResponse{
		AccessToken:     newtoken.AccessToken,
		RefreshToken:    newtoken.RefreshToken,
		TransmissionKey: newtoken.TransmissionKey,
	}

	c.JSON(http.StatusOK, authLoginResponse)
}

// CheckToken godoc
// @Summary CheckToken
// @Description CheckToken header example
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Token header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 "success"
// @Router /auth/check [post]
func (base *Controller) CheckToken(c *gin.Context) {
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

	claims := token.Claims.(jwt.MapClaims)
	email := claims["email"]

	// Check if user exist in database and credentials are true

	c.JSON(http.StatusOK, email)
}
