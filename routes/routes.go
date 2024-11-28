package routes

import (
	"Nookhub/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, userHandler handlers.UserHandler) {
	v1 := router.Group("/api/v1")
	{
		user := v1.Group("/users")
		{
			user.POST("/register", userHandler.RegisterUser)
			user.GET("/login", userHandler.LoginUser)
			user.GET("/inputAvailable", userHandler.IsEmailOrUsernameAvailable)
		}
	}
}
