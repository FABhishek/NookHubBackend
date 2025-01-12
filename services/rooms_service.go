package services

import (
	"Nookhub/models"
	"Nookhub/repositories"
)

type RoomsService interface {
	GetRooms(userId int) ([]models.Rooms, error)
	CreateRoom(createrId int, roomName string) (int, error)
	GetHomies(roomId int) ([]models.Homies, error)
}

type roomsService struct {
	roomsRepository repositories.RoomsRepository
}

func NewRoomsService(roomsRepository repositories.RoomsRepository) *roomsService {
	return &roomsService{roomsRepository: roomsRepository}
}

func (s *roomsService) GetRooms(userId int) ([]models.Rooms, error) {
	return s.roomsRepository.GetRooms(userId)
}

func (s *roomsService) CreateRoom(createrId int, roomName string) (int, error) {
	return s.roomsRepository.CreateRoom(createrId, roomName)
}

func (s *roomsService) GetHomies(roomId int) ([]models.Homies, error) {
	return s.roomsRepository.GetHomies(roomId)
}
