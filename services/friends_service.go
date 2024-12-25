package services

import (
	"Nookhub/models"
	"Nookhub/repositories"
)

type FriendsService interface {
	FetchFriends(userId int) (models.FriendList, error)
	FindUser(input string) (models.Friend, error) // its about searching other users, will rename
	AddFriend(friendreuest models.FriendRequest) (bool, error)
	RequestStatus(friend models.FriendRequest) (bool, error)
}

type friendsService struct {
	friendsRepository repositories.FriendsRepository
}

// FetchFriends implements FriendsService.
func (f *friendsService) FetchFriends(userId int) (models.FriendList, error) {
	return f.friendsRepository.FetchFriends(userId)
}

// AddFriend implements FriendsService.
func (f *friendsService) AddFriend(friendrequest models.FriendRequest) (bool, error) {
	return f.friendsRepository.AddFriend(friendrequest)
}

// FindFriend implements FriendsService.
func (f *friendsService) FindUser(username string) (models.Friend, error) {
	return f.friendsRepository.FindUser(username)
}

// RequestStatus implements FriendsService.
func (f *friendsService) RequestStatus(friendrequest models.FriendRequest) (bool, error) {
	return f.friendsRepository.RequestStatus(friendrequest)
}

func NewFriendsService(friendsRepository repositories.FriendsRepository) *friendsService {
	return &friendsService{friendsRepository: friendsRepository}
}
