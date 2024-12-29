package handlers

import (
	"Nookhub/services"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type FriendChatHandler interface {
	HandleConnections(c *gin.Context)
}

type friendChatHandler struct {
	friendChatService services.FriendChatService
}

func NewFriendChatHandler(friendChatService services.FriendChatService) *friendChatHandler {
	return &friendChatHandler{friendChatService: friendChatService}
}

var clients = make(map[string]*websocket.Conn)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections;
	},
}

func (h *friendChatHandler) HandleConnections(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer ws.Close()

	userId := c.Query("userid")
	chatId := c.Query("chatid")
	user := c.Query("user")

	if userId == "" {
		c.JSON(http.StatusBadRequest, "UserId cannot be null or empty")
		return
	}

	defer func() {
		delete(clients, userId)
		log.Printf(`user disconnected %s`, userId)
	}()

	h.friendChatService.HandleConnections(ws, userId, chatId, user)
}
