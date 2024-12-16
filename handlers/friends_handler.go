package handlers

import (
	"Nookhub/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FriendsHandler interface {
	FetchFriends(c *gin.Context)
	FindFriend(c *gin.Context)
	AddFriend(c *gin.Context)
	RequestStatus(c *gin.Context)
}

type friendsHandler struct {
	friendsService services.FriendsService
}

func NewFriendsHandler(friendsService services.FriendsService) *friendsHandler {
	return &friendsHandler{friendsService: friendsService}
}

// FetchFriends implements FriendsRepository.
func (h *friendsHandler) FetchFriends(c *gin.Context) {
	userId := c.DefaultQuery("userid", "")
	id, err := strconv.Atoi(userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "There is some issue while processing the given user"})
	}

	friendList, err := h.friendsService.FetchFriends(id)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"FriendList": friendList})
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}
}

// AddFriend implements FriendsRepository.
func (h *friendsHandler) AddFriend(c *gin.Context) {
	panic("unimplemented")
}

// FindFriend implements FriendsRepository.
func (h *friendsHandler) FindFriend(c *gin.Context) {
	panic("unimplemented")
}

// RequestStatus implements FriendsRepository.
func (h *friendsHandler) RequestStatus(c *gin.Context) {
	panic("unimplemented")
}
