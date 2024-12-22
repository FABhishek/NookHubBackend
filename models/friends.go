package models

type FriendRequest struct {
	UserId     int    `json:"userid"`
	Username   string `json:"username"` //why do we need names? bhul gya bc
	FriendId   int    `json:"friendid"`
	FriendName string `json:"friendname"`
	Status     string `json:"requeststaus"`
}

type Friend struct {
	FriendId   string `json:"friendid"`
	FriendName string `json:"friendname"`
}

type FriendList struct {
	Friends []Friend `json:"friends"`
}
