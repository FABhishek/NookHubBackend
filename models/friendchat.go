package models

import "time"

type FriendChat struct {
	Sender     string    `json:"sender"`
	Recipient  string    `json:"recipient"`
	Content    string    `json:"content"`
	Date       time.Time `json:"date"`
	FriendName string    `json:"friend"`
}

type Chat struct {
	Text string
	Date time.Time
}

type MessageData struct {
	Sender  string
	Message string
	Date    string
}
