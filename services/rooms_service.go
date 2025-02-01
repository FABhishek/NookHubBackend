package services

import (
	"Nookhub/models"
	"Nookhub/repositories"
	"fmt"
)

type RoomsService interface {
	GetRooms(userId int) ([]models.Room, error)
	CreateRoom(createrId int, roomName string) (int, error)
	JoinRoom(roomId int, userId int) (bool, error)
	LeaveRoom(roomId int, userId int) (bool, error)
	DeleteRoom(roomId int, userId int) (bool, error)
	SearchRoom(roomName string) (models.Room, error)
	GetHomies(roomId int) ([]models.Homies, error)
	CheckRoomIdentity(roomId int) bool
	IsRoomAvailable(roomname string) (bool, error)
}

type roomsService struct {
	roomsRepository repositories.RoomsRepository
}

func NewRoomsService(roomsRepository repositories.RoomsRepository) *roomsService {
	return &roomsService{roomsRepository: roomsRepository}
}

func (s *roomsService) GetRooms(userId int) ([]models.Room, error) {
	return s.roomsRepository.GetRooms(userId)
}

func (s *roomsService) CreateRoom(createrId int, roomName string) (int, error) {
	return s.roomsRepository.CreateRoom(createrId, roomName)
}

func (s *roomsService) JoinRoom(roomId int, userId int) (bool, error) {
	return s.roomsRepository.JoinRoom(roomId, userId)
}

func (s *roomsService) LeaveRoom(roomId int, userId int) (bool, error) {
	//check if user is admin of that particular room
	var isAdmin = s.roomsRepository.IsAdmin(roomId, userId)
	if isAdmin {
		return false, fmt.Errorf("please make somene else admin before leaving the room")
	}
	return s.roomsRepository.LeaveRoom(roomId, userId)
}

func (s *roomsService) DeleteRoom(roomId int, userId int) (bool, error) {
	return s.roomsRepository.DeleteRoom(roomId, userId)
}

func (s *roomsService) SearchRoom(roomName string) (models.Room, error) {
	return s.roomsRepository.SearchRoom(roomName)
}

func (s *roomsService) GetHomies(roomId int) ([]models.Homies, error) {
	return s.roomsRepository.GetHomies(roomId)
}

func (s *roomsService) CheckRoomIdentity(roomId int) bool {
	isRoomValid := s.roomsRepository.CheckRoomIdentity(roomId)
	return isRoomValid
}

func (s *roomsService) IsRoomAvailable(roomname string) (bool, error) {
	return s.roomsRepository.IsRoomAvailable(roomname)
}
