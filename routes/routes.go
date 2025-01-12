package routes

import (
	"Nookhub/handlers"
	jwtutil "Nookhub/utilities"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine,
	signupHandler handlers.SignupHandler,
	friendsHandler handlers.FriendsHandler,
	friendChatHandler handlers.FriendChatHandler,
	roomsHandler handlers.RoomsHandler) {

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
			friends.GET("/searchuser", jwtutil.AuthenticateMiddleware, friendsHandler.FindUser) //local friend search
			friends.POST("/requestsent", jwtutil.AuthenticateMiddleware, friendsHandler.AddFriend)
			friends.PUT("/requeststatus", jwtutil.AuthenticateMiddleware, friendsHandler.RequestStatus)
			friends.GET("/pendingrequests", jwtutil.AuthenticateMiddleware, friendsHandler.PendingRequests)
			// we will pass the query param on that basis will delete the entry if request is declined, in put itself
		}

		friendChat := v1.Group("/dashboard/friends/friendchat")
		{
			friendChat.GET("/ws", friendChatHandler.HandleConnections)
			friendChat.GET("/:chatid/messages", friendChatHandler.RetreiveMessages)
		}

		Rooms := v1.Group("/dashboard/rooms")
		{
			Rooms.GET("/", jwtutil.AuthenticateMiddleware, roomsHandler.GetRooms)
			Rooms.POST("/create", jwtutil.AuthenticateMiddleware, roomsHandler.CreateRoom) // add queryparam as roomname
			Rooms.GET("/:roomid/gethomies", jwtutil.AuthenticateMiddleware, roomsHandler.GetHomies)
			Rooms.POST("/join/:roomid", jwtutil.AuthenticateMiddleware, roomsHandler.JoinRoom)
			Rooms.DELETE("/leave/:roomid", jwtutil.AuthenticateMiddleware, roomsHandler.LeaveRoom)
			Rooms.DELETE("/deleteroom/:roomid", jwtutil.AuthenticateMiddleware, roomsHandler.LeaveRoom)
		}

		// RoomChat := v1.Group("/dashboard/rooms/:roomid/chat") {

		// }
	}
}
