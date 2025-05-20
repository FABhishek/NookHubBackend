package repositories

import (
	"Nookhub/models"
	"context"
	"database/sql"
	"log"
	"strconv"
)

type RoomChatRepository interface {
	StoreChatData(message models.RoomMessage)
}

type roomChatRepository struct {
	db         *sql.DB
	configData []byte // firebase config settings.
}

func NewRoomChatRepository(db *sql.DB, configData []byte) *roomChatRepository {
	return &roomChatRepository{db: db, configData: configData}
}

func (r *roomChatRepository) StoreChatData(msg models.RoomMessage) {
	ctx := context.Background()

	FirebaseClient, _ := initializeFirebaseClient(ctx, r.configData)

	var data models.MessageData
	data.Sender = msg.SenderName
	data.Message = msg.Content
	data.Date = msg.Date.String()

	_, err := FirebaseClient.Collection("roomChats").
		Doc(strconv.Itoa(msg.RoomId)).
		Collection("messages").
		Doc(data.Date).
		Set(ctx, data)

	if err != nil {
		log.Printf("Failed to store message: %v", err)
	} else {
		log.Printf("Message stored successfully for chatID: %s", msg.RoomId)
	}

	defer FirebaseClient.Close()
}
