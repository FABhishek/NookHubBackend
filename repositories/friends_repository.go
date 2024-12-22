package repositories

import (
	"Nookhub/models"
	"database/sql"
	"fmt"
)

type FriendsRepository interface {
	FetchFriends(userId int) (models.FriendList, error)
	FindFriend(input int) (models.FriendList, error)
	AddFriend(friendrequest models.FriendRequest) (bool, error)
	RequestStatus(friend models.Friend)
}

type friendsRepository struct {
	db *sql.DB
}

// FetchFriends implements FriendsRepository.
func (r *friendsRepository) FetchFriends(userId int) (models.FriendList, error) {

	// Prepare the call to the stored procedure
	stmt, err := r.db.Prepare("SELECT * from func_getUserFriends($1)")
	if err != nil {
		return models.FriendList{}, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()
	var friendList models.FriendList
	// Retrieve the OUT parameter value
	rows, err := stmt.Query(userId)
	if err != nil {
		return models.FriendList{}, fmt.Errorf("error executing function: %w", err)
	}

	rows.Columns()
	defer rows.Close()

	for rows.Next() {
		var friend models.Friend
		err := rows.Scan(&friend.FriendId, &friend.FriendName)
		if err != nil {
			return models.FriendList{}, fmt.Errorf("error scanning friends 😭: %v", err)
		}
		friendList.Friends = append(friendList.Friends, friend)
	}

	if err := rows.Err(); err != nil {
		return models.FriendList{}, fmt.Errorf("some error occured while reading data from DB: %v", err)
	}

	return friendList, nil
}

// AddFriend implements FriendsRepository.
func (r *friendsRepository) AddFriend(friendrequest models.FriendRequest) (bool, error) {

	stmt, err := r.db.Prepare("SELECT func_addFriendRequest($1, $2)")
	if err != nil {
		return false, fmt.Errorf("error executing the procedure: %w", err)
	}

	defer stmt.Close()
	var success bool
	err = stmt.QueryRow(friendrequest.UserId, friendrequest.FriendId).Scan(&success)

	if err != nil {
		return false, fmt.Errorf("some error occured while reading data from DB")
	} else {
		return success, err
	}
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
