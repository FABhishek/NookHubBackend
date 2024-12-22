package handlers

import (
	"Nookhub/services"
	"net/http"

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
	userId, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	friendList, err := h.friendsService.FetchFriends(userId.(int))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"FriendList": friendList.Friends})
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
