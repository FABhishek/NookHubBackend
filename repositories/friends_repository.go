package repositories

import (
	"Nookhub/models"
	"database/sql"
)

type FriendsRepository interface {
	FetchFriends(input int)
	FindFriend(input int) (models.FriendList, error)
	AddFriend(friend models.Friend)
	RequestStatus(friend models.Friend)
}

type friendsRepository struct {
	db *sql.DB
}

// FetchFriends implements FriendsRepository.
func (s *friendsRepository) FetchFriends(input int) {
	panic("unimplemented")
}

// AddFriend implements FriendsRepository.
func (s *friendsRepository) AddFriend(friend models.Friend) {
	panic("unimplemented")
}

// FindFriend implements FriendsRepository.
func (s *friendsRepository) FindFriend(input int) (models.FriendList, error) {
	panic("unimplemented")
}

// RequestStatus implements FriendsRepository.
func (s *friendsRepository) RequestStatus(friend models.Friend) {
	panic("unimplemented")
}

func NewFriendsRepository(db *sql.DB) *friendsRepository {
	return &friendsRepository{db: db}
}
