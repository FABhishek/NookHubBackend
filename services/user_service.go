package services

import (
	"Nookhub/models"
	"Nookhub/repositories"
)

type UserService interface {
	RegisterUser(user models.User) (int,error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) RegisterUser(user models.User) (int,error) {
	// Business logic, e.g., hashing password
	return s.userRepo.Create(user)
}
