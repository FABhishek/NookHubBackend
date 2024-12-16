package services

import (
	"Nookhub/models"
	"Nookhub/repositories"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type SignupService interface {
	RegisterUser(user models.RegisterUser) (int, error)
	LoginUser(user models.LoginUser) (string, int, error)
	IsAvailable(input string, inputType string) (bool, error)
}

type signupService struct {
	signupRepository repositories.SignupRepository
}

func NewSignupService(signupRepository repositories.SignupRepository) SignupService {
	return &signupService{signupRepository: signupRepository}
}

// signup related services
func (s *signupService) RegisterUser(user models.RegisterUser) (int, error) {

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return 0, fmt.Errorf("problem creating the hash of password: %w", err)
	}
	user.Password = hashedPassword
	return s.signupRepository.CreateUser(user)
}

// EmailChecker implements UserService.
func (s *signupService) IsAvailable(input string, inputType string) (bool, error) {
	return s.signupRepository.IsAvailable(input, inputType)
}

// login related service
// we will get all the users from server to client side to filter it on client side only, inefficient but okay for small amt of data
func (s *signupService) LoginUser(user models.LoginUser) (string, int, error) {
	passHash, username, userId, err := s.signupRepository.LoginUser(user)
	if err != nil {
		return username, -1, fmt.Errorf("some error occured: %w", err)
	} else if userId > 0 {
		if CheckPasswordHash(user.Password, passHash) {
			return username, userId, nil
		} else {
			return "existsButPWNotMatched", userId, fmt.Errorf("provided password does not match")
		}
	}
	return username, userId, fmt.Errorf("user doesn't exist please check username")
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// need to call the function while logging in: first need to retrieve the password hash from table on the basis of userEmail
// or we can get the id and then retrieve it later and then compare the password with it.
func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
