package repositories

import (
	"Nookhub/models"
	"database/sql"
	"fmt"
	"log"
)

type RoomsRepository interface {
	GetRooms(userId int) ([]models.Rooms, error)
	GetHomies(roomId int) ([]models.Homies, error)
	CreateRoom(createrId int, roomName string) (int, error)
}

type roomsRepository struct {
	db *sql.DB
}

func NewRoomsRepository(db *sql.DB) *roomsRepository {
	return &roomsRepository{db: db}
}

func (r *roomsRepository) GetRooms(userId int) ([]models.Rooms, error) {
	stmt, err := r.db.Prepare("SELECT * FROM func_getRoomsByUser($1)")

	if err != nil {
		log.Printf("Some error occured while preparing statement %v", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userId)

	if err != nil {
		log.Printf("Some error occured while getting rooms: %v", err)
		return nil, err
	}
	var rooms []models.Rooms
	for rows.Next() {
		var room models.Rooms
		err := rows.Scan(&room.RoomId, &room.RoomName, &room.RoomIcon)
		if err != nil {
			return nil, fmt.Errorf("error fetching homies 😭: %w", err)
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (r *roomsRepository) CreateRoom(createrId int, roomName string) (int, error) {
	var roomId int
	stmt, err := r.db.Prepare("Select func_CreateRoom($1, $2)")

	if err != nil {
		log.Printf("Some error occured while preparing statement %v", err)
		return 0, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(roomName, createrId).Scan(&roomId)

	if err != nil {
		log.Printf("Some error occured in DB: %v", err)
		return roomId, err
	} else {
		return roomId, err
	}
}

func (r *roomsRepository) GetHomies(roomId int) ([]models.Homies, error) {
	stmt, err := r.db.Prepare("SELECT * FROM func_GetHomies($1)")

	if err != nil {
		log.Printf("Some error occured while preparing statement %v", err)
		return nil, err
	}
	defer stmt.Close()
	var homiesList []models.Homies

	rows, err := stmt.Query(roomId)

	if err != nil {
		log.Printf("Some error occured while preparing statement %v", err)
		return nil, err
	}

	rows.Columns()
	defer rows.Close()

	for rows.Next() {
		var user models.Homies
		err := rows.Scan(&user.Id, &user.Name, &user.Pfp, &user.IsAdmin)
		if err != nil {
			return nil, fmt.Errorf("error fetching homies 😭: %w", err)
		}
		homiesList = append(homiesList, user)
	}
	return homiesList, err
}
