package handlers

import (
	"Nookhub/models"
	"Nookhub/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type RoomChatHandler interface {
	HandleConnections(c *gin.Context)
}

type roomChatHandler struct {
	roomChatService services.RoomChatService
}

func NewRoomChatHandler(roomChatService services.RoomChatService) *roomChatHandler {
	return &roomChatHandler{roomChatService: roomChatService}
}

var upgraderConn = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections;
	},
}

// I think we can implement a pub-sub here if we get the roomid in request then we will fetch all the users for which
// that message is intended for, and broadcast for all of them.

func (h *roomChatHandler) HandleConnections(c *gin.Context) {
	ws, err := upgraderConn.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("There is an error updgrading the http to ws %v", err)
		c.JSON(http.StatusInternalServerError, "Some error occured while upgrading the connection")
		return
	}

	go h.runSeperately(c, ws)
}

func (h *roomChatHandler) runSeperately(c *gin.Context, ws *websocket.Conn) {

	var request models.RoomMessage

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	userId := c.Param("userId")

	if isAnyEmpty(request.SenderName, userId) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input, either of the required param is empty"})
	}

	// var roomId := request.RoomId

	h.roomChatService.HandleConnections(ws, request, userId)
}
