package main

import (
	"Nookhub/config"
	"Nookhub/db"
	"Nookhub/handlers"
	"Nookhub/repositories"
	"Nookhub/routes"
	"Nookhub/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize database
	db.Initialize()

	// Create a new Gin router
	r := gin.Default()

	// cors
	r.Use(cors.Default())
	// Setup dependency injections
	userRepo := repositories.NewUserRepository(db.DB)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Setup routes
	routes.SetupRoutes(r, userHandler)

	// Start the server
	r.Run(":8080")
}
