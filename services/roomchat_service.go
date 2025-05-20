package services

import (
	"Nookhub/models"
	"Nookhub/repositories"
	"context"
	"encoding/json"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type RoomChatService interface {
	HandleConnections(ws *websocket.Conn, request models.RoomMessage, userId string)
}

type roomChatService struct {
	roomChatRepository repositories.RoomChatRepository
	redisClient        *redis.Client
}

func NewRoomChatService(roomChatReposiory repositories.RoomChatRepository, redisClient *redis.Client) *roomChatService {
	return &roomChatService{roomChatRepository: roomChatReposiory, redisClient: redisClient}
}

type Room struct {
	clients map[*websocket.Conn]string // Map of active clients
	mu      sync.Mutex                 // Mutex to protect clients map
}

// Global rooms map
var (
	rooms      = make(map[int]*Room) // Maps roomID to Room struct
	roomsMutex sync.Mutex            // Mutex for the rooms map
)

func (s *roomChatService) HandleConnections(ws *websocket.Conn, request models.RoomMessage, userId string) {
	roomID := request.RoomId

	roomsMutex.Lock()
	if _, exists := rooms[roomID]; !exists {
		rooms[roomID] = &Room{clients: make(map[*websocket.Conn]string)}
	}
	room := rooms[roomID]
	roomsMutex.Unlock()

	// Add clients to the room
	room.mu.Lock()
	room.clients[ws] = userId
	room.mu.Unlock()

	log.Println("User", userId, "connected to room:", roomID)

	// deliver missed messages via redis
	s.deliverMissedMessages(userId, request.RoomId, ws)

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("Client disconnected:", err)
			break
		}
		s.storeOrBroadcastMessage(room, userId, request.RoomId, msg) // Send message to all clients in the room, or store in redis if (offline)
		// store message in nosql db as well.
	}

	defer ws.Close()

	room.mu.Lock()
	delete(room.clients, ws)
	room.mu.Unlock()
}

func (s *roomChatService) storeOrBroadcastMessage(room *Room, sender string, roomId int, msg []byte) {
	room.mu.Lock()
	defer room.mu.Unlock()

	offlineUsers := make(map[string]bool) // Track offline users

	for client, userId := range room.clients {
		err := client.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Error sending message to user", userId, ":", err)
			client.Close()
			delete(room.clients, client)
			offlineUsers[userId] = true
		}
	}

	// Store message for offline users in redis
	for userId := range offlineUsers {
		s.storeMessageInRedis(userId, sender, roomId, msg)
	}
}

func (s *roomChatService) deliverMissedMessages(userId string, roomId int, ws *websocket.Conn) {
	ctx := context.Background()
	redisKey := strconv.Itoa(roomId) + userId

	messages, err := s.redisClient.LRange(ctx, redisKey, 0, -1).Result()
	if err != nil {
		log.Println("Error fetching missed messages:", err)
		return
	}

	// Send each missed message to the user
	for _, msg := range messages {
		err := ws.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Println("Error delivering missed message to user", userId, ":", err)
			return
		}
	}

	// Clear messages from Redis after sending
	s.redisClient.Del(ctx, redisKey)
}

func (s *roomChatService) storeMessageInRedis(userId string, sender string, roomId int, msg []byte) {
	ctx := context.Background()
	userId_str, _ := strconv.Atoi(userId)

	message := models.RoomMessage{
		SenderName: sender,
		RoomId:     roomId,
		Content:    string(msg),
		Date:       time.Now(),
		SenderId:   userId_str,
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Println("Error marshalling message:", err)
		return
	}

	redisKey := strconv.Itoa(roomId) + userId
	s.redisClient.RPush(ctx, redisKey, data) // Append message to the user's missed message list
	// s.roomChatRepository.StoreChatData(data); // store data in firebase.
}
