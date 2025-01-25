package services

import (
	"Nookhub/models"
	"Nookhub/repositories"
)

type RoomsService interface {
	GetRooms(userId int) ([]models.Room, error)
	CreateRoom(createrId int, roomName string) (int, error)
	JoinRoom(roomId int, userId int) (bool, error)
	LeaveRoom(roomId int, userId int) (bool, error)
	DeleteRoom(roomId int, userId int) (bool, error)
	SearchRoom(roomName string) (models.Room, error)
	GetHomies(roomId int) ([]models.Homies, error)
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
