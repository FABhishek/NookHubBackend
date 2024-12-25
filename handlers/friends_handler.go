package handlers

import (
	"Nookhub/models"
	"Nookhub/services"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FriendsHandler interface {
	FetchFriends(c *gin.Context)  // will fetch the user's friends
	FindUser(c *gin.Context)      // will fetch the searched user from db
	AddFriend(c *gin.Context)     // will add the friend and set the status as pending
	RequestStatus(c *gin.Context) // will update the status of request if approved or declined or withdraw
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
// API can return following http Responses
// 401: unauthorized: authentication issue
// 200: success: successfully sent the request
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
		c.JSON(http.StatusBadRequest, "Something went wrong, please try again")
		return
	}

	c.JSON(http.StatusOK, gin.H{"friend request sent to": request.FriendName})
}

// FindFriend implements FriendsRepository.
func (h *friendsHandler) FindUser(c *gin.Context) {
	userId := checkCookies(c)

	username := c.DefaultQuery("username", "")
	userid := c.DefaultQuery("userid", "")
	int_userid, err := strconv.Atoi(userid)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if userId != int_userid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Something went wrong please login again"})
		return
	}

	if len(username) < 3 || username == "" {
		c.JSON(http.StatusBadRequest, "Please enter a valid username")
		return
	}

	user, err := h.friendsService.FindUser(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"result": user})
		return
	}
}

// This api will send the friend request to other users
// API can return following http Responses
// 401: unauthorized: authentication issue
// 200: success: accepted or rejected the request
func (h *friendsHandler) RequestStatus(c *gin.Context) {
	userId := checkCookies(c)

	var request models.FriendRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if userId != request.UserId {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Something went wrong please login again"})
		return
	}

	res, err := h.friendsService.RequestStatus(request)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if !res {
		c.JSON(http.StatusBadRequest, "Something is not right, please try again")
		return
	} else {
		if request.Status == "approved" {
			c.JSON(http.StatusOK, fmt.Sprintf("Friend added: %s", request.FriendName))
			return
		} else {
			c.JSON(http.StatusOK, fmt.Sprintf("Friend request declined: %s", request.FriendName))
			return
		}
	}
}

func checkCookies(c *gin.Context) int {
	UserId, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return 0
	}
	return UserId.(int)
}
