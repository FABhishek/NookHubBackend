package repositories

import "database/sql"

type FriendChatRepository interface {
}

type friendChatRepository struct {
	db *sql.DB
}

func NewFriendChatRepository(db *sql.DB) *friendChatRepository {
	return &friendChatRepository{db: db}
}
