package repositories

import (
	"Nookhub/models"
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

var FirebaseClient *firestore.Client

type FriendChatRepository interface {
	StoreChatData(msg models.FriendChat, chatid string)
	RetreiveMessages(chatid string, startAfter string) ([]models.MessageData, error)
}

type friendChatRepository struct {
	chatDB     *firestore.Client
	configData []byte // firebase config settings.
}

func NewFriendChatRepository(chatDB *firestore.Client, configData []byte) *friendChatRepository {
	return &friendChatRepository{chatDB: chatDB, configData: configData}
}

// to store messages in firebase
func (r *friendChatRepository) StoreChatData(msg models.FriendChat, chatid string) {
	ctx := context.Background()

	FirebaseClient, _ := initializeFirebaseClient(ctx, r.configData)

	var data models.MessageData
	data.Sender = msg.Sender
	data.Message = msg.Content
	data.Date = msg.Date.String()

	_, err := FirebaseClient.Collection("chats").
		Doc(chatid).
		Collection("messages").
		Doc(data.Date).
		Set(ctx, data)

	if err != nil {
		log.Printf("Failed to store message: %v", err)
	} else {
		log.Printf("Message stored successfully for chatID: %s", chatid)
	}

	defer FirebaseClient.Close()
}

// to retrieve stored messages from firestore
func (r *friendChatRepository) RetreiveMessages(chatid string, startAfter string) ([]models.MessageData, error) {
	ctx := context.Background()

	// Create a Firestore client
	FirebaseClient, err := initializeFirebaseClient(ctx, r.configData)

	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
		return nil, err
	}

	var storedChat []models.MessageData

	crap, err := FirebaseClient.Collection("chats").
		Doc(chatid).
		Collection("messages").
		OrderBy("Date", firestore.Asc).
		StartAfter(startAfter).
		Limit(20).
		Documents(ctx).GetAll()

	if err != nil {
		log.Printf("Error occured while reading some data")
		return nil, err
	}

	defer FirebaseClient.Close()

	for i := len(crap) - 1; i >= 0; i-- {
		message := crap[i]

		data := message.Data()
		if data == nil {
			continue
		}

		date, _ := data["Date"].(string)

		messageText, _ := data["Message"].(string)

		sender, _ := data["Sender"].(string)

		storedChat = append(storedChat, models.MessageData{Date: date, Message: messageText, Sender: sender})
	}
	return storedChat, nil
}

func initializeFirebaseClient(ctx context.Context, configData []byte) (*firestore.Client, error) {
	client, err := firestore.NewClient(ctx, "nookhubchats", option.WithCredentialsJSON(configData))
	if err != nil {
		return nil, err
	}
	return client, nil
}
