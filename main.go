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
	router := gin.Default()

	// cors
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Allow the frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true, // Allow cookies to be sent with cross-origin requests
	}))

	// Setup dependency injections for signup
	signupRepository := repositories.NewSignupRepository(db.DB)
	signupService := services.NewSignupService(signupRepository)
	signupHandler := handlers.NewSignupHandler(signupService)

	// Setup dependency injections for friends
	friendsRepository := repositories.NewFriendsRepository(db.DB)
	friendsService := services.NewFriendsService(friendsRepository)
	friendsHandler := handlers.NewFriendsHandler(friendsService)
	// Setup routes
	routes.SetupRoutes(router, signupHandler, friendsHandler)

	// Start the server
	router.Run(":8080")
}
