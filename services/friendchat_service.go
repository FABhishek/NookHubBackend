package services

import (
	"Nookhub/repositories"
	"context"

	"github.com/gorilla/websocket"

	"Nookhub/models"
	"log"

	"github.com/redis/go-redis/v9"
)

type FriendChatService interface {
	HandleConnections(ws *websocket.Conn, userId string, chatId string, user string)
}

type friendChatService struct {
	friendChatRepository repositories.FriendChatRepository
}

func NewFriendChatService(friendChatRepository repositories.FriendChatRepository) *friendChatService {
	return &friendChatService{friendChatRepository: friendChatRepository}
}

var clients = make(map[string]*websocket.Conn)

func (s *friendChatService) HandleConnections(ws *websocket.Conn, userId string, chatId string, user string) {

	clients[userId] = ws
	log.Printf(`user connected %s`, userId)

	redisStore := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // No password set
		DB:       0,  // Use default DB
		Protocol: 2,  // Connection protocol
	})
	ctx := context.Background()

	// will check if user has some messages that got missed while he was away

	check, err := redisStore.Exists(ctx, user+chatId).Result()
	if err != nil {
		log.Printf("value does not exists in redis or got some error while doing operation %v", err)
	} else if check != 0 {
		checkIfUserHasSomeMessagesAlready(ctx, ws, user+chatId, redisStore)
	}

	sendRealTimeMessageToFriend(ctx, ws, redisStore, chatId)

}

func checkIfUserHasSomeMessagesAlready(ctx context.Context, ws *websocket.Conn, userChatId string, redisStore *redis.Client) {
	messages := []string{}
	for {
		message, _ := redisStore.RPop(ctx, userChatId).Result()
		messages = append(messages, message)
		log.Print(message)
		length, _ := redisStore.LLen(ctx, userChatId).Result()
		if length == 0 {
			redisStore.Del(ctx, userChatId)
			break
		}
	}
	ws.WriteJSON(messages)

}

func sendRealTimeMessageToFriend(ctx context.Context, ws *websocket.Conn, redisStore *redis.Client, chatId string) {
	for {
		var msg models.FriendChat

		err := ws.ReadJSON(&msg)

		var chat models.Chat
		chat.Text = msg.Content
		chat.Date = msg.Date

		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		// Route the message to the recipient
		if recipientConn, exists := clients[msg.Recipient]; exists {
			err = recipientConn.WriteJSON(chat)

			if err != nil {
				log.Printf("Error sending message to %s: %v", msg.Recipient, err)
			}
		} else {
			// we will store the messages in redis and db if recipient is unavailable
			err := redisStore.LPush(ctx, msg.FriendName+chatId, msg.Content).Err() // will act as queue
			if err != nil {
				panic(err)
			}

			log.Printf("Recipient %s not connected", msg.Recipient)
		}

		//Route the message to the sender itslef
		if senderConn, exists := clients[msg.Sender]; exists {
			err = senderConn.WriteJSON(msg.Content)
			if err != nil {
				log.Printf("Error sending message to %s: %v", msg.Sender, err)
			}
		} else {
			log.Printf("Recipient %s not connected", msg.Sender)
		}
	}
}
