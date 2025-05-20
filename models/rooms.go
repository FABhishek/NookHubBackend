package models

import "time"

type Room struct {
	RoomName string
	RoomId   int
	RoomIcon *string
}

type RoomMessage struct {
	SenderName string //
	RoomId     int    //room id for which message is intended for
	Content    string
	Date       time.Time //message publish date
	SenderId   int
}
