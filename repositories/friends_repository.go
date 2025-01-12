package repositories

import (
	"Nookhub/models"
	"database/sql"
	"fmt"
)

type FriendsRepository interface {
	FetchFriends(userId int) (models.FriendList, error)
	FindUser(frindname string, userid int) (models.Friend, error)
	AddFriend(friendrequest models.FriendRequest) (bool, error)
	RequestStatus(friendrequest models.FriendRequest) (bool, error)
	PendingRequests(userId int) (models.FriendList, error)
}

type friendsRepository struct {
	db *sql.DB
}

func NewFriendsRepository(db *sql.DB) *friendsRepository {
	return &friendsRepository{db: db}
}

// FetchFriends implements FriendsRepository.
func (r *friendsRepository) FetchFriends(userId int) (models.FriendList, error) {

	// Prepare the call to the stored procedure
	stmt, err := r.db.Prepare("SELECT * FROM func_getUserFriends($1)")
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
		err := rows.Scan(&friend.FriendId, &friend.FriendName, &friend.ChatId, &friend.ProfilePhoto, &friend.Status)
		if err != nil {
			return models.FriendList{}, fmt.Errorf("error scanning friends 😭: %w", err)
		}
		friendList.Friends = append(friendList.Friends, friend)
	}

	if err := rows.Err(); err != nil {
		return models.FriendList{}, fmt.Errorf("some error occured while reading data from DB: %w", err)
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
		return false, fmt.Errorf("some error occured while reading data from DB %w", err)
	} else {
		return success, err
	}
}

// FindFriend implements FriendsRepository.
func (r *friendsRepository) FindUser(friendname string, userid int) (models.Friend, error) {

	stmt, err := r.db.Prepare("SELECT * From func_findUser($1, $2)")

	if err != nil {
		return models.Friend{}, fmt.Errorf("error executing the procedure: %w", err)
	}

	defer stmt.Close()
	var user models.Friend
	err = stmt.QueryRow(friendname, userid).Scan(&user.FriendName, &user.FriendId, &user.ProfilePhoto, &user.Status)

	if err != nil {
		return models.Friend{}, err
	} else {
		return user, nil
	}
}

// RequestStatus implements FriendsRepository.
func (r *friendsRepository) RequestStatus(friendrequest models.FriendRequest) (bool, error) {

	stmt, err := r.db.Prepare("SELECT func_changeRequestStatus($1, $2, $3)")
	if err != nil {
		return false, fmt.Errorf("error executing the procedure: %w", err)
	}

	defer stmt.Close()
	var success bool
	err = stmt.QueryRow(friendrequest.UserId, friendrequest.FriendId, friendrequest.Status).Scan(&success)

	if err != nil {
		return false, fmt.Errorf("some error occured while reading data from DB %w", err)
	} else {
		return success, err
	}
}

func (r *friendsRepository) PendingRequests(userId int) (models.FriendList, error) {
	stmt, err := r.db.Prepare("SELECT * from func_getUserPendingRequests($1)")
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
		err := rows.Scan(&friend.FriendId, &friend.FriendName, &friend.ChatId, &friend.ProfilePhoto, &friend.Status)
		if err != nil {
			return models.FriendList{}, fmt.Errorf("error scanning pending requests 😭: %w", err)
		}
		friendList.Friends = append(friendList.Friends, friend)
	}

	if err := rows.Err(); err != nil {
		return models.FriendList{}, fmt.Errorf("some error occured while reading data from DB: %w", err)
	}

	return friendList, nil
}
