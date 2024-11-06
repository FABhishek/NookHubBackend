package services

import (
	"Nookhub/models"
	"Nookhub/repositories"
	"fmt"

	"golang.org/x/crypto/bcrypt"
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

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return 0, fmt.Errorf("problem creating the hash of password: %w", err)
	}
	user.Password = hashedPassword
	return s.userRepo.CreateUser(user)
}

// we will get all the users from server to client side to filter it on client side only, inefficient but okay for small amt of data
func (s *userService) LoginUser(user models.LoginUser) (int, error) {

	return s.userRepo.LoginUser(user)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// need to call the function while logging in: first need to retrieve the password hash from table on the basis of userEmail
// or we can get the id and then retrieve it later and then compare the password with it.
// func CheckPasswordHash(password string, hash string) bool {
//     err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
//     return err == nil
// }
