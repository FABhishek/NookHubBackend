package handlers

import (
	"Nookhub/models"
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

// Login related APIs
/*
API can return following http Responses
401: unauthorized: authentication issue
200: success: successfully fetched the friends
*/
func (h *friendsHandler) FetchFriends(c *gin.Context) {
	userId := checkCookies(c)

	friendList, err := h.friendsService.FetchFriends(userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"FriendList": friendList.Friends})
		return
	}
}

// This api will send the friend request to other users
func (h *friendsHandler) AddFriend(c *gin.Context) {
	userId := checkCookies(c)

	var request models.FriendRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if userId != request.UserId {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Somehthing went wrong please login again"})
		return
	}

	success, err := h.friendsService.AddFriend(request)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !success {
		c.JSON(http.StatusBadRequest, "Something is not right, please try again")
		return
	}

	c.JSON(http.StatusOK, gin.H{"friend request sent to": request.FriendName})
}

// FindFriend implements FriendsRepository.
func (h *friendsHandler) FindFriend(c *gin.Context) {
	panic("unimplemented")
}

// RequestStatus implements FriendsRepository.
func (h *friendsHandler) RequestStatus(c *gin.Context) {
	panic("unimplemented")
}

func checkCookies(c *gin.Context) int {
	UserId, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return 0
	}
	return UserId.(int)
}
