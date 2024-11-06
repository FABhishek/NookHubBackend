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
}

type userHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) UserHandler {
	return &userHandler{userService: userService}
}

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

	userId, err := h.userService.RegisterUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	message := fmt.Sprintf("User registered successfully with ID %d", userId)

	c.JSON(http.StatusOK, gin.H{"message": message})
}

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

	userId, err := h.userService.LoginUser(user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if userId <= 0 {
		c.JSON(http.StatusForbidden, gin.H{"message": "User not registered, please Signup!!"})
	}

	// if(storedPassword != user.Password){ wrong password entered. }

	c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully"})
}
