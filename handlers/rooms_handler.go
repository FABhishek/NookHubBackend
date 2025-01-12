package handlers

import (
	"Nookhub/services"
	jwtutil "Nookhub/utilities"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoomsHandler interface {
	GetRooms(c *gin.Context)   // will get the rooms a particular user is participant of
	CreateRoom(c *gin.Context) // whoever creates the room will be the admin of that room
	JoinRoom(c *gin.Context)   // will make the user a participant of the room
	LeaveRoom(c *gin.Context)  // will remove the user from a room
	DeleteRoom(c *gin.Context) // delete room can only be performed by admin
	GetHomies(c *gin.Context)  // will get all the participants in a room
}

type roomsHandler struct {
	roomsService services.RoomsService
}

func NewRoomsHandler(roomsService services.RoomsService) *roomsHandler {
	return &roomsHandler{roomsService: roomsService}
}

func (h *roomsHandler) GetRooms(c *gin.Context) {
	userId := jwtutil.CheckCookies(c)

	rooms, err := h.roomsService.GetRooms(userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Some error occured: %v", err)})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": rooms})
		return
	}
}

func (h *roomsHandler) CreateRoom(c *gin.Context) {
	createrId := jwtutil.CheckCookies(c)
	roomName := c.DefaultQuery("roomname", "")

	roomId, err := h.roomsService.CreateRoom(createrId, roomName)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("some error occuered: %v", err)})
		return
	} else {
		c.JSON(http.StatusOK, fmt.Sprintf("room created successfully with name: %s and Id: %d", roomName, roomId))
		return
	}
}

func (h *roomsHandler) JoinRoom(c *gin.Context) {
	panic("unimplemented")
}

func (h *roomsHandler) LeaveRoom(c *gin.Context) {
	panic("unimplemented")
}

func (h *roomsHandler) DeleteRoom(c *gin.Context) {
	panic("unimplemented")
}

func (h *roomsHandler) GetHomies(c *gin.Context) {
	roomid := c.Param("roomid")
	int_roomId, _ := strconv.Atoi(roomid)

	homies, err := h.roomsService.GetHomies(int_roomId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, gin.H{"Homies": homies})
	}
}
