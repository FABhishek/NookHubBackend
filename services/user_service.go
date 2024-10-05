package services

import (
	"Nookhub/models"
	"Nookhub/repositories"
)

type UserService interface {
	RegisterUser(user models.RegisterUser) (int, error)
	LoginUser(user models.LoginUser) (int, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) RegisterUser(user models.RegisterUser) (int, error) {
	// encrypt the password and it will be decrypted only while signing in
	return s.userRepo.CreateUser(user)
}

func (s *userService) LoginUser(user models.LoginUser) (int, error) {

	return s.userRepo.LoginUser(user)
}
