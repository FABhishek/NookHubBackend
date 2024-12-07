package handlers

import (
	"Nookhub/models"
	"Nookhub/services"
	"net/http"

	"fmt"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	RegisterUser(c *gin.Context)
	LoginUser(c *gin.Context)
	IsEmailOrUsernameAvailable(c *gin.Context)
}

type userHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) UserHandler {
	return &userHandler{userService: userService}
}

// signup  related APIs
/*
API can return following http Responses
403: forbidded: User already registerd
200: success: successfully registered
*/
func (h *userHandler) RegisterUser(c *gin.Context) {
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

	userId, err := h.userService.RegisterUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	message := fmt.Sprintf("User registered successfully with ID %d", userId)

	c.JSON(http.StatusOK, gin.H{"message": message})
}

func (h *userHandler) IsEmailOrUsernameAvailable(c *gin.Context) {

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

	isValid, err := h.userService.IsAvailable(input, inputType)
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
func (h *userHandler) LoginUser(c *gin.Context) {
	var user models.LoginUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if isAnyEmpty(user.Email, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input, either of the required param is empty"})
	}

	username, userId, err := h.userService.LoginUser(user)

	if userId <= 0 && err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": fmt.Sprintf("User not registered, please register: %s", err)})
		return
	} else if userId > 0 && username == "existsButPWNotMatched" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": fmt.Sprintf("Wrong password entered %s", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"username": username, "userId": userId})
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
