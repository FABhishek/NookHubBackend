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
	// Setup dependency injections for signup
	signupRepository := repositories.NewSignupRepository(db.DB)
	signupService := services.NewSignupService(signupRepository)
	signupHandler := handlers.NewSignupHandler(signupService)

	// Setup dependency injections for friends
	friendsRepository := repositories.NewFriendsRepository(db.DB)
	friendsService := services.NewFriendsService(friendsRepository)
	friendsHandler := handlers.NewFriendsHandler(friendsService)
	// Setup routes
	routes.SetupRoutes(r, signupHandler, friendsHandler)

	// Start the server
	r.Run(":8080")
}
