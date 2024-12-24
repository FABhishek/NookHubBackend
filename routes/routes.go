package routes

import (
	"Nookhub/handlers"
	jwtutil "Nookhub/utilities"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine,
	signupHandler handlers.SignupHandler,
	friendsHandler handlers.FriendsHandler) {

	v1 := router.Group("/api/v1")
	{
		signup := v1.Group("/users")
		{
			signup.POST("/register", signupHandler.RegisterUser)
			signup.POST("/login", signupHandler.LoginUser)
			signup.GET("/inputAvailable", signupHandler.IsEmailOrUsernameAvailable)
		}

		friends := v1.Group("/dashboard/friends")
		{
			friends.GET("/fetchfriends", jwtutil.AuthenticateMiddleware, friendsHandler.FetchFriends)
			friends.GET("/search", friendsHandler.FindFriend) //local friend search
			friends.POST("/requestsent", jwtutil.AuthenticateMiddleware, friendsHandler.AddFriend)
			friends.PUT("/acceptrequest", friendsHandler.RequestStatus)
			// we will pass the query param on that basis will delete the entry if request is declined, in put itself
		}
	}
}
