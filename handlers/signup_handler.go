package handlers

import (
	"Nookhub/models"
	"Nookhub/services"
	jwtutil "Nookhub/utilities"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const passwordNotMatched = "existsButPWNotMatched"

type SignupHandler interface {
	RegisterUser(c *gin.Context)
	LoginUser(c *gin.Context)
	IsEmailOrUsernameAvailable(c *gin.Context)
}

type signupHandler struct {
	signupService services.SignupService
}

func NewSignupHandler(signupService services.SignupService) SignupHandler {
	return &signupHandler{signupService: signupService}
}

// signup  related APIs
/*
API can return following http Responses
403: forbidded: User already registerd
200: success: successfully registered
*/
func (h *signupHandler) RegisterUser(c *gin.Context) {
	var user models.RegisterUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	fmt.Println(user)
	if isAnyEmpty(user.Username, user.Username, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input, either of the required param is empty"})
		return
	}

	userId, err := h.signupService.RegisterUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	message := fmt.Sprintf("User registered successfully with ID %d", userId)

	c.JSON(http.StatusOK, gin.H{"message": message})

	tokenString, err := jwtutil.CreateToken(user.Username, userId)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating the authentication token")
	}
	c.SetCookie(
		"token",
		tokenString,
		3600,
		"/",
		"localhost",
		false, //make sure to make it true later in https
		true)
}

func (h *signupHandler) IsEmailOrUsernameAvailable(c *gin.Context) {

	email := c.DefaultQuery("email", "")
	username := c.DefaultQuery("username", "")

	if isAnyEmpty(email) && isAnyEmpty(username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": `Invalid, input cannot be empty`})
	}

	var inputType string
	var input string
	if isAnyEmpty(email) {
		inputType = "username"
		input = username
	} else {
		inputType = "email"
		input = email
	}

	isValid, err := h.signupService.IsAvailable(input, inputType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if isValid {
		c.JSON(http.StatusOK, gin.H{"Valid": true})
		return
	} else {
		c.JSON(http.StatusConflict, gin.H{"Valid": false})
		return
	}
}

// Login related APIs
/*
API can return following http Responses
403: forbidded: wrong email/password provied
200: success: successfully logged in
*/
func (h *signupHandler) LoginUser(c *gin.Context) {
	var user models.LoginUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if isAnyEmpty(user.Email, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input, either of the required param is empty"})
	}

	username, userId, err := h.signupService.LoginUser(user)

	if userId <= 0 && err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "User not registered, please register"})
		return
	} else if userId > 0 && username == passwordNotMatched {
		c.JSON(http.StatusUnauthorized, gin.H{"message": fmt.Sprintf("Wrong password entered %s", err)})
		return
	}
	tokenString, err := jwtutil.CreateToken(username, userId)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating the authentication token")
	}
	c.SetCookie(
		"token",
		tokenString,
		3600,
		"/",
		"localhost",
		false, //make sure to make it true later in https
		true)
	c.JSON(http.StatusOK, gin.H{"username": username, "userId": userId, "message": "User successfully logged in and cookies have been set 😋"})
}

// private functions
func isAnyEmpty(strings ...string) bool {
	for _, str := range strings {
		if str == "" {
			return true
		}
	}
	return false
}
