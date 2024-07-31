package services

import (
	"Nookhub/models"
	"Nookhub/repositories"
)

type UserService interface {
	RegisterUser(user models.User) (int, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) RegisterUser(user models.User) (int, error) {
	// encrypt the password and it will be decrypted only while signing in
	return s.userRepo.Create(user)
}
