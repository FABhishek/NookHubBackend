package handlers

import (
	"Nookhub/services"
	"fmt"
	"strings"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type FriendChatHandler interface {
	HandleConnections(c *gin.Context)
	RetreiveMessages(c *gin.Context)
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

func (h *friendChatHandler) RetreiveMessages(c *gin.Context) {
	chatid := c.Param("chatid")
	startAfter := c.DefaultQuery("startAfter", "")

	data, err := h.friendChatService.RetreiveMessages(chatid, startAfter)

	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Some error occured: %v", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *friendChatHandler) HandleConnections(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("There is an error updgrading the http to ws %v", err)
		c.JSON(http.StatusInternalServerError, "Some error occuered while upgrading the connection")
		return
	}

	go h.runSeperately(c, ws) // using go routines, multithreading basically
}

func (h *friendChatHandler) runSeperately(c *gin.Context, ws *websocket.Conn) {
	userId := c.Query("userid")
	chatId := c.Query("chatid")
	user := c.Query("user")

	if strings.TrimSpace(userId) == "" {
		c.JSON(http.StatusBadRequest, "UserId cannot be null or empty")
		return
	}

	defer func() {
		delete(clients, userId)
		log.Printf("user disconnected %s", userId)
	}()

	defer ws.Close()
	h.friendChatService.HandleConnections(ws, userId, chatId, user)
}
