package repositories

import (
	"Nookhub/models"
	"database/sql"
	"fmt"
	"log"
)

type RoomsRepository interface {
	GetRooms(userId int) ([]models.Room, error)
	CreateRoom(createrId int, roomName string) (int, error)
	JoinRoom(roomId int, userId int) (bool, error)
	LeaveRoom(roomId int, userId int) (bool, error)
	DeleteRoom(roomId int, userId int) (bool, error)
	SearchRoom(roomName string) (models.Room, error)
	GetHomies(roomId int) ([]models.Homies, error)
	IsAdmin(roomId int, userId int) bool
	CheckRoomIdentity(roomId int) bool
	IsRoomAvailable(roomname string) (bool, error)
}

type roomsRepository struct {
	db *sql.DB
}

func NewRoomsRepository(db *sql.DB) *roomsRepository {
	return &roomsRepository{db: db}
}

func (r *roomsRepository) GetRooms(userId int) ([]models.Room, error) {
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
	var rooms []models.Room
	for rows.Next() {
		var room models.Room
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
	stmt, err := r.db.Prepare("SELECT func_CreateRoom($1, $2)")

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

func (r *roomsRepository) JoinRoom(roomId int, userId int) (bool, error) {
	stmt, err := r.db.Prepare("SLECET func_JoinRoom($1, $2)")

	if err != nil {
		log.Printf("Some error occured while preparing statement %v", err)
		return false, err
	}

	defer stmt.Close()

	var joined bool
	err = stmt.QueryRow(roomId, userId).Scan(&joined)

	if err != nil {
		log.Printf("Some error occured in DB: %v", err)
		return false, err
	}

	return joined, err
}

func (r *roomsRepository) LeaveRoom(roomId int, userId int) (bool, error) {
	stmt, err := r.db.Prepare("SLECET func_LeaveRoom($1, $2)")

	if err != nil {
		log.Printf("Some error occured while preparing statement %v", err)
		return false, err
	}

	defer stmt.Close()

	var left bool
	err = stmt.QueryRow(roomId, userId).Scan(&left)

	if err != nil {
		log.Printf("Some error occured in DB: %v", err)
		return false, err
	}

	return left, err
}

func (r *roomsRepository) DeleteRoom(roomId int, userId int) (bool, error) {
	stmt, err := r.db.Prepare("SLECET func_DeleteRoom($1, $2)")

	if err != nil {
		log.Printf("Some error occured while preparing statement %v", err)
		return false, err
	}

	defer stmt.Close()

	var isDeleted bool
	err = stmt.QueryRow(roomId, userId).Scan(&isDeleted)

	if err != nil {
		log.Printf("Some error occured in DB: %v", err)
		return false, err
	}

	return isDeleted, err
}

func (r *roomsRepository) SearchRoom(roomName string) (models.Room, error) {
	stmt, err := r.db.Prepare("SLECET func_SearchRoom($1)")

	if err != nil {
		log.Printf("Some error occured while preparing statement %v", err)
		return models.Room{}, err
	}

	defer stmt.Close()

	var room models.Room
	err = stmt.QueryRow(roomName).Scan(&room)

	if err != nil {
		log.Printf("Some error occured in DB: %v", err)
		return models.Room{}, err
	}

	return room, err
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

func (r *roomsRepository) IsAdmin(roomId int, userId int) bool {
	var admin = -1
	err := r.db.QueryRow("SELECT admin FROM roomdetails WHERE roomId = $1", roomId).Scan(&admin)

	if err != nil {
		log.Printf("Some error occured while preparing statement %v", err)
		return false
	}

	if admin == userId {
		return true
	}
	return false
}

func (r *roomsRepository) CheckRoomIdentity(roomId int) bool {
	var rows = 0
	err := r.db.QueryRow("SELECT count(*) FROM roomdetails WHERE roomid = $1", roomId).Scan(&rows)

	if err != nil {
		log.Printf("Some error occured while preparing statement %v", err)
		return false
	}

	if rows == 0 {
		return false
	}
	return true
}

func (r *roomsRepository) IsRoomAvailable(roomname string) (bool, error) {
	var rows = 0
	err := r.db.QueryRow("SELECT count(*) FROM roomdetails WHERE roomname = $1", roomname).Scan(&rows)

	if err != nil {
		log.Printf("Some error occured while preparing statement %v", err)
		return false, err
	}

	if rows == 0 {
		return true, nil
	}
	return false, nil
}
