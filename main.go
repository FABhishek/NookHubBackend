package main

import (
	"Nookhub/config"
	"Nookhub/db"
	"Nookhub/handlers"

	"Nookhub/repositories"
	"Nookhub/routes"
	"Nookhub/services"

	"context"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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

	redisStore := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // No password set
		DB:       0,  // Use default DB
	})

	_, err := redisStore.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	// Setup dependency injections for signup
	signupRepository := repositories.NewSignupRepository(db.DB)
	signupService := services.NewSignupService(signupRepository)
	signupHandler := handlers.NewSignupHandler(signupService)

	// Setup dependency injections for friends
	friendsRepository := repositories.NewFriendsRepository(db.DB)
	friendsService := services.NewFriendsService(friendsRepository)
	friendsHandler := handlers.NewFriendsHandler(friendsService)

	// Setup dependency injection for friendsChat
	friendChatRepository := repositories.NewFriendChatRepository(db.DB)
	friendChatService := services.NewFriendChatService(friendChatRepository, redisStore)
	friendChatHandler := handlers.NewFriendChatHandler(friendChatService)

	// Setup routes
	routes.SetupRoutes(router, signupHandler, friendsHandler, friendChatHandler)

	// Start the server
	router.Run("0.0.0.0:8080")
}
