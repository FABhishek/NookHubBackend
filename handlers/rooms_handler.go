package handlers

import (
	"Nookhub/services"
	jwtutil "Nookhub/utilities"
	"database/sql"
	"errors"
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
	SearchRoom(c *gin.Context)
	GetHomies(c *gin.Context) // will get all the participants in a room
	IsRoomAvailable(c *gin.Context)
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
	} else if roomId == -1 {
		c.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("room with similar name might already exists, try different name: %v", err)})
		return
	} else {
		c.JSON(http.StatusOK, fmt.Sprintf("room created successfully with name: %s and Id: %d", roomName, roomId))
		return
	}
}

// if users joins a room then he will see leave option inside the room somewhere
// if he searches the room then join option will not appear
func (h *roomsHandler) JoinRoom(c *gin.Context) {
	roomId := c.Param("roomid")
	roomId_No, _ := strconv.Atoi(roomId)
	userId := jwtutil.CheckCookies(c)

	joined, err := h.roomsService.JoinRoom(roomId_No, userId)

	if err != nil {
		if !joined {
			c.JSON(http.StatusBadRequest, gin.H{"error": "please try again after refreshing"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("user %d sucessfully joined the room %d", userId, roomId_No)})
}

// you will get userId from cookies and roomId in url same as join
func (h *roomsHandler) LeaveRoom(c *gin.Context) {
	roomId := c.Param("roomid")
	roomId_No, _ := strconv.Atoi(roomId)
	userId := jwtutil.CheckCookies(c)

	left, err := h.roomsService.LeaveRoom(roomId_No, userId)

	if err != nil {
		if !left {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("user %d sucessfully left the room %d", userId, roomId_No)})
}

// only admin can delete room
func (h *roomsHandler) DeleteRoom(c *gin.Context) {
	roomId := c.Param("roomid")
	roomId_No, err := strconv.Atoi(roomId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Inavalid room id provided make sure it's an integer"})
		return
	}
	userId := jwtutil.CheckCookies(c)

	isRoomValid := h.roomsService.CheckRoomIdentity(roomId_No)
	if !isRoomValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Inavalid room id provided, room with id %d doesn't exist", roomId_No)})
		return
	}

	isDeleted, err := h.roomsService.DeleteRoom(roomId_No, userId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	} else if !isDeleted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You don't have the perms to delete, please ask admin to do so"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("sucessfully deleted the room %d", roomId_No)})
}

// if user is already in the room then room will not have the join option instead leave will be there
// and vice versa in case of not joined.
func (h *roomsHandler) SearchRoom(c *gin.Context) {
	roomName := c.Param("roomname")

	if roomName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "roomName cannot be empty"})
	}
	room, err := h.roomsService.SearchRoom(roomName)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, gin.H{"room": room})
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

func (h *roomsHandler) IsRoomAvailable(c *gin.Context) {
	roomname := c.DefaultQuery("roomname", "")

	availability, err := h.roomsService.IsRoomAvailable(roomname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else if !availability {
		c.JSON(http.StatusConflict, gin.H{"message": fmt.Sprintf("room with name: %s already exists please give different name", roomname)})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "room name available"})
	}
}
