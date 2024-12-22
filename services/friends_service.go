package services

import (
	"Nookhub/models"
	"Nookhub/repositories"
)

type FriendsService interface {
	FetchFriends(userId int) (models.FriendList, error)
	FindFriend(input int) (models.FriendList, error) // its about searching other users, will rename
	AddFriend(friend models.Friend)
	RequestStatus(friend models.Friend)
}

type friendsService struct {
	friendsRepository repositories.FriendsRepository
}

// FetchFriends implements FriendsService.
func (f *friendsService) FetchFriends(userId int) (models.FriendList, error) {
	return f.friendsRepository.FetchFriends(userId)
}

// AddFriend implements FriendsService.
func (f *friendsService) AddFriend(friend models.Friend) {
	panic("unimplemented")
}

// FindFriend implements FriendsService.
func (f *friendsService) FindFriend(input int) (models.FriendList, error) {
	panic("unimplemented")
}

// RequestStatus implements FriendsService.
func (f *friendsService) RequestStatus(friend models.Friend) {
	panic("unimplemented")
}

func NewFriendsService(friendsRepository repositories.FriendsRepository) *friendsService {
	return &friendsService{friendsRepository: friendsRepository}
}
