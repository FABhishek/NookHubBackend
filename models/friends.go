package models

type FriendRequest struct {
	UserId     int    `json:"userid"`
	Username   string `json:"username"` //why do we need names? bhul gya bc
	FriendId   int    `json:"friendid"`
	FriendName string `json:"friendname"`
	Status     string `json:"requeststaus"`
}

type Friend struct {
	FriendId     int     `json:"friendid"`
	FriendName   string  `json:"friendname"`
	ProfilePhoto *string `json:"profilephoto"` //nullable
	ChatId       *string `json:"chatid"`
	Status       *string `json:"status"`
}

type FriendList struct {
	Friends []Friend `json:"friends"`
}
